package cli_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/lebenasa/space/cli"
	"github.com/lebenasa/space/service"
)

func TestListBucket(t *testing.T) {
	argv := []string{
		"cli", "list",
	}

	err := cli.Run(argv)
	if err != nil {
		t.Errorf("case 1 got error %v", err)
	}

	spaceKey := service.SpaceKey
	service.SpaceKey = "test"

	err = cli.Run(argv)
	service.SpaceKey = spaceKey
	if err == nil {
		t.Error("case 2 got no error, want error")
	}
}

func TestListObjects(t *testing.T) {
	devBucket, err := service.GetBucket("dev")
	if err != nil {
		t.Error(err)
	}

	argv := []string{
		"cli", "list", devBucket,
	}

	err = cli.Run(argv)
	if err != nil {
		t.Errorf("case 1 got error %v", err)
	}

	argv = []string{
		"cli", "list", fmt.Sprintf("%v/test", devBucket),
	}

	err = cli.Run(argv)
	if err != nil {
		t.Errorf("case 1 got error %v", err)
	}
}

func setupPushFolder(t *testing.T) (path string, objectNames []string) {
	path = "./cli"
	prefix := "test"

	argv := []string{
		"cli", "push", "-r", "--prefix", prefix, path,
	}
	err := cli.Run(argv)
	if err != nil {
		t.Errorf("push folder setup got error %v", err)
	}

	filepath.Walk(path, func(fpath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		relpath, err := filepath.Rel(path, fpath)
		if err != nil {
			return err
		}
		objectName := fmt.Sprintf("%v/%v", prefix, relpath)
		objectNames = append(objectNames, objectName)
		return nil
	})
	fmt.Println(objectNames)

	return
}

func teardownPushFolder(t *testing.T, objectNames []string) {
	argv := []string{
		"cli", "remove",
	}
	argv = append(argv, objectNames...)
	err := cli.Run(argv)
	if err != nil {
		t.Errorf("push folder teardown got error %v", err)
	}
}

func TestPushFolder(t *testing.T) {
	_, objectNames := setupPushFolder(t)
	teardownPushFolder(t, objectNames)
}
