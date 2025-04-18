package SelectelS3

import (
	"VK_posts/internal/config"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"log"
	"mime/multipart"
	"strings"
)

type SelectelS3 struct {
	session    *session.Session
	bucketName string
}

func NewSelectelS3(cfg *config.Config) *SelectelS3 {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Selectel.Region),
		Credentials: credentials.NewStaticCredentials(cfg.Selectel.AccessKey, cfg.Selectel.SecretKey, ""),
		Endpoint:    aws.String(cfg.Selectel.Endpoint),
	})
	if err != nil {
		log.Fatalf("Unable to create session, %v", err)
	}
	return &SelectelS3{session: sess, bucketName: cfg.Selectel.BucketName}
}
func (s3 *SelectelS3) SendImage(img *multipart.FileHeader) (string, error) {
	namePerm := strings.Split(img.Filename, ".")[1]
	file, err := img.Open()
	if err != nil {
		return "", err
	}
	UUID := uuid.New().String()
	objectKey := "photos" + "/" + UUID + "." + namePerm
	uploader := s3manager.NewUploader(s3.session)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s3.bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		return "", err
	}
	selectelPath := "https://8b3c3d07-9708-4c63-93a0-758ae7a7ee89.selstorage.ru" + "/" + objectKey

	return selectelPath, nil

}
