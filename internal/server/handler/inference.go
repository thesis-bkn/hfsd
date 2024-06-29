package handler

import (
	"image"
	"image/color"
	"image/draw"
	"io"
	"mime/multipart"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/ztrue/tracerr"

	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/thesis-bkn/hfsd/internal/entity"
	"github.com/thesis-bkn/hfsd/internal/errors"
	"github.com/thesis-bkn/hfsd/internal/s3"
	"github.com/thesis-bkn/hfsd/internal/utils"
)

type InferenceHandler struct {
	cc       chan<- entity.Task
	validate *validator.Validate
	s3Client *s3.Client
	client   database.Client
	cfg      *config.Config
}

func NewInferenceHandler(
    taskqueue chan<- entity.Task,
	cfg *config.Config,
	s3Client *s3.Client,
	client database.Client,
	validate *validator.Validate,
) *InferenceHandler {
	return &InferenceHandler{
		cc:       taskqueue,
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

	mask := image.NewRGBA(image.Rect(0, 0, 512, 512))
	draw.Draw(mask, mask.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)
	for _, obj := range req.Mask {
		drawCircle(mask, int(obj.X), int(obj.Y), int(obj.R), color.White)
	}

	tx, err := h.client.Conn().Begin(c.Request().Context())
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
	if err := utils.SaveBase64ImageToFile(req.Image, inf.ImagePath()); err != nil {
		return tracerr.Wrap(err)
	}

	if err := utils.SaveRGBAImageToJPEG(mask, inf.MaskPath()); err != nil {
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
