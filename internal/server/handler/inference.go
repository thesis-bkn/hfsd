package handler

import (
	"encoding/base64"
	"fmt"
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

type maskObject struct {
	X float32 `json:"x"`
	Y float32 `json:"y"`
	R float32 `json:"r"`
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
	image, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return tracerr.Wrap(err)
	}

	// FIXME: use req.Masks to create mask image
	mask := image

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

	if err := i.s3Client.UploadImage(image, imageURL); err != nil {
		return tracerr.Wrap(err)
	}
	if err := i.s3Client.UploadImage(mask, maskUrl); err != nil {
		return tracerr.Wrap(err)
	}

	if err := query.InsertAsset(c.Request().Context(), database.InsertAssetParams{
		TaskID:   taskID,
		Order:    0,
		Prompt:   req.Prompt,
		Image:    image,
		ImageUrl: imageURL,
		Mask:     mask,
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
