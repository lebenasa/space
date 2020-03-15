package space_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/lebenasa/space"
	"github.com/lebenasa/space/service"
	"github.com/minio/minio-go/v6"
)

func TestNew(t *testing.T) {
	_, err := space.New()
	if err != nil {
		t.Errorf("case 1 got %v, want nil error", err)
	}

	endpoint := service.SPACE_ENDPOINT
	service.SPACE_ENDPOINT = "https://foo.bar.com"
	_, err = space.New()
	service.SPACE_ENDPOINT = endpoint
	if err == nil {
		t.Errorf("case 2 got %v, want error", err)
	}
}

func TestNewFromClient(t *testing.T) {
	client, _ := minio.New(service.SPACE_ENDPOINT, service.SPACE_KEY, service.SPACE_SECRET, true)
	s := space.NewFromClient(client)
	s.SetAppInfo("test", "0.0.0")
	if _, err := s.ListBuckets(); err != nil {
		t.Errorf("case 1 got %v", err)
	}
}

func setupSpace(t *testing.T) (space.Space, string) {
	s, err := space.New()
	if err != nil {
		t.Errorf("setup space fail: %v", err)
	}
	bucket, err := service.GetBucket("dev")
	if err != nil {
		t.Errorf("setup bucket fail: %v", err)
	}
	return s, bucket
}

func TestListBuckets(t *testing.T) {
	s, _ := setupSpace(t)
	_, err := s.ListBuckets()
	if err != nil {
		t.Errorf("case 1 got %v", err)
	}
}

func TestListObjects(t *testing.T) {
	s, bucket := setupSpace(t)

	objectName := "test/space.go"
	fileName := "./space.go"
	_, err := s.PutFile(context.Background(), bucket, objectName, fileName, space.PutObjectOptions{})
	if err != nil {
		t.Errorf("setup put file fail: %v", err)
	}

	objectNames, err := s.ListObjects(bucket, "test", true)
	if err != nil {
		t.Errorf("case 1 got %v", err)
	}
	found := false
	for _, name := range objectNames {
		if name.Key == objectName {
			found = true
		}
	}
	if !found {
		t.Errorf("case 1 object found is %v,, want true", found)
	}

	err = s.Remove(bucket, objectName)
	if err != nil {
		t.Errorf("teardown fail: %v", err)
	}
}

func setupPut(objectName, objectContent string, s space.Space, bucket string) error {
	content := strings.NewReader(objectContent)
	length, err := s.Put(context.Background(), bucket, objectName, content, content.Size(), space.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("put setup got %v", err)
	}
	if length != content.Size() {
		return fmt.Errorf("put setup got length %v, want %v", length, content.Size())
	}
	return nil
}

func teardownPut(objectName string, s space.Space, bucket string) error {
	err := s.Remove(bucket, objectName)
	if err != nil {
		return fmt.Errorf("put teardown got %v", err)
	}
	return nil
}

func TestPut(t *testing.T) {
	s, bucket := setupSpace(t)
	objectName := "test/put.txt"
	objectContent := "test content"
	err := setupPut(objectName, objectContent, s, bucket)
	if err != nil {
		t.Error(err)
	}
	err = teardownPut(objectName, s, bucket)
	if err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	s, bucket := setupSpace(t)
	objectName := "test/get.txt"
	objectContent := "test content"
	err := setupPut(objectName, objectContent, s, bucket)
	if err != nil {
		t.Error(err)
	}

	object, err := s.Get(context.Background(), bucket, objectName, space.GetObjectOptions{})
	if err != nil {
		t.Errorf("case 1 got erorr %v", err)
	}

	buf := make([]byte, len(objectContent))
	if _, err = io.ReadFull(object, buf); err != nil {
		t.Errorf("case 1 got error %v", err)
	}
	bufstr := string(buf)
	if bufstr != objectContent {
		t.Errorf("case 1 read %v, want %v", bufstr, objectContent)
	}

	err = teardownPut(objectName, s, bucket)
	if err != nil {
		t.Error(err)
	}
	err = os.RemoveAll("./tmp")
	if err != nil {
		t.Error(err)
	}
}

func setupPutFile(objectName, filePath string, s space.Space, bucket string) error {
	_, err := s.PutFile(context.Background(), bucket, objectName, filePath, space.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("setup put file got error %v", err)
	}
	return nil
}

func TestPutFile(t *testing.T) {
	s, bucket := setupSpace(t)
	objectName := "test/space.go"
	filePath := "./space.go"
	err := setupPutFile(objectName, filePath, s, bucket)
	if err != nil {
		t.Error(err)
	}

	err = teardownPut(objectName, s, bucket)
	if err != nil {
		t.Error(err)
	}
}

func TestGetFile(t *testing.T) {
	s, bucket := setupSpace(t)
	objectName := "test/space.go"
	filePath := "./space.go"
	outPath := "./tmp/space.go"
	err := setupPutFile(objectName, filePath, s, bucket)
	if err != nil {
		t.Error(err)
	}

	err = s.GetFile(context.Background(), bucket, objectName, outPath, space.GetObjectOptions{})
	if err != nil {
		t.Errorf("case 1 got error %v", err)
	}

	length := 256
	original := make([]byte, length)
	remote := make([]byte, length)
	originalFile, err := os.Open(filePath)
	if err != nil {
		t.Errorf("case 1 original file: %v", err)
	}
	if _, err := io.ReadFull(originalFile, original); err != nil {
		t.Errorf("case 1 original file: %v", err)
	}
	remoteFile, err := os.Open(filePath)
	if err != nil {
		t.Errorf("case 1 remote file: %v", err)
	}
	if _, err := io.ReadFull(remoteFile, remote); err != nil {
		t.Errorf("case 1 remote file: %v", err)
	}

	ostr := string(original)
	rstr := string(remote)
	if rstr != ostr {
		t.Errorf("case 1 got %v, want %v", rstr, ostr)
	}

	err = teardownPut(objectName, s, bucket)
	if err != nil {
		t.Error(err)
	}
}

func TestStat(t *testing.T) {
	s, bucket := setupSpace(t)
	objectName := "test/stat.txt"
	objectContent := "test content"
	err := setupPut(objectName, objectContent, s, bucket)
	if err != nil {
		t.Error(err)
	}

	info, err := s.Stat(bucket, objectName, space.StatObjectOptions{})
	if err != nil {
		t.Errorf("case 1 got error %v", err)
	}

	if info.Key != objectName {
		t.Errorf("case 1 got name %v, want %v", info.Key, objectName)
	}

	err = teardownPut(objectName, s, bucket)
	if err != nil {
		t.Error(err)
	}
}

func TestRemove(t *testing.T) {
	TestPut(t)
}

func TestRemoveObjects(t *testing.T) {
	s, bucket := setupSpace(t)
	objectNames := []string{
		"test/object1.txt", "test/foo/object2.txt",
	}
	objectContent := "test content"

	for _, objectName := range objectNames {
		if err := setupPut(objectName, objectContent, s, bucket); err != nil {
			t.Error(err)
		}
	}

	err := s.RemoveObjects(context.Background(), bucket, objectNames)
	if err != nil {
		t.Errorf("case 1 got error %v", err)
	}
}

func TestTag(t *testing.T) {
	s, bucket := setupSpace(t)
	objectName := "test/tag.txt"
	objectContent := "test content"

	err := setupPut(objectName, objectContent, s, bucket)
	if err != nil {
		t.Error(err)
	}

	tags := map[string]string{
		"foo": "bar",
	}
	err = s.PutTag(context.Background(), bucket, objectName, tags)
	if err != nil {
		t.Errorf("case 1 got error %v", err)
	}
	_, err = s.GetTag(context.Background(), bucket, objectName)
	if err != nil {
		t.Errorf("case 2 got error %v", err)
	}
	s.RemoveTag(context.Background(), bucket, objectName)

	err = teardownPut(objectName, s, bucket)
	if err != nil {
		t.Error(err)
	}
}
