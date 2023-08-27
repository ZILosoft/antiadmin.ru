package main

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/adrg/frontmatter"
	"github.com/chloyka/chloyka.com/tools/content-translator/internals/translator"
	"html/template"
	"os"
	"strings"
)

var infoText = "===================================== \n Translate text if translation not found \n Usage: translate /path/to/md-file.md \n====================================="

type Post struct {
	Title string   `yaml:"title"`
	Tags  []string `yaml:"tags"`
	Date  string   `yaml:"date"`
	Draft bool     `yaml:"draft"`

	Body string `yaml:"-"`
}

//go:embed templates
var templatesFS embed.FS

func main() {
	if len(os.Args) < 2 || os.Args[1] == "" {
		printHelp()
		return
	}

	filePath := os.Args[1]

	err := processFile(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func processFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("provided file %s does not exist", filePath)
	}

	filePaths := strings.Split(filePath, ".")
	if len(filePaths) < 3 {
		return nil
	}

	lang := translator.Language("")
	err := lang.Parse(filePaths[len(filePaths)-2])
	if err != nil {
		return err
	}
	filePaths[len(filePaths)-2] = string(lang.GetAntipode())
	newFilePath := strings.Join(filePaths, ".")

	if _, err = os.Stat(newFilePath); err == nil {
		// Translations file exists, do nothing
		return nil
	}

	token := os.Getenv("OPENAI_TOKEN")
	if token == "" {
		return fmt.Errorf("OPENAI_TOKEN env variable is not set")
	}

	client := translator.NewTranslator(token)
	file, _ := os.ReadFile(filePath)
	if len(file) == 0 {
		return fmt.Errorf("file %s is empty", filePath)
	}

	var head Post
	body, err := frontmatter.Parse(bytes.NewReader(file), &head)
	if err != nil {
		return err
	}

	bodyTranslation, err := client.Translate(string(body), lang, lang.GetAntipode())
	if err != nil {
		return err
	}

	titleTranslation, err := client.Translate(head.Title, lang, lang.GetAntipode())
	if err != nil {
		return err
	}

	templateContent, err := templatesFS.ReadFile("templates/post.md.tmpl")
	if err != nil {
		return err
	}

	tmpl := template.Must(template.New("post.md").Parse(string(templateContent)))

	var outputBuffer bytes.Buffer
	err = tmpl.Execute(&outputBuffer, &Post{
		Body:  string(bodyTranslation),
		Date:  head.Date,
		Draft: head.Draft,
		Tags:  head.Tags,
		Title: string(titleTranslation),
	})
	if err != nil {
		return err
	}

	err = os.WriteFile(newFilePath, outputBuffer.Bytes(), 0644)
	if err != nil {
		return err
	}

	// always exit with code 1, because we want to commit new files manually
	return fmt.Errorf("translation file %s for file %s created", newFilePath, filePath)
}

func printHelp() {
	fmt.Println(infoText)
	os.Exit(1)
}
