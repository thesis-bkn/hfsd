package handler

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"mime/multipart"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/ztrue/tracerr"

	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/thesis-bkn/hfsd/internal/entity"
	"github.com/thesis-bkn/hfsd/internal/errors"
	"github.com/thesis-bkn/hfsd/internal/s3"
	"github.com/thesis-bkn/hfsd/internal/utils"
	"github.com/thesis-bkn/hfsd/internal/worker"
)

type InferenceHandler struct {
	cc       chan<- interface{}
	validate *validator.Validate
	s3Client *s3.Client
	client   database.Client
	cfg      *config.Config
}

func NewInferenceHandler(
	w *worker.Worker,
	cfg *config.Config,
	s3Client *s3.Client,
	client database.Client,
	validate *validator.Validate,
) *InferenceHandler {
	return &InferenceHandler{
		cc:       w.C,
		client:   client,
		s3Client: s3Client,
		cfg:      cfg,
		validate: validate,
	}
}

type maskObject struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	R float64 `json:"r"`
}

type inferRequest struct {
	ModelID   string        `json:"model"`
	Prompt    string        `json:"prompt"`
	Image     string        `json:"image"`
	NegPrompt string        `json:"negPrompt"`
	Mask      []*maskObject `json:"mask"`
}

func (h *InferenceHandler) SubmitInferenceTask(c echo.Context) error {
	var req inferRequest
	if err := c.Bind(&req); err != nil {
		return tracerr.Wrap(err)
	}

	sourceModel, err := h.client.Query().GetModelByID(c.Request().Context(), req.ModelID)
	if err != nil {
		c.Error(errors.ErrNotFound)
		return tracerr.Wrap(err)
	}

	sourceModelAgg := entity.NewModelFromDB(&sourceModel)
	inf := entity.NewInference(sourceModelAgg, req.Prompt, req.NegPrompt)

	b64 := req.Image[strings.IndexByte(req.Image, ',')+1:]
	imageB, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return tracerr.Wrap(err)
	}

	mask := image.NewRGBA(image.Rect(0, 0, 512, 512))
	draw.Draw(mask, mask.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)
	for _, obj := range req.Mask {
		drawCircle(mask, int(obj.X), int(obj.Y), int(obj.R), color.White)
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, mask, nil); err != nil {
		return tracerr.Wrap(err)
	}
	maskB := buf.Bytes()

	tx, err := h.client.Conn().BeginTx(c.Request().Context(), pgx.TxOptions{})
	if err != nil {
		return tracerr.Wrap(err)
	}
	defer func() {
		if err != nil {
			tx.Rollback(c.Request().Context())
		} else {
			tx.Commit(c.Request().Context())
		}
	}()

	query := h.client.Query().WithTx(tx)
	if err := utils.SavePNG(inf.ImagePath(), imageB); err != nil {
		return tracerr.Wrap(err)
	}

	if err := utils.
		SavePNG(inf.MaskPath(), maskB); err != nil {
		return tracerr.Wrap(err)
	}

	if err := query.InsertInference(
		c.Request().Context(),
		inf.Insertion()); err != nil {
		return tracerr.Wrap(err)
	}

	h.cc <- inf

	return nil
}

func readAll(f *multipart.FileHeader) ([]byte, error) {
	// Open the file
	src, err := f.Open()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	defer src.Close()

	// Read the file content into a byte slice
	fileData, err := io.ReadAll(src)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return fileData, nil
}

func drawCircle(img draw.Image, x, y, r int, c color.Color) {
	for dx := -r; dx < r; dx++ {
		for dy := -r; dy < r; dy++ {
			if dx*dx+dy*dy < r*r {
				img.Set(x+dx, y+dy, c)
			}
		}
	}
}
