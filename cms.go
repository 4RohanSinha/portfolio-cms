package main

import (
	"cms/firebase_service"
	"cms/version_control"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

func main() {

	var newVersionId string
	var newVersionDescr string

	if len(os.Args) < 2 {
		panic("No command provided")
	}

	command := os.Args[1]
	args := os.Args[2:]

	s_meta := version_control.Metadata{
		Timestamp:   time.Now().String(),
		Description: "initial snapshot",
	}

	switch command {
	case "restore":
		var fileName string
		cmdFlags := flag.NewFlagSet("restore", flag.ExitOnError)
		cmdFlags.StringVar(&newVersionId, "id", "", "ID of version to restore")
		cmdFlags.Parse(args)

		err := filepath.WalkDir(".vc/versions/", func(path string, info os.DirEntry, err error) error {
			if !info.IsDir() && strings.HasPrefix(info.Name(), newVersionId) {
				fileName = info.Name()
				return nil
			}

			if err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			panic(err)
		}

		version_control.Restore(fileName)
	case "save":
		cmdFlags := flag.NewFlagSet("save", flag.ExitOnError)
		cmdFlags.StringVar(&newVersionDescr, "descr", "", "Description of changes")
		cmdFlags.Parse(args)

		v_meta := version_control.Metadata{
			Timestamp:   time.Now().String(),
			Description: newVersionDescr,
		}

		head, err := version_control.Head()

		if err != nil {
			s := version_control.Snapshot{
				Mdata:   s_meta,
				Muuid:   uuid.NewString(),
				Content: make(map[string]string),
			}

			version_control.DumpVersion(s)
			//version_control.Restore(s_meta.Title)
			head = s
		}

		delta, err := version_control.CalculateDelta(head, v_meta, uuid.NewString())

		if err != nil {
			panic(err)
		}

		version_control.DumpVersion(delta)
		version_control.Restore(delta.Muuid)
	case "upload":
		var collection string
		cmdFlags := flag.NewFlagSet("upload", flag.ExitOnError)
		cmdFlags.StringVar(&collection, "mode", "", "[prod/dev]")
		cmdFlags.Parse(args)

		app, err := firebase_service.App()

		if err != nil {
			panic(err)
		}

		if collection != "prod" && collection != "dev" {
			panic("Invalid mode picked.")
		}

		err = firebase_service.SetHead(app, "posts-"+collection)

		if err != nil {
			panic(err)
		}
	default:
		panic("Unknown command")
	}
}
