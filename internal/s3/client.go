package s3

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ztrue/tracerr"
)

// S3Client struct represents the S3 image uploader
type S3Client struct {
	sess       *session.Session
	bucketName string
}

// NewS3Client creates a new instance of S3ImageUploader
func NewS3Client(bucketName string) *S3Client {
	sess := session.Must(session.NewSession())
	return &S3Client{
		sess:       sess,
		bucketName: bucketName,
	}
}

func (client *S3Client) GeneratePresignedURL(s3Key string) (string, error) {
	svc := s3.New(client.sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(client.bucketName),
		Key:    aws.String(s3Key),
	})

	url, err := req.Presign(1 * time.Hour)
	if err != nil {
		return "", tracerr.Wrap(err)
	}

	return url, nil
}
