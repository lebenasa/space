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

func setupPushFolder(t *testing.T) (path string, objectNames []string) {
	// path = "./cli"
	// prefix := "test"

	// argv := []string{
	// 	"cli", "push", "-r", "--prefix", prefix, path,
	// }
	return
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
