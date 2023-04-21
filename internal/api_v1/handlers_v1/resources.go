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

	id := chi.URLParam(r, "id")
	dir := filepath.Join(path, "resources", "product_images")

	files, err := os.ReadDir(dir)
	if err != nil {
		rs.App.Logger.NewWarn("error in read files directory", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	var foundFile string
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

	http.ServeFile(w, r, filepath.Join(dir, foundFile))

}

func (rs *Resolver) ResourcesGetAvatarImage(w http.ResponseWriter, r *http.Request) {}
func (rs *Resolver) ResourcesGetSvgFile(w http.ResponseWriter, r *http.Request)     {}
