package main

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"log/slog"
	"mesh/src/components"
	"mesh/src/handlers/app"
	"mesh/src/handlers/board"
	"mesh/src/services"
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

func main() {
	// Create logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Create services
	counterService := services.NewCounterService()

	// Create handlers
	boardHandler := board.New(logger, counterService)
	appHandler := app.New(logger, boardHandler)

	// Initialize component registry (optional - for templ template functions)
	components.SetRegistry(appHandler, boardHandler)

	// Index page handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Add component registry to context for template functions
		ctx := components.WithRegistry(r.Context(), &components.Registry{
			AppHandler:   appHandler,
			BoardHandler: boardHandler,
		})
		r = r.WithContext(ctx)

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
		appComponent := appHandler.RenderComponent(r.Context())
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
	})

	// Route handlers with registry context middleware
	http.HandleFunc("/app", func(w http.ResponseWriter, r *http.Request) {
		ctx := components.WithRegistry(r.Context(), &components.Registry{
			AppHandler:   appHandler,
			BoardHandler: boardHandler,
		})
		appHandler.ServeHTTP(w, r.WithContext(ctx))
	})

	http.HandleFunc("/board", func(w http.ResponseWriter, r *http.Request) {
		ctx := components.WithRegistry(r.Context(), &components.Registry{
			AppHandler:   appHandler,
			BoardHandler: boardHandler,
		})
		boardHandler.ServeHTTP(w, r.WithContext(ctx))
	})

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("Server starting on :8000")
	log.Println("Routes:")
	log.Println("  GET  / - Main page with counter demo")
	log.Println("  GET  /app - App component")
	log.Println("  GET  /board - Board component with counter")
	log.Println("  POST /board - Increment counter (HTMX)")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
