package cli

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lebenasa/space"

	"github.com/jedib0t/go-pretty/table"
	"github.com/urfave/cli/v2"
)

func handleEnum(val string, enums []string) (value string, err error) {
	for _, enum := range enums {
		if val == enum {
			return val, err
		}
	}
	return "", fmt.Errorf("Invalid argument %v, possible values: %v", val, enums)
}

func handleEnvFlag(val string) (string, error) {
	return handleEnum(val, []string{
		"dev",
		"live",
	})
}

func listObjects(s space.Space, bucket, prefix string) error {
	fmt.Printf("Listing objects from %v with prefix '%v'\n", bucket, prefix)
	objects, err := s.ListObjects(bucket, prefix, true)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Object", "Size", "Last modified"})
	for _, object := range objects {
		t.AppendRow([]interface{}{object.Key, object.Size, object.LastModified})
	}
	t.SetStyle(table.StyleColoredBlueWhiteOnBlack)
	t.Render()
	return nil
}

func listBuckets(s space.Space) error {
	buckets, err := s.ListBuckets()
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Bucket", "Created on"})
	for _, bucket := range buckets {
		t.AppendRow([]interface{}{bucket.Name, bucket.CreationDate})
	}
	t.SetStyle(table.StyleColoredBlueWhiteOnBlack)
	t.Render()

	return nil
}

func listInternalAction(c *cli.Context) error {
	bucket, prefix := parseBucketAndPrefix(c.Args().First())

	s, err := space.New()
	if err != nil {
		return err
	}

	if bucket != "" {
		return listObjects(s, bucket, prefix)
	}

	return listBuckets(s)
}

func listAction(c *cli.Context) error {
	env, err := handleEnvFlag(c.String("env"))
	if err != nil {
		return err
	}

	s, err := space.New()
	if err != nil {
		return err
	}

	prefix := c.Args().First()
	objects, err := s.List(env, prefix)
	if err != nil {
		return err
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Object", "Size", "Last modified"})
	for _, object := range objects {
		t.AppendRow([]interface{}{object.Key, object.Size, object.LastModified})
	}
	t.SetStyle(table.StyleColoredBlueWhiteOnBlack)
	t.Render()

	return nil
}

func pushFolder(folder string, s space.Space, env string, prefix string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*60*time.Second)
	defer cancel()

	// TODO: verify uploaded files
	objectNames, err := s.UploadFolder(ctx, folder, env, prefix)
	for _, name := range objectNames {
		fmt.Println(name)
	}
	return err
}

func pushFile(fileName string, s space.Space, env string, prefix string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*60*time.Second)
	defer cancel()

	fi, err := os.Stat(fileName)
	if err != nil {
		return err
	}

	if fi.IsDir() {
		return fmt.Errorf("%v is a directory, push with --recursive flag", fileName)
	}

	// TODO: verify uploaded file
	objectName, err := s.UploadFile(ctx, fileName, env, prefix)
	fmt.Println(objectName)
	return err
}

func pushAction(c *cli.Context) error {
	env, err := handleEnvFlag(c.String("env"))
	if err != nil {
		return err
	}

	s, err := space.New()
	if err != nil {
		return err
	}

	s = s.WithTags(parseTags(c.String("tags")))

	fp := c.Args().Get(0)
	if fp == "" {
		return fmt.Errorf("Invalid file/folder: '%v'", fp)
	}

	prefix := c.String("prefix")
	if c.Bool("recursive") {
		return pushFolder(fp, s, env, prefix)
	}
	return pushFile(fp, s, env, prefix)
}

func parseBucketAndPrefix(text string) (bucket, prefix string) {
	split := strings.SplitN(text, "/", 2)
	if len(split) == 2 {
		bucket = split[0]
		prefix = split[1]
	} else {
		bucket = text
	}
	return bucket, prefix
}

func parseTags(text string) (tags map[string]string) {
	pairs := strings.Split(text, ",")
	if pairs[0] == text {
		return
	}
	kv := func(keyval []string) (string, string) {
		if len(keyval) != 2 {
			return "", ""
		}
		return strings.TrimSpace(keyval[0]), strings.TrimSpace(keyval[1])
	}
	for _, pair := range pairs {
		key, val := kv(strings.Split(pair, ":"))
		if key == "" {
			continue
		}
		tags[key] = val
	}
	return tags
}

func removeAction(c *cli.Context) error {
	env, err := handleEnvFlag(c.String("env"))
	if err != nil {
		return err
	}

	s, err := space.New()
	if err != nil {
		return err
	}

	objectNames := make([]string, c.Args().Len())
	for i := range objectNames {
		objectNames[i] = c.Args().Get(i)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*60*time.Second)
	defer cancel()

	return s.RemoveFiles(ctx, env, objectNames)
}

// Run using arguments from `argv`.
func Run(argv []string) (err error) {
	envFlag := cli.StringFlag{
		Name:  "env",
		Value: "dev",
		Usage: "Specify Space environment",
	}

	listInternalCommand := cli.Command{
		Name:      "list-internal",
		Usage:     "List available buckets or objects in Space. Not a good idea for production bucket.",
		ArgsUsage: "If given, list all objects in {bucket}/{prefix}, otherwise list all buckets",
		HideHelp:  true,
		Flags: []cli.Flag{
			&envFlag,
		},
		Action: listInternalAction,
	}

	listCommand := cli.Command{
		Name:      "list",
		Usage:     "List available objects in Space.",
		ArgsUsage: "Prefix",
		Flags: []cli.Flag{
			&envFlag,
		},
		Action: listAction,
	}

	pushCommand := cli.Command{
		Name:      "push",
		Usage:     "Upload file/folder to Space",
		ArgsUsage: "File or folder path to upload",
		Flags: []cli.Flag{
			&envFlag,
			&cli.BoolFlag{
				Name:    "recursive",
				Aliases: []string{"r"},
				Usage:   "Upload a folder recursively",
				Value:   false,
			},
			&cli.StringFlag{
				Name:    "prefix",
				Aliases: []string{"p"},
				Usage:   "Object name's prefix.",
				Value:   "",
			},
			&cli.StringFlag{
				Name:    "tags",
				Aliases: []string{"t"},
				Usage:   "Add tags, e.g. \"version: 0.0, type: app\"",
				Value:   "",
			},
		},
		Action: pushAction,
	}

	removeCommand := cli.Command{
		Name:      "remove",
		Aliases:   []string{"rm"},
		Usage:     "Remove file(s) in Space",
		ArgsUsage: "Files to be removed",
		Flags: []cli.Flag{
			&envFlag,
		},
		Action: removeAction,
	}

	app := &cli.App{
		Name:  "space",
		Usage: "Work with Space and assets",
		Commands: []*cli.Command{
			&listInternalCommand,
			&listCommand,
			&pushCommand,
			&removeCommand,
		},
	}

	err = app.Run(argv)
	return err
}
