package s3

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/thesis-bkn/hfsd/internal/config"
	"github.com/ztrue/tracerr"
)

// Client struct represents the S3 image uploader
type Client struct {
	sess       *session.Session
	bucketName string
}

// NewS3Client creates a new instance of S3ImageUploader
func NewS3Client(cfg *config.Config) *Client {
	awsCfg := aws.NewConfig().
		WithEndpoint(cfg.EndpointUrl).
		WithCredentials(credentials.NewStaticCredentials(
			cfg.AwsAccessKeyID,
			cfg.AwsSecretAccessKey,
			"")).
		WithRegion("hn")
	return &Client{
		bucketName: cfg.Bucket,
		sess:       session.Must(session.NewSession(awsCfg)),
	}
}

// UploadImage uploads an image to Amazon S3 and returns a pre-signed URL
func (uploader *Client) UploadImage(imageData []byte, s3Key string) error {
	// Upload image to S3
	if _, err := s3.New(uploader.sess).PutObject(&s3.PutObjectInput{
		Bucket: aws.String(uploader.bucketName),
		Key:    aws.String(s3Key),
		Body:   bytes.NewReader(imageData),
	}); err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}
