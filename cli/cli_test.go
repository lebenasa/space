package cli_test

import (
	"testing"

	"github.com/lebenasa/space/cli"
	"github.com/lebenasa/space/service"
)

func TestListBucket(t *testing.T) {
	argvFn := func() []string {
		return []string{
			"cli", "list",
		}
	}

	err := cli.Run(argvFn)
	if err != nil {
		t.Errorf("case 1 got error %v", err)
	}

	spaceKey := service.SpaceEndpoint
	service.SpaceKey = "test"

	err = cli.Run(argv)
	service.SpaceKey = spaceKey
	if err == nil {
		t.Error("case 2 got no error, want error")
	}
}

func setupPushFolder(t *testing.T) (path string, objectNames []string) {
	path = "./cli"
	prefix := "test"

	argv := []string{
		"cli", "push", "-r", "--prefix", prefix, path,
	}
	argvFn := func() []string {
		return argv
	}
}

func TestListObjects(t *testing.T) {
	devBucket := service.GetBucket("dev")
	argv := []string{
		"cli", "list", devBucket,
	}
	argvFn := func() []string {
		return argv
	}

	err := cli.Run(argvFn)
	if err != nil {
		t.Errorf("case 1 got error %v", err)
	}

	argv = []string{
		"cli", "list", fmt.Sprintf("%v/test", devBucket),
	}

	err := cli.Run(argvFn)
	if err != nil {
		t.Errorf("case 1 got error %v", err)
	}
}
