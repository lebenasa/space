package cli

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lebenasa/space"

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

	fmt.Println("Object\t\t\tSize\t\tLast modified")
	for _, object := range objects {
		fmt.Printf("%v\t\t%v\t\t%v\n", object.Key, object.Size, object.LastModified)
	}
	return nil
}

func listBuckets(s space.Space) error {
	fmt.Println("Listing all buckets")
	buckets, err := s.ListBuckets()
	if err != nil {
		return err
	}

	fmt.Println("Bucket\t\t\tCreated on")
	for _, bucket := range buckets {
		fmt.Printf("%v\t\t\t%v\n", bucket.Name, bucket.CreationDate)
	}

	return nil
}

func listAction(c *cli.Context) error {
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

func pushFolder(folder string, s space.Space, env string, prefix string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*60*time.Second)
	defer cancel()

	// TODO: verify uploaded files
	fmt.Printf("Uploading %v\n", folder)
	objectNames, err := s.UploadFolder(ctx, folder, env, prefix)
	if err != nil {
		return err
	}

	for _, objectName := range objectNames {
		fmt.Printf("Uploaded %v\n", objectName)
	}

	return nil
}

func pushFile(fileName string, s space.Space, env string, prefix string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*60*time.Second)
	defer cancel()

	// TODO: verify uploaded file
	fmt.Printf("Uploading %v\n", fileName)
	objectName, err := s.UploadFile(ctx, fileName, env, prefix)
	if err != nil {
		return err
	}

	fmt.Printf("Uploaded to %v\n", objectName)
	return nil
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

// Run using arguments from `argv` provider function.
func Run(argv []string) (err error) {
	envFlag := cli.StringFlag{
		Name:  "env",
		Value: "dev",
		Usage: "Specify Space environment",
	}

	listCommand := cli.Command{
		Name:      "list",
		Usage:     "List available buckets or objects in Space. Not a good idea for production bucket.",
		ArgsUsage: "If given, list all objects in {bucket}/{prefix}, otherwise list all buckets",
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

	app := &cli.App{
		Name:  "space",
		Usage: "Work with Space and assets",
		Commands: []*cli.Command{
			&listCommand,
			&pushCommand,
		},
	}

	err = app.Run(argv)
	return err
}
