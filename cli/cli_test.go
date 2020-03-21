package cli_test

import (
	"fmt"
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
		t.Errorf("setup push folder got error %v", err)
	}

	filepath.Walk(path, func(fpath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		relpath = filepath.Rel(path, fpath)
		objectName = fmt.Sprintf("%v/%v", prefix, fpath)
		objectNames = append(objectNames, objectName)
		return nil
	})
	fmt.Println(objectNames)

	return
}

func teardownPushFolder(objectNames []string) error {
	return nil
}

func TestPushFolder(t *testing.T) {
	path, objectNames := setupPushFolder(t)
	teardownPushFolder(objectNames)
}
