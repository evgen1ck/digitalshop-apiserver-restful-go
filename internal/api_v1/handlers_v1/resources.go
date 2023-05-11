package handlers_v1

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"test-server-go/internal/api_v1"
	tl "test-server-go/internal/tools"
)

func (rs *Resolver) ResourcesGetProductImage(w http.ResponseWriter, r *http.Request) {
	path, err := tl.GetExecutablePath()
	if err != nil {
		rs.App.Logger.NewWarn("error in get executable path", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	dir := filepath.Join(path, "resources", "product_images")
	files, err := os.ReadDir(dir)
	if err != nil {
		rs.App.Logger.NewWarn("error in read files directory", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	var foundFile string
	id := chi.URLParam(r, "id")
	for _, file := range files {
		if !file.IsDir() {
			filename := file.Name()
			fileID := strings.TrimSuffix(filename, filepath.Ext(filename))
			if fileID == id {
				foundFile = filename
				break
			}
		}
	}

	if foundFile == "" {
		api_v1.RedRespond(w, http.StatusNotFound, "Not found", "This file not found")
		return
	}

	w.Header().Set("Content-Type", "image")
	http.ServeFile(w, r, filepath.Join(dir, foundFile))
}

func (rs *Resolver) ResourcesGetSvgFile(w http.ResponseWriter, r *http.Request) {
	path, err := tl.GetExecutablePath()
	if err != nil {
		rs.App.Logger.NewWarn("error in get executable path", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	dir := filepath.Join(path, "resources", "svg_files")
	fileName := chi.URLParam(r, "id") + ".svg"
	fullPath := filepath.Join(dir, fileName)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		api_v1.RedRespond(w, http.StatusNotFound, "Not found", "This file not found")
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	http.ServeFile(w, r, fullPath)
}
