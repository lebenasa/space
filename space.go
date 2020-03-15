package space

import (
	"context"
	"fmt"
	"io"

	"github.com/lebenasa/space/service"
	"github.com/minio/minio-go/v6"
)

// Space access client to limit what can be done programatically to our Spaces.
type Space struct {
	client *minio.Client
	tags   map[string]string
}

// Object represents an open object.
type Object = minio.Object

// BucketInfo contains bucket's metadata.
type BucketInfo = minio.BucketInfo

// ObjectInfo contains object's metadata.
type ObjectInfo = minio.ObjectInfo

// PutObjectOptions specifies additional headers when putting object to Space.
type PutObjectOptions = minio.PutObjectOptions

// GetObjectOptions specifies additional headers when getting object from Space.
type GetObjectOptions = minio.GetObjectOptions

// StatObjectOptions specifies additional headers when stating object in Space.
type StatObjectOptions = minio.StatObjectOptions

// New space client.
// Requires generated `service` module that's not tracked by git.
func New() (space Space, err error) {
	client, err := minio.New(service.SPACE_ENDPOINT, service.SPACE_KEY, service.SPACE_SECRET, true)
	if err != nil {
		return space, err
	}

	space.client = client
	return
}

// NewFromClient via `minio.New`.
func NewFromClient(client *minio.Client) (space Space) {
	space.client = client
	return
}

// SetAppInfo adds custom application details to User-Agent.
func (s Space) SetAppInfo(appName, appVersion string) {
	s.client.SetAppInfo(appName, appVersion)
}

// ListBuckets in current endpoint.
func (s Space) ListBuckets() ([]BucketInfo, error) {
	return s.client.ListBuckets()
}

// ListObjects inside a bucket.
func (s Space) ListObjects(bucketName string, objectPrefix string, recursive bool) (objects []ObjectInfo, err error) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	objectCh := s.client.ListObjectsV2(bucketName, objectPrefix, recursive, doneCh)
	for object := range objectCh {
		if object.Err != nil {
			return nil, object.Err
		}
		objects = append(objects, object)
	}

	return objects, err
}

// Put object to Space.
func (s Space) Put(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64, options PutObjectOptions) (int64, error) {
	return s.client.PutObjectWithContext(ctx, bucketName, objectName, reader, objectSize, options)
}

// Get object from Space.
func (s Space) Get(ctx context.Context, bucketName, objectName string, options GetObjectOptions) (*Object, error) {
	return s.client.GetObjectWithContext(ctx, bucketName, objectName, options)
}

// PutFile to Space (upload a file).
func (s Space) PutFile(ctx context.Context, bucketName, objectName, filePath string, options PutObjectOptions) (length int64, err error) {
	return s.client.FPutObjectWithContext(ctx, bucketName, objectName, filePath, options)
}

// GetFile from Space (download a file).
func (s Space) GetFile(ctx context.Context, bucketName, objectName, filePath string, options GetObjectOptions) error {
	return s.client.FGetObjectWithContext(ctx, bucketName, objectName, filePath, options)
}

// Stat of an object in Space.
func (s Space) Stat(bucketName, objectName string, options StatObjectOptions) (ObjectInfo, error) {
	return s.client.StatObject(bucketName, objectName, options)
}

// Remove object in Space.
func (s Space) Remove(bucketName, objectName string) error {
	return s.client.RemoveObject(bucketName, objectName)
}

// RemoveObjects in Space.
func (s Space) RemoveObjects(ctx context.Context, bucketName string, objectNames []string) (err error) {
	objectsCh := make(chan string)

	go func() {
		defer close(objectsCh)
		for _, name := range objectNames {
			objectsCh <- name
		}
	}()

	for rErr := range s.client.RemoveObjectsWithContext(ctx, bucketName, objectsCh) {
		err = fmt.Errorf("%v\nFailed to remove %v: %v", err, rErr.ObjectName, rErr)
	}

	return err
}

// PutTag on an object in Space.
func (s Space) PutTag(ctx context.Context, bucketName, objectName string, tags map[string]string) error {
	return s.client.PutObjectTaggingWithContext(ctx, bucketName, objectName, tags)
}

// GetTag of an object in Space. Returned string is in XML format.
func (s Space) GetTag(ctx context.Context, bucketName, objectName string) (string, error) {
	return s.client.GetObjectTaggingWithContext(ctx, bucketName, objectName)
}

// RemoveTag from an object in Space.
func (s Space) RemoveTag(ctx context.Context, bucketName, objectName string) error {
	return s.client.RemoveObjectTaggingWithContext(ctx, bucketName, objectName)
}
