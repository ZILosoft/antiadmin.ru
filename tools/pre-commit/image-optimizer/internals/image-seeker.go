package internals

import (
	"fmt"
	"os"
	"strings"
)

func SearchForImages(dirEntry string) ([]string, error) {
	entries, err := os.ReadDir(dirEntry)

	if err != nil {
		return nil, err
	}
	images, err := filterImageFiles(entries, dirEntry)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func filterImageFiles(entries []os.DirEntry, currentPath string) ([]string, error) {
	files := make([]string, 0)

	for _, file := range entries {
		entry := fmt.Sprintf("%s/%s", currentPath, file.Name())
		if !file.IsDir() {
			if isImage(file.Name()) {
				files = append(files, entry)
			}
		} else {
			dir, err := os.ReadDir(entry)
			if err != nil {
				return nil, err
			}

			subEntries, err := filterImageFiles(dir, entry)
			if err != nil {
				return nil, err
			}
			files = append(files, subEntries...)
		}
	}

	return files, nil
}

func isImage(filename string) bool {
	return strings.HasSuffix(filename, ".jpg") || strings.HasSuffix(filename, ".jpeg") || strings.HasSuffix(filename, ".png")
}
