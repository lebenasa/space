package space_test

import (
	"context"
	"testing"

	"github.com/lebenasa/space"
	"github.com/lebenasa/space/service"
	"github.com/minio/minio-go/v6"
)

func TestNew(t *testing.T) {
	_, err := space.New()
	if err != nil {
		t.Errorf("case 1 %v want nil error", err)
	}

	endpoint := service.SPACE_ENDPOINT
	service.SPACE_ENDPOINT = "https://foo.bar.com"
	_, err = space.New()
	service.SPACE_ENDPOINT = endpoint
	if err == nil {
		t.Errorf("case 2 %v want error", err)
	}
}

func TestNewFromClient(t *testing.T) {
	client, _ := minio.New(service.SPACE_ENDPOINT, service.SPACE_KEY, service.SPACE_SECRET, true)
	s := space.NewFromClient(client)
	s.SetAppInfo("test", "0.0.0")
	if _, err := s.ListBuckets(); err != nil {
		t.Errorf("case 1 %v want no error", err)
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
		t.Errorf("case 1 %v want no error", err)
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
		t.Errorf("case 1 %v want no error", err)
	}
	found := false
	for _, name := range objectNames {
		if name.Key == objectName {
			found = true
		}
	}
	if !found {
		t.Errorf("case 1 object found is %v, want true", found)
	}

	err = s.Remove(bucket, objectName)
	if err != nil {
		t.Errorf("teardown fail: %v", err)
	}
}
