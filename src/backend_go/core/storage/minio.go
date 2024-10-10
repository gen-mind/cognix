package storage

import (
	"cognix.ch/api/v2/core/utils"
	"context"
	minio "github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
)

// Package containing MinIO related code.
// MinioConfig represents the configuration settings for connecting to a MinIO server.
type (
	MinioConfig struct {
		AccessKey       string `env:"MINIO_ACCESS_KEY"`
		SecretAccessKey string `env:"MINIO_SECRET_ACCESS_KEY"`
		Endpoint        string `env:"MINIO_ENDPOINT"`
		UseSSL          bool   `env:"MINIO_USE_SSL"`
		BucketName      string `env:"MINIO_BUCKET_NAME"`
		Region          string `env:"MINIO_REGION"`
	}
	FileStorageClient interface {
		Upload(ctx context.Context, bucket, filename, contentType string, reader io.Reader) (string, string, error)
		GetObject(ctx context.Context, bucket, filename string, writer io.Writer) error
		DeleteObject(ctx context.Context, bucket, filename string) error
	}
	minIOClient struct {
		Region string
		client *minio.Client
	}
	minIOMockClient struct{}
)

// DeleteObject removes an object from the specified bucket with the given filename.
// It utilizes the RemoveObject method of the underlying minio.Client to perform the deletion.
// The ForceDelete option is set to true to bypass the object lock and delete it forcefully.
// It returns an error if the deletion fails.
func (c *minIOClient) DeleteObject(ctx context.Context, bucket, filename string) error {
	return c.client.RemoveObject(ctx, bucket, filename, minio.RemoveObjectOptions{
		ForceDelete: true,
	})
}

// checkOrCreateBucket checks if a bucket exists and creates it if it does not exist
// Parameters:
// - ctx: the context.Context object for the request
// - bucketName: the name of the bucket to check or create
// Returns:
// - error: an error object if any error occurs, otherwise nil
func (c *minIOClient) checkOrCreateBucket(ctx context.Context, bucketName string) error {
	ok, err := c.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}
	// create bucket if not exists
	if !ok {
		return c.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			Region: c.Region,
		})
	}
	return nil
}

// Upload uploads a file to the specified bucket in the minio storage.
//
// Parameters:
//   - ctx: The context.Context object.
//   - bucket: The name of the bucket where the file will be uploaded.
//   - filename: The name of the file to be uploaded.
//   - contentType: The content type of the file.
//   - reader: The io.Reader object representing the file content.
//
// Returns:
//   - string: The key of the uploaded file in the minio storage.
//   - string: The CRC32C checksum of the uploaded file.
//   - error: The error that occurred during the upload process, if any.
//
// The method first checks if the specified bucket exists and creates it if it doesn't.
// Then it uploads the file to the minio storage using the provided reader and options.
// The key and checksum of the uploaded file are returned upon successful upload.
// Any error that occurred during the upload process is returned.
func (c *minIOClient) Upload(ctx context.Context, bucket, filename, contentType string, reader io.Reader) (string, string, error) {
	// verify is bucket exists. create if not exists
	if err := c.checkOrCreateBucket(ctx, bucket); err != nil {
		return "", "", err
	}

	// save file in minio
	res, err := c.client.PutObject(ctx, bucket, filename, reader, -1,
		minio.PutObjectOptions{
			ContentType: contentType,
			NumThreads:  0,
		})
	if err != nil {
		return "", "", utils.Internal.Wrapf(err, "cannot upload file: %s", err.Error())
	}
	return res.Key, res.ChecksumCRC32C, nil
}

// GetObject retrieves an object from the specified bucket and saves it to the given writer.
// It returns an error if there was a problem retrieving or saving the object.
func (c *minIOClient) GetObject(ctx context.Context, bucket, filename string, writer io.Writer) error {
	object, err := c.client.GetObject(ctx, bucket, filename, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer object.Close()
	_, err = io.Copy(writer, object)
	if err != nil {
		return err
	}
	return nil
}

// NewMinIOClient creates a new instance of FileStorageClient using the provided MinioConfig.
// It returns the created FileStorageClient and an error if any occurred.
func NewMinIOClient(cfg *MinioConfig) (FileStorageClient, error) {

	minioClient, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(cfg.AccessKey,
			cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	return &minIOClient{
		Region: cfg.Region,
		client: minioClient,
	}, nil
}
