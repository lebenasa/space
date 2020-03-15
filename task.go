package space

import (
	".space/service"
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"time"
)

// Collection of routine tasks like uploading a folder.
// Suitable for CLI.

// Environments in Space to work with. Maps to bucket name.
var environments = map[string]string{
	"dev":     "dev",
	"staging": "dev",
	"live":    "dev",
}

// WithTags that will be set to all files uploaded with `Upload*` functions.
func (s Space) WithTags(tags map[string]string) (s Space) {
	s.tags = tags
	return s
}

// GetBucket name from given environment name.
func GetBucket(env string) (string, error) {
	envs := make([]string, len(environments))
	for key, _ := range environments {
		envs = append(envs, key)
		if key == env {
			return env, nil
		}
	}
	return "", fmt.Errorf("Invalid environment %v, possible values: %v", env, envs)
}

// UploadFile into Space. For large file (>100 MB) please use `UploadBigFile`.
// If Space is created using `WithTags`, apply those tags into uploaded file.
func (s Space) UploadFile(ctx context.Context, fp, env, prefix string) (objectName string, err error) {
	bucket, err := GetBucket(env)
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
	bucket, err := GetBucket(env)
	if err != nil {
		return
	}

	filePaths := make([]string)
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
		objectName, err := s.UploadFileWithContext(ctx, filePath, env, prefix)
		if err != nil {
			return
		}
		objectNames = append(objectNames, objectName)
	}

	return
}
