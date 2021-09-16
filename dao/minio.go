package dao

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"path/filepath"
	"time"

	minio "github.com/minio/minio-go"
)

type MinioStorage struct {
	Client *minio.Client
	Bucket string
}

func NewMinioStorage(client *minio.Client, bket string) *MinioStorage {
	return &MinioStorage{
		Client: client,
		Bucket: bket,
	}
}

func (m *MinioStorage) Get(ctx context.Context, path string) ([]byte, error) {
	object, err := m.Client.GetObjectWithContext(ctx, m.Bucket, path, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	var buff bytes.Buffer
	writer := bufio.NewWriter(&buff)
	_, err = io.Copy(writer, object)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func (m *MinioStorage) GetURL(ctx context.Context, path string, timeout time.Duration) (string, error) {
	name := filepath.Base(path)
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", fmt.Sprintf(`attachment; filename="%s"`, name))

	presignedUrl, err := m.Client.PresignedGetObject(m.Bucket, path, timeout, reqParams)
	if err != nil {
		return "", err
	}

	return presignedUrl.String(), nil
}

func (m *MinioStorage) Save(ctx context.Context, data []byte, path, ext string) error {
	// return ioutil.WriteFile(path, data, 0644)
	return putImageToMinio(ctx, m.Client, m.Bucket, data, path, ext)
}

func (m *MinioStorage) Move(ctx context.Context, from, path string) error {
	srcOpts := minio.NewSourceInfo(m.Bucket, from, nil)

	// Destination object
	dstOpts, err := minio.NewDestinationInfo(m.Bucket, path, nil, nil)
	if err != nil {
		return err
	}

	err = m.Client.CopyObject(dstOpts, srcOpts)
	if err != nil {
		return err
	}
	err = m.Client.RemoveObject(m.Bucket, from)
	if err != nil {
		return err
	}
	return nil
}

func putImageToMinio(ctx context.Context, minioClient *minio.Client, bucket string, d []byte, objectName, ext string) error {
	// Upload the zip file

	if ok, err := minioClient.BucketExists(bucket); !ok {
		return fmt.Errorf("Bucket %s not exists %v \n", bucket, err)
	}

	reader := bytes.NewReader(d)
	_, err := minioClient.PutObject(bucket, objectName, reader, reader.Size(), minio.PutObjectOptions{ContentType: fmt.Sprintf("image/%s", ext)})
	if err != nil {
		return err
	}
	return nil
}
