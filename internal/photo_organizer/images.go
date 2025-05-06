package photoorganizer

import (
	"os"
	"path/filepath"
	"strings"
)

func FindImagesInFolder(folder string) ([]string, error) {
	var images []string
	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !isImageFile(path) {
			return nil
		}
		images = append(images, path)
		return nil
	})
	return images, err
}

var supportedExtensions = map[string]struct{}{
	".jpg": {}, ".jpeg": {}, ".png": {},
}

func isImageFile(filename string) bool {
	_, ok := supportedExtensions[strings.ToLower(filepath.Ext(filename))]
	return ok
}
