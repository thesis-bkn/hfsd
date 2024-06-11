package handler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"mime/multipart"
	"path"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/labstack/echo/v4"
	"github.com/ztrue/tracerr"

	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/thesis-bkn/hfsd/internal/s3"
	"github.com/thesis-bkn/hfsd/internal/worker"
)

const (
	modelIDsAlphabet = `0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ`
)

type InferenceHandler struct {
	w        *worker.Worker
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
		w:        w,
		client:   client,
		s3Client: s3Client,
		cfg:      cfg,
		validate: validate,
	}
}

type maskObject struct {
	X int `json:"x"`
	Y int `json:"y"`
	R int `json:"r"`
}

type inferRequest struct {
	Model  string        `json:"model"`
	Prompt string        `json:"prompt"`
	Image  string        `json:"image"`
	Mask   []*maskObject `json:"mask"`
}

func (i *InferenceHandler) SubmitInferenceTask(c echo.Context) error {
	var req inferRequest
	if err := c.Bind(&req); err != nil {
		return tracerr.Wrap(err)
	}

	b64 := req.Image[strings.IndexByte(req.Image, ',')+1:]
	imageB, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return tracerr.Wrap(err)
	}

	mask := image.NewRGBA(image.Rect(0, 0, 512, 512))
	draw.Draw(mask, mask.Bounds(), &image.Uniform{color.Black}, image.Point{}, draw.Src)
	for _, obj := range req.Mask {
		drawCircle(mask, obj.X, obj.Y, obj.R, color.White)
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, mask, nil); err != nil {
		return tracerr.Wrap(err)
	}
	maskB := buf.Bytes()

	tx, err := i.client.Conn().BeginTx(c.Request().Context(), pgx.TxOptions{})
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

	query := i.client.Query().WithTx(tx)
	order := 0
	var taskID int32
	if taskID, err = query.InsertInferenceTask(c.Request().Context(), req.Model); err != nil {
		return tracerr.Wrap(err)
	}

	imageURL := path.Join(i.cfg.ImagePath, fmt.Sprintf("%d-%d", order, taskID))
	maskUrl := path.Join(i.cfg.MaskPath, fmt.Sprintf("%d-%d", order, taskID))

	if err := i.s3Client.UploadImage(imageB, imageURL); err != nil {
		return tracerr.Wrap(err)
	}
	if err := i.s3Client.UploadImage(maskB, maskUrl); err != nil {
		return tracerr.Wrap(err)
	}

	if err := query.InsertAsset(c.Request().Context(), database.InsertAssetParams{
		TaskID:   taskID,
		Order:    0,
		Prompt:   req.Prompt,
		Image:    imageB,
		ImageUrl: imageURL,
		Mask:     maskB,
		MaskUrl: pgtype.Text{
			String: maskUrl,
			Valid:  true,
		},
	}); err != nil {
		return tracerr.Wrap(err)
	}

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
