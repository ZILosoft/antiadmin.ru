package internals

import (
	"errors"
	"fmt"
	"github.com/h2non/bimg"
	"strings"
)

type ImageOptimizer struct{}

func NewImageOptimizer() *ImageOptimizer {
	return &ImageOptimizer{}
}

func (io *ImageOptimizer) Optimize(entry string) error {
	if strings.HasSuffix(entry, ".jpg") || strings.HasSuffix(entry, ".jpeg") {
		return io.OptimizeJpeg(entry)
	} else if strings.HasSuffix(entry, ".png") {
		return io.OptimizePng(entry)
	}

	return errors.New(fmt.Sprintf("file %s has unsapported format", entry))
}

func (io *ImageOptimizer) OptimizeJpeg(fileName string) error {
	buffer, err := bimg.Read(fileName)
	if err != nil {
		return err
	}
	process, err := bimg.NewImage(buffer).Process(bimg.Options{
		Quality:       70,
		StripMetadata: true,
	})
	if err != nil {
		return err
	}

	err = bimg.Write(fileName, process)
	if err != nil {
		return err
	}

	return nil
}

func (io *ImageOptimizer) OptimizePng(fileName string) error {
	buffer, err := bimg.Read(fileName)
	if err != nil {
		return err
	}

	// Process and optimize the image using bimg
	process, err := bimg.NewImage(buffer).Process(bimg.Options{
		Quality:       70,
		StripMetadata: true,
	})
	if err != nil {
		return err
	}

	// Write the optimized image to the file
	return bimg.Write(fileName, process)
}
