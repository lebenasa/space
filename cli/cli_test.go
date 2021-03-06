package cli_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/lebenasa/space/cli"
	"github.com/lebenasa/space/service"
)

func TestListBucket(t *testing.T) {
	argv := []string{
		"cli", "list-internal",
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
		"cli", "list-internal", devBucket,
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

func TestList(t *testing.T) {
	argv := []string{
		"cli", "list",
	}

	err := cli.Run(argv)
	if err != nil {
		t.Errorf("case 1 got error %v", err)
	}
}

func setupPushFolder(t *testing.T) {
	path := "."
	prefix := "test/cli"

	argv := []string{
		"cli", "push", "-r", "--prefix", prefix, path,
	}
	err := cli.Run(argv)
	if err != nil {
		t.Errorf("push folder setup got error %v", err)
	}

	err = cli.Run([]string{"cli", "list"})
	if err != nil {
		t.Errorf("case 1 got error %v", err)
	}

	return
}

func teardownPushFolder(t *testing.T) {
	argv := []string{
		"cli", "remove",
		"test/cli/cli.go",
		"test/cli/cli_test.go",
		"test/cli/main/main.go",
	}
	err := cli.Run(argv)
	if err != nil {
		t.Errorf("push folder teardown got error %v", err)
	}
}

func setupDownload(t *testing.T) {
	argv := []string{
		"cli", "pull",
		"-o", "./tmp/cli.go",
		"test/cli/cli.go",
	}
	err := cli.Run(argv)
	if err != nil {
		t.Errorf("download setup got error %v", err)
	}
}

func teardownDownload(t *testing.T) {
	err := os.RemoveAll("./tmp")
	if err != nil {
		t.Errorf("download teardown got error %v", err)
	}
}

func TestPushAndDownloadAndRemoveFolder(t *testing.T) {
	setupPushFolder(t)
	setupDownload(t)
	teardownDownload(t)
	teardownPushFolder(t)
}
