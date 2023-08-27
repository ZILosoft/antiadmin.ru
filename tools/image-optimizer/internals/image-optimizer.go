package internals

import (
	"fmt"
	"github.com/h2non/bimg"
	"strings"
	"sync"
	"time"
)

type ImageOptimizer struct {
	files    chan string
	err      chan error
	progress chan int
	wg       *sync.WaitGroup
}

func NewImageOptimizer() *ImageOptimizer {
	return &ImageOptimizer{
		wg:       &sync.WaitGroup{},
		files:    make(chan string),
		err:      make(chan error),
		progress: make(chan int),
	}
}

func (io *ImageOptimizer) Optimize(entries []string) {
	io.wg.Add(len(entries))

	go func() {
		for _, fileName := range entries {
			io.files <- fileName
		}
	}()

	go func() {
		done := 0
		for pg := range io.progress {
			io.wg.Done()
			done += pg
			fmt.Printf(fmt.Sprintf("\rProgress: %d/%d", done, len(entries)))
		}
	}()

	go func() {
		for file := range io.files {
			time.Sleep(time.Second)
			go func(file string) {
				if strings.HasSuffix(file, ".jpg") || strings.HasSuffix(file, ".jpeg") {
					err := io.OptimizeJpeg(file)
					if err != nil {
						fmt.Println(err)
					}
					time.Sleep(time.Second)
					io.progress <- 1
				} else {
					err := io.OptimizePng(file)
					if err != nil {
						fmt.Println(err)
					}
					io.progress <- 1
				}
			}(file)
		}
	}()

	io.wg.Wait()
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

	//if len(process) < len(buffer) {
	//	fmt.Println(fmt.Sprintf("\r Optimized %s", fileName))
	//} else {
	//	fmt.Println(fmt.Sprintf("\r Skipped %s", fileName))
	//	return nil
	//}

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

	//if len(process) < len(buffer) {
	//	fmt.Println(fmt.Sprintf("\r Optimized %s", fileName))
	//} else {
	//	fmt.Println(fmt.Sprintf("\r Skipped %s", fileName))
	//	return nil
	//}

	// Write the optimized image to the file
	return bimg.Write(fileName, process)
}
