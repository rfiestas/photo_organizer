package photoorganizer

import (
	"fmt"
)

func Generate(imageList []string, describer Describer) error {
	for _, image := range imageList {
		fmt.Println("Processing:", image)
		description, tags, err := describer.Describe(image)
		if err != nil {
			return fmt.Errorf("failed describing image %s: %v", image, err)

		}

		if err := generateXMPFile(image, description, tags); err != nil {
			return fmt.Errorf("failed to write XMP for %s: %v", image, err)
		}
	}
	return nil
}
