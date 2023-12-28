package s3

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"net/http"
)

type FileStore struct {
	s3  *s3.S3
	cfg Config
}

type Config struct {
	AWSSession *session.Session
	Bucket     string
}

func New(cfg Config) *FileStore {
	return &FileStore{
		s3:  s3.New(cfg.AWSSession),
		cfg: cfg,
	}
}

func (f *FileStore) PutObject(object []byte, key string) error {
	_, err := f.s3.PutObject(&s3.PutObjectInput{
		Bucket:               aws.String(f.cfg.Bucket),
		Key:                  aws.String(key),
		ACL:                  aws.String("private"),
		Body:                 bytes.NewReader(object),
		ContentLength:        aws.Int64(int64(len(object))),
		ContentType:          aws.String(http.DetectContentType(object)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
	})
	return err
}

func (f *FileStore) GetObject(key string) ([]byte, error) {
	object, err := f.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(f.cfg.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer object.Body.Close()
	return io.ReadAll(object.Body)
}

func (f *FileStore) DeleteObject(key string) error {
	_, err := f.s3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(f.cfg.Bucket),
		Key:    aws.String(key),
	})
	return err
}
