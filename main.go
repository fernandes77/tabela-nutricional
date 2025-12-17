package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"tabela-nutricional/api"
)

func main() {
	// Get the executable directory to find the database
	execPath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	execDir := filepath.Dir(execPath)

	// Try to find the database in different locations
	dbPaths := []string{
		"tabelas-nutricionais.db",
		filepath.Join(execDir, "tabelas-nutricionais.db"),
		"./tabelas-nutricionais.db",
	}

	var database *api.DB
	for _, dbPath := range dbPaths {
		database, err = api.New(dbPath)
		if err == nil {
			log.Printf("Using database: %s", dbPath)
			break
		}
	}

	if database == nil {
		log.Fatal("Could not find database file")
	}
	defer database.Close()

	handler := api.NewHandler(database)

	// API routes
	http.HandleFunc("/api/foods", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			handler.HandleCORS(w, r)
			return
		}
		handler.GetFoods(w, r)
	})

	// Serve static files from web/dist
	distPath := "web/dist"
	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		// Try relative to executable
		distPath = filepath.Join(execDir, "web/dist")
	}

	fs := http.FileServer(http.Dir(distPath))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Check if the request is for an API route
		if len(r.URL.Path) >= 4 && r.URL.Path[:4] == "/api" {
			http.NotFound(w, r)
			return
		}

		// Try to serve the file directly
		path := filepath.Join(distPath, r.URL.Path)

		// If the path doesn't exist and doesn't have an extension, try with .html
		if _, err := os.Stat(path); os.IsNotExist(err) {
			htmlPath := path + ".html"
			if _, err := os.Stat(htmlPath); err == nil {
				r.URL.Path = r.URL.Path + ".html"
			} else {
				// Try index.html for directory paths
				indexPath := filepath.Join(path, "index.html")
				if _, err := os.Stat(indexPath); err == nil {
					r.URL.Path = filepath.Join(r.URL.Path, "index.html")
				}
			}
		}

		fs.ServeHTTP(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Printf("Server starting on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
