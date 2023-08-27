package main

import (
	"fmt"
	"github.com/chloyka/chloyka.com/tools/image-optimizer/internals"
	"os"
)

var infoText = "===================================== \n Optimize image by provided path\n Usage: optimize /path/to/image.jpeg \n====================================="

func main() {
	if len(os.Args) < 2 || os.Args[1] == "" {
		printHelp()
		return
	}

	fileName := os.Args[1]

	// Check if file exists
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		fmt.Printf("provided file %s does not exist", fileName)
		os.Exit(1)
	}

	optimizer := internals.NewImageOptimizer()
	optimizer.Optimize(fileName)
}

func printHelp() {
	fmt.Println(infoText)
	os.Exit(1)
}
