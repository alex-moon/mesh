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

// ViteManifestEntry represents an entry in the Vite manifest
type ViteManifestEntry struct {
	File string   `json:"file"`
	CSS  []string `json:"css,omitempty"`
}

// ViteManifest represents the entire Vite manifest
type ViteManifest map[string]ViteManifestEntry

// TemplateData holds the data to be passed to the template
type TemplateData struct {
	Css string
	Js  string
	App template.HTML
}

// loadViteManifest reads and parses the Vite manifest file
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

// getAssetsFromManifest extracts CSS and JS file paths from the manifest
func getAssetsFromManifest(manifest ViteManifest) (string, string) {
	var cssFiles []string
	var jsFile string

	// Look for main.ts entry point
	if entry, exists := manifest["src/main.ts"]; exists {
		jsFile = entry.File
		cssFiles = append(cssFiles, entry.CSS...)
	}

	// Look for additional CSS entries
	for key, entry := range manifest {
		if strings.HasSuffix(key, ".scss") || strings.HasSuffix(key, ".css") {
			cssFiles = append(cssFiles, entry.File)
		}
	}

	// Return the first CSS file found (or combine if needed)
	var cssFile string
	if len(cssFiles) > 0 {
		cssFile = cssFiles[0] // Using the first CSS file
	}

	return cssFile, jsFile
}

// IndexHandler handles the index page request
func IndexHandler(registry *components.Registry) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Load the Vite manifest
		manifest, err := loadViteManifest()
		if err != nil {
			log.Printf("Error loading Vite manifest: %v", err)
			// Continue without manifest for development
			manifest = make(ViteManifest)
		}

		// Get CSS and JS file paths from manifest
		cssFile, jsFile := getAssetsFromManifest(manifest)

		// Load the index.html template
		tmplContent, err := os.ReadFile("index.html")
		if err != nil {
			log.Printf("Error loading index.html: %v", err)
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}

		// Create a new template and define the templ-content template
		tmpl := template.New("index")
		tmpl, err = tmpl.Parse(string(tmplContent))
		if err != nil {
			log.Printf("Error parsing template: %v", err)
			http.Error(w, "Error parsing template", http.StatusInternalServerError)
			return
		}

		// Render the app component using the handler
		buf := new(bytes.Buffer)
		appComponent := registry.AppHandler.RenderComponent(r.Context())
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