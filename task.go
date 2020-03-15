package space

// Collection of routine tasks like uploading a folder.
// Suitable for CLI.

import (
	"context"
	"os"
	"path"
	"path/filepath"

	"github.com/lebenasa/space/service"
)

// WithTags that will be set to all files uploaded with `Upload*` functions.
func (s Space) WithTags(tags map[string]string) Space {
	s.tags = tags
	return s
}

// UploadFile into Space. For large file (>100 MB) please use `UploadBigFile`.
// If Space is created using `WithTags`, apply those tags into uploaded file.
func (s Space) UploadFile(ctx context.Context, fp, env, prefix string) (objectName string, err error) {
	bucket, err := service.GetBucket(env)
	if err != nil {
		return
	}
	filename := filepath.Base(fp)
	objectName = path.Join(prefix, filename)

	_, err = s.PutFile(ctx, bucket, objectName, fp, PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return
	}

	if len(s.tags) == 0 {
		return
	}
	err = s.PutTag(ctx, bucket, objectName, s.tags)
	return
}

// UploadFolder into Space. Do not use if there's a large file (>100 MB) inside the folder.
func (s Space) UploadFolder(ctx context.Context, fp, env, prefix string) (objectNames []string, err error) {
	filePaths := []string{}
	filepath.Walk(fp, func(fpath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		// TODO: skips ignored files
		filePaths = append(filePaths, fpath)

		return nil
	})

	// TODO: do this concurrently
	for _, filePath := range filePaths {
		objectName, errr := s.UploadFile(ctx, filePath, env, prefix)
		if errr != nil {
			return objectNames, errr
		}
		objectNames = append(objectNames, objectName)
	}

	return
}
