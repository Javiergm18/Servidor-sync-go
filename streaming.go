
package main

import (
    "net/http"
    "os"
    "io"
    "path/filepath"
)

func StreamVideoHandler(w http.ResponseWriter, r *http.Request) {
    videoName := r.URL.Query().Get("video")
    videoPath := filepath.Join(baseDir, videoName)
    
    videoFile, err := os.Open(videoPath)
    if err != nil {
        http.Error(w, "Video no encontrado", http.StatusNotFound)
        return
    }
    defer videoFile.Close()

    w.Header().Set("Content-Type", "video/mp4")
    io.Copy(w, videoFile)
}
