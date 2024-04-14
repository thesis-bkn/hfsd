package handler

import (
	"fmt"
	"io"
	"mime/multipart"
	"path"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/ztrue/tracerr"

	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/thesis-bkn/hfsd/internal/database"
	"github.com/thesis-bkn/hfsd/internal/s3"
)

const (
	modelIDsAlphabet = `0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ`
)

type InferenceHandler struct {
	validate *validator.Validate
	s3Client *s3.Client
	client   database.Client
	cfg      *config.Config
}

func NewInferenceHandler(
	cfg *config.Config,
	s3Client *s3.Client,
	client database.Client,
	validate *validator.Validate,
) *InferenceHandler {
	return &InferenceHandler{
		client:   client,
		s3Client: s3Client,
		cfg:      cfg,
		validate: validate,
	}
}

func (i *InferenceHandler) SubmitInferenceTask(c echo.Context) error {
	model := c.FormValue("model")
	prompt := c.FormValue("prompt")

	image, err := c.FormFile("image")
	if err != nil {
		return tracerr.Wrap(err)
	}

	mask, err := c.FormFile("mask")
	if err != nil {
		return tracerr.Wrap(err)
	}

	imageB, err := readAll(image)
	if err != nil {
		return tracerr.Wrap(err)
	}
	maskB, err := readAll(mask)
	if err != nil {
		return tracerr.Wrap(err)
	}

	taskID, err := gonanoid.Generate(modelIDsAlphabet, 5)
	if err != nil {
		return tracerr.Wrap(err)
	}

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
	if err := query.InsertInferenceTask(c.Request().Context(), model); err != nil {
		return tracerr.Wrap(err)
	}

	imageURL := path.Join(i.cfg.ImagePath, fmt.Sprintf("%d-%s", order, taskID))
	maskUrl := path.Join(i.cfg.MaskPath, fmt.Sprintf("%d-%s", order, taskID))

	if err := i.s3Client.UploadImage(imageB, imageURL); err != nil {
		return tracerr.Wrap(err)
	}
	if err := i.s3Client.UploadImage(maskB, maskUrl); err != nil {
		return tracerr.Wrap(err)
	}

	if err := query.InsertAsset(c.Request().Context(), database.InsertAssetParams{
		Prompt:   prompt,
		Order:    0,
		Image:    imageB,
		ImageUrl: imageURL,
		Mask:     maskB,
		MaskUrl:  maskUrl,
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
