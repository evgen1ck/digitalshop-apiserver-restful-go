package handlers_v1

import (
	"encoding/json"
	"net/http"
	"test-server-go/internal/api_v1"
	"test-server-go/internal/storage"
	tl "test-server-go/internal/tools"
)

func (rs *Resolver) AdminNull(w http.ResponseWriter, r *http.Request)           {}
func (rs *Resolver) AdminCreateProduct(w http.ResponseWriter, r *http.Request)  {}
func (rs *Resolver) AdminProductsUpdate(w http.ResponseWriter, r *http.Request) {}
func (rs *Resolver) AdminProductsDelete(w http.ResponseWriter, r *http.Request) {}

func (rs *Resolver) AdminGetVariants(w http.ResponseWriter, r *http.Request) {
	products, err := storage.GetAdminVariants(r.Context(), rs.App.Postgres, rs.App.Config.App.Service.Url.Server, "")
	if err != nil {
		rs.App.Logger.NewWarn("error in get variants", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	api_v1.RespondOK(w, products)
}

func (rs *Resolver) AdminGetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := storage.AdminGetProducts(r.Context(), rs.App.Postgres)
	if err != nil {
		rs.App.Logger.NewWarn("error in get variants", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	api_v1.RespondOK(w, products)
}

func (rs *Resolver) AdminGetServices(w http.ResponseWriter, r *http.Request) {
	services, err := storage.AdminGetServices(r.Context(), rs.App.Postgres, rs.App.Config.App.Service.Url.Server)
	if err != nil {
		rs.App.Logger.NewWarn("error in get services", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	api_v1.RespondOK(w, services)
}

func (rs *Resolver) AdminGetStates(w http.ResponseWriter, r *http.Request) {
	states, err := storage.AdminGetStates(r.Context(), rs.App.Postgres)
	if err != nil {
		rs.App.Logger.NewWarn("error in get states", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	api_v1.RespondOK(w, states)
}

func (rs *Resolver) AdminGetItems(w http.ResponseWriter, r *http.Request) {
	items, err := storage.AdminGetItems(r.Context(), rs.App.Postgres)
	if err != nil {
		rs.App.Logger.NewWarn("error in get items", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	api_v1.RespondOK(w, items)
}

func (rs *Resolver) AdminGetTypes(w http.ResponseWriter, r *http.Request) {
	types, err := storage.AdminGetTypes(r.Context(), rs.App.Postgres)
	if err != nil {
		rs.App.Logger.NewWarn("error in get types", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	api_v1.RespondOK(w, types)
}

func (rs *Resolver) AdminGetSubtypes(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}

	typeName := r.FormValue("type_name")
	if typeName == "" {
		api_v1.RespondWithUnprocessableEntity(w, "Type_name: the parameter value is empty")
		return
	}

	subtypes, err := storage.AdminGetSubtypes(r.Context(), rs.App.Postgres, typeName)
	if err != nil {
		rs.App.Logger.NewWarn("error in get subtypes", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	api_v1.RespondOK(w, subtypes)
}

type AdminCreateVariantData struct {
	ProductName     string  `json:"product_name"`
	VariantName     string  `json:"variant_name"`
	Service         string  `json:"service_name"`
	State           string  `json:"state_name"`
	Subtype         string  `json:"subtype_name"`
	Item            string  `json:"item_name"`
	Mask            string  `json:"mask"`
	Price           string  `json:"price"`
	DiscountMoney   *string `json:"discount_money"`
	DiscountPercent *string `json:"discount_percent"`
}

func (rs *Resolver) AdminCreateVariant(w http.ResponseWriter, r *http.Request) {
	// Block 0 - decode data
	var data AdminCreateVariantData
	decodeErr := json.NewDecoder(r.Body).Decode(&data)
	if decodeErr != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}

	// Block 1 - data validation
	if err := tl.Validate(data.ProductName, tl.IsNotBlank(), tl.IsMinMaxLen(MinTextLength, MaxTextLength), tl.IsNotContainsConsecutiveSpaces(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Product name: "+err.Error())
		return
	}
	if err := tl.Validate(data.VariantName, tl.IsNotBlank(), tl.IsMinMaxLen(MinTextLength, MaxTextLength), tl.IsNotContainsConsecutiveSpaces(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Variant name: "+err.Error())
		return
	}
	if err := tl.Validate(data.Service, tl.IsNotBlank(), tl.IsMinMaxLen(MinTextLength, MaxTextLength), tl.IsNotContainsConsecutiveSpaces(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Service: "+err.Error())
		return
	}
	if err := tl.Validate(data.State, tl.IsNotBlank(), tl.IsMinMaxLen(MinTextLength, MaxTextLength), tl.IsNotContainsConsecutiveSpaces(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "State: "+err.Error())
		return
	}
	if err := tl.Validate(data.Subtype, tl.IsNotBlank(), tl.IsMinMaxLen(MinTextLength, MaxTextLength), tl.IsNotContainsConsecutiveSpaces(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Subtype: "+err.Error())
		return
	}
	if err := tl.Validate(data.Item, tl.IsNotBlank(), tl.IsMinMaxLen(MinTextLength, MaxTextLength), tl.IsNotContainsConsecutiveSpaces(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Item: "+err.Error())
		return
	}
	if err := tl.Validate(data.Mask, tl.IsNotBlank(), tl.IsMinMaxLen(MinTextLength, MaxTextLength), tl.IsNotContainsConsecutiveSpaces(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Mask: "+err.Error())
		return
	}
	if err := tl.Validate(data.Price, tl.IsNotBlank(), tl.IsMoney(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Price: "+err.Error())
		return
	}
	if data.DiscountMoney != nil && *data.DiscountMoney != "" {
		if err := tl.Validate(*data.DiscountMoney, tl.IsMoney(), tl.IsTrimmedSpace()); err != nil {
			api_v1.RespondWithUnprocessableEntity(w, "Discount money: "+err.Error())
			return
		}
	}
	if data.DiscountPercent != nil && *data.DiscountPercent != "" {
		if err := tl.Validate(*data.DiscountPercent, tl.IsValidInteger(false), tl.IsNotContainsConsecutiveSpaces(), tl.IsTrimmedSpace()); err != nil {
			api_v1.RespondWithUnprocessableEntity(w, "Discount percent: "+err.Error())
			return
		}
	}

	_, jwtData, err := api_v1.ContextGetAuthenticated(r)
	if err != nil {
		rs.App.Logger.NewWarn("error in took jwt data", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	if err = storage.CreateAdminVariant(r.Context(), rs.App.Postgres, data.ProductName, data.VariantName, data.Service, data.State, data.Subtype, data.Item, data.Mask, data.Price, *data.DiscountMoney, *data.DiscountPercent, jwtData.AccountUuid); err != nil {
		api_v1.RespondWithInternalServerError(w)
		rs.App.Logger.NewWarn("Error in create admin variant", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
