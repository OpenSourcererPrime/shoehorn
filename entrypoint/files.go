package entrypoint

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"github.com/OpenSourcererPrime/shoehorn/config"
)

func (ep *EntryPoint) generateAllFiles() {
	for _, gen := range ep.appConfig.Generate {
		generateFile(gen)
	}
}

func generateFile(gen config.GenerateConfig) {
	// Ensure output directory exists
	outputDir := filepath.Dir(filepath.Join(gen.Path, gen.Name))
	err := os.MkdirAll(outputDir, 0o755)
	if err != nil {
		log.Printf("Failed to create output directory %s: %v", outputDir, err)
		return
	}

	outputPath := filepath.Join(gen.Path, gen.Name)

	switch gen.Strategy {
	case "append":
		// Read and concatenate all input files
		var buffer bytes.Buffer
		for _, input := range gen.Inputs {
			data, err := os.ReadFile(input.Path)
			if err != nil {
				log.Printf("Failed to read input file %s: %v", input.Path, err)
				continue
			}
			buffer.Write(data)
			// Add newline if not present at the end of the file
			if len(data) > 0 && data[len(data)-1] != '\n' {
				buffer.WriteString("\n")
			}
		}

		err = os.WriteFile(outputPath, buffer.Bytes(), 0o644)
		if err != nil {
			log.Printf("Failed to write output file %s: %v", outputPath, err)
		} else {
			log.Printf("Successfully generated %s (append strategy)", outputPath)
		}

	case "template":
		// Read the template file
		templateData, err := os.ReadFile(gen.Template)
		if err != nil {
			log.Printf("Failed to read template file %s: %v", gen.Template, err)
			return
		}

		// Create a template context with input files
		context := make(map[string]string)
		for _, input := range gen.Inputs {
			data, err := os.ReadFile(input.Path)
			if err != nil {
				log.Printf("Failed to read input file %s: %v", input.Path, err)
				context[input.Name] = "" // Set empty content if file can't be read
			} else {
				context[input.Name] = string(data)
			}
		}

		// Process the template
		tmpl, err := template.New("output").Parse(string(templateData))
		if err != nil {
			log.Printf("Failed to parse template %s: %v", gen.Template, err)
			return
		}

		var buffer bytes.Buffer
		err = tmpl.Execute(&buffer, context)
		if err != nil {
			log.Printf("Failed to execute template for %s: %v", outputPath, err)
			return
		}

		err = os.WriteFile(outputPath, buffer.Bytes(), 0o644)
		if err != nil {
			log.Printf("Failed to write output file %s: %v", outputPath, err)
		} else {
			log.Printf("Successfully generated %s (template strategy)", outputPath)
		}
	}
}
