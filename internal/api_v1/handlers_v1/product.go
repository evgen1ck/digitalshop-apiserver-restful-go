package handlers_v1

import (
	"net/http"
	"test-server-go/internal/api_v1"
	"test-server-go/internal/storage"
	tl "test-server-go/internal/tools"
)

func (rs *Resolver) ProductsDataForMainpage(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id != "" {
		if err := tl.Validate(id, tl.UuidFieldValidators(true)...); err != nil {
			api_v1.RespondWithUnprocessableEntity(w, "Id: "+err.Error())
			return
		}
	}
	sortBy := r.FormValue("sort_by")
	sortType := r.FormValue("sort_type")
	searchText := r.FormValue("search")

	// Block 1 - get products for mainpage
	products, err := storage.GetProductsForMainpage(r.Context(), rs.App.Postgres, rs.App.Config.App.Service.Url.Server, id, searchText, sortBy, sortType)
	if err != nil {
		rs.App.Logger.NewWarn("error in get products for mainpage", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	// Block 2 - send the result
	api_v1.RespondOK(w, products)
}
