package src

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"mesh/src/components"
	"net/http"
	"os"
	"strings"
)

type ViteManifestEntry struct {
	File string   `json:"file"`
	CSS  []string `json:"css,omitempty"`
}

type ViteManifest map[string]ViteManifestEntry

type TemplateData struct {
	Css string
	Js  string
	App template.HTML
}

func loadViteManifest() (ViteManifest, error) {
	data, err := os.ReadFile("static/.vite/manifest.json")
	if err != nil {
		return nil, err
	}

	var manifest ViteManifest
	err = json.Unmarshal(data, &manifest)
	if err != nil {
		return nil, err
	}

	return manifest, nil
}

func getAssetsFromManifest(manifest ViteManifest) (string, string) {
	var cssFiles []string
	var jsFile string

	if entry, exists := manifest["src/main.ts"]; exists {
		jsFile = entry.File
		cssFiles = append(cssFiles, entry.CSS...)
	}

	for key, entry := range manifest {
		if strings.HasSuffix(key, ".scss") || strings.HasSuffix(key, ".css") {
			cssFiles = append(cssFiles, entry.File)
		}
	}

	var cssFile string
	if len(cssFiles) > 0 {
		cssFile = cssFiles[0] // Using the first CSS file
	}

	return cssFile, jsFile
}

func IndexHandler(registry *components.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		manifest, err := loadViteManifest()
		if err != nil {
			log.Printf("Error loading Vite manifest: %v", err)
			manifest = make(ViteManifest)
		}

		cssFile, jsFile := getAssetsFromManifest(manifest)

		tmplContent, err := os.ReadFile("index.html")
		if err != nil {
			log.Printf("Error loading index.html: %v", err)
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}

		tmpl := template.New("index")
		tmpl, err = tmpl.Parse(string(tmplContent))
		if err != nil {
			log.Printf("Error parsing template: %v", err)
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			return
		}

		buf := new(bytes.Buffer)
		appComponent := registry.AppHandler.RenderComponent()
		err = appComponent.Render(r.Context(), buf)
		if err != nil {
			log.Printf("Error rendering app template: %v", err)
			http.Error(w, "Error rendering app template", http.StatusInternalServerError)
			return
		}

		data := TemplateData{
			Css: cssFile,
			Js:  jsFile,
			App: template.HTML(buf.String()),
		}

		// Execute the template
		err = tmpl.Execute(w, data)
		if err != nil {
			log.Printf("Error executing template: %v", err)
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
	}
}
