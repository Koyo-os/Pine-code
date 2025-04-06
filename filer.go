package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func compareFiles(file1, file2 string) (bool, error) {
	content1, err := os.ReadFile(file1)
	if err != nil {
		return false, err
	}

	content2, err := os.ReadFile(file2)
	if err != nil {
		return false, err
	}

	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(string(content1), string(content2), false)

	if len(diffs) == 1 && diffs[0].Type == diffmatchpatch.DiffEqual {
		return true, nil
	}
	return false, nil
}

func compareDirs(dir1, dir2 string) (bool, error) {
	filesInDir1 := make(map[string]struct{})
	err := filepath.Walk(dir1, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(dir1, path)
			filesInDir1[relPath] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return false, err
	}

	err = filepath.Walk(dir2, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, _ := filepath.Rel(dir2, path)

			if _, exists := filesInDir1[relPath]; !exists {
				return fmt.Errorf("файл %s отсутствует в первой директории", relPath)
			}

			file1 := filepath.Join(dir1, relPath)
			file2 := filepath.Join(dir2, relPath)
			ok, err := compareFiles(file1, file2)
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("файлы %s различаются", relPath)
			}

			delete(filesInDir1, relPath)
		}
		return nil
	})
	if err != nil {
		return false, err
	}

	if len(filesInDir1) > 0 {
		return false, fmt.Errorf("в первой директории есть лишние файлы: %v", filesInDir1)
	}

	return true, nil
}
