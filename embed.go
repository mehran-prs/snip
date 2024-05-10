package main

import (
	"embed"
	"io"
	"os"
	"path/filepath"
)

//go:embed autocompletes
var autocompletes embed.FS

func WriteAutocompleteScript(w io.Writer, fname string) error {
	bytes, err := autocompletes.ReadFile(filepath.Join("autocompletes", fname))
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(bytes)
	return err
}
