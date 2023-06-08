package unupdatednixpkgs

import (
	"fmt"
	"os"
	"path/filepath"
)

func findUnupdatedPkgs(path string) ([]string, error) {
	fmt.Println("Path : " + path)
	findNixPkgs, err := findNixPkgs(path)
	if err != nil {
		return nil, err
	}

	return findNixPkgs, nil
}

func findNixPkgs(path string) ([]string, error) {
	var result []string
	err := filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() == "default.nix" {
			result = append(result, filePath)
		}
		return nil
	})

	return result, err
}
