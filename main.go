package main

import (
	"log"
	"mesh/src/components"
	"net/http"
)

// Index page handler
func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := components.App().Render(r.Context(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
