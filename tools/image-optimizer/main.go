package main

import (
	"fmt"
	"github.com/chloyka/chloyka.com/tools/image-optimizer/internals"
	"github.com/spf13/cobra"
	"os"
)

var defaultCmd = &cobra.Command{
	Use:   "optimize",
	Short: "Optimize images in a directory",
	Long:  "===================================== \n Optimize images in a directory \n=====================================",
}

func main() {
	defaultCmd.Flags().StringP("dir", "d", "./", "Directory to optimize")
	err := defaultCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	dir := defaultCmd.Flags().Lookup("dir")

	// Get all files in dir recursively
	files, err := internals.SearchForImages(dir.Value.String())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Optimize all files
	optimizer := internals.NewImageOptimizer()
	optimizer.Optimize(files)
}
