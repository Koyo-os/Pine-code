package main

import (
	"fmt"
	"io"
	"os"
)

const templ = `
type: "piner"
addr: "%s"

exclude_files :
  - ""
exclude_dirs :
  - "test"
  - "cmd"
  - "build"
include_dirs :
  - "all"
temp_dir : "temp"
	`

func initPiner(addr string) error {
	file, err := os.Create("pine.yaml")
	if err != nil {
		return err
	}

	text := fmt.Sprintf(templ, addr)
	_, err = io.WriteString(file, text)
	if err != nil {
		return err
	}

	return nil
}
