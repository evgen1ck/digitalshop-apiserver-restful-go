package handlers_v1

import (
	"net/http"
	"test-server-go/internal/api_v1"
	"test-server-go/internal/storage"
)

func (rs *Resolver) AdminGetProducts(w http.ResponseWriter, r *http.Request)    {}
func (rs *Resolver) AdminCreateProduct(w http.ResponseWriter, r *http.Request)  {}
func (rs *Resolver) AdminProductsUpdate(w http.ResponseWriter, r *http.Request) {}
func (rs *Resolver) AdminProductsDelete(w http.ResponseWriter, r *http.Request) {}

func (rs *Resolver) AdminProductsServicesGet(w http.ResponseWriter, r *http.Request) {
	services, err := storage.GetProductServices(r.Context(), rs.App.Postgres, rs.App.Config.App.Service.Url.Server)
	if err != nil {
		rs.App.Logger.NewWarn("error in get product services", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	api_v1.RespondWithCreated(w, services)
}
