// main.go
package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	// Cargar plantilla HTML
	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	// Rutas del servidor
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	})

	// Servir archivos estáticos como CSS o JS
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Rutas para APIs
	http.HandleFunc("/sync", JWTMiddleware(SyncDirectoryHandler)) // Sincronización protegida con autenticación JWT
	http.HandleFunc("/list", JWTMiddleware(ListFilesHandler))     // Listado de archivos
	http.HandleFunc("/upload", JWTMiddleware(UploadFileHandler))  // Subida de archivos
	http.HandleFunc("/stream", JWTMiddleware(StreamVideoHandler)) // Streaming de videos
	http.HandleFunc("/login", LoginHandler)                       // Ruta de login para autenticación

	// Iniciar servidor en puerto 8080
	log.Println("Servidor iniciado en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
