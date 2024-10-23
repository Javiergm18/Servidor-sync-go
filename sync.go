// sync.go
package main

import (
	"database/sql"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var baseDir = "./sync" // Asegúrate de definir el directorio base

// Conexión a la base de datos
func connectDB() (*sql.DB, error) {
	dsn := "usuario:contraseña@tcp(127.0.0.1:3306)/ServiciosArchivos" // Cambiar
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

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

	// Guardar metadatos en la base de datos
	db, err := connectDB()
	if err != nil {
		http.Error(w, "Error de conexión a la base de datos", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	idcreador := 1 // Cambiar esto según lógica de autenticación

	_, err = db.Exec("INSERT INTO Archivos (nombre, extension, path, tamano, idcreador) VALUES (?, ?, ?, ?, ?)",
		header.Filename, filepath.Ext(header.Filename), filePath, header.Size, idcreador)
	if err != nil {
		http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Archivo " + header.Filename + " subido con éxito.\n"))
}

// Listar archivos en el servidor
func ListFilesHandler(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		http.Error(w, "Error de conexión a la base de datos", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT nombre FROM Archivos")
	if err != nil {
		http.Error(w, "Error en la consulta", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var fileList []string
	for rows.Next() {
		var nombre string
		if err := rows.Scan(&nombre); err != nil {
			http.Error(w, "Error al escanear resultados", http.StatusInternalServerError)
			return
		}
		fileList = append(fileList, nombre)
	}

	w.Write([]byte(strings.Join(fileList, "\n")))
}
func SyncDirectoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	dirName := r.URL.Query().Get("dirname")
	dirPath := filepath.Join(baseDir, dirName)

	// Crear el directorio si no existe
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		os.MkdirAll(dirPath, os.ModePerm)
	}

	r.ParseMultipartForm(10 << 20) // Límite de tamaño 10 MB
	files := r.MultipartForm.File["files"]

	db, err := connectDB()
	if err != nil {
		http.Error(w, "Error de conexión a la base de datos", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Error al abrir el archivo", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Ruta donde se guardará el archivo
		filePath := filepath.Join(dirPath, fileHeader.Filename)

		outFile, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Error al crear el archivo en el servidor", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		// Copiar el archivo cargado al servidor
		_, err = io.Copy(outFile, file)
		if err != nil {
			http.Error(w, "Error al guardar el archivo", http.StatusInternalServerError)
			return
		}

		// Guardar metadatos en la base de datos
		idcreador := 1 // Cambiar esto según  lógica de autenticación

		_, err = db.Exec("INSERT INTO Archivos (nombre, extension, path, tamano, idcreador) VALUES (?, ?, ?, ?, ?)",
			fileHeader.Filename, filepath.Ext(fileHeader.Filename), filePath, fileHeader.Size, idcreador)
		if err != nil {
			http.Error(w, "Error al guardar en la base de datos", http.StatusInternalServerError)
			return
		}
	}

	w.Write([]byte("Sincronización completada con éxito.\n"))
}
