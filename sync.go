// sync.go
package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var baseDir = "./sync" // Asegúrate de definir el directorio base

// Función para manejar la carga de archivos
func UploadFileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No se pudo leer el archivo", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filePath := filepath.Join(baseDir, header.Filename)

	outFile, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "No se pudo crear el archivo en el servidor", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	// Copiar el archivo cargado al servidor
	_, err = io.Copy(outFile, file)
	if err != nil {
		http.Error(w, "Error al guardar el archivo", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Archivo " + header.Filename + " subido con éxito.\n"))
}

// Listar archivos en el servidor
func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(baseDir)
	if err != nil {
		http.Error(w, "Error al leer directorios", http.StatusInternalServerError)
		return
	}

	var fileList []string
	for _, f := range files {
		fileList = append(fileList, f.Name())
	}

	w.Write([]byte(strings.Join(fileList, "\n")))
}

// Sincronización de archivos y directorios
func SyncDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	dirName := r.URL.Query().Get("dirname")
	dirPath := filepath.Join(baseDir, dirName)

	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, os.ModePerm)
	}

	r.ParseMultipartForm(10 << 20) // Límite de tamaño 10 MB
	files := r.MultipartForm.File["files"]

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Error al abrir archivo", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		dst, err := os.Create(filepath.Join(dirPath, fileHeader.Filename))
		if err != nil {
			http.Error(w, "Error al crear archivo", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		io.Copy(dst, file)
	}

	w.Write([]byte("Sincronización completada"))
}
