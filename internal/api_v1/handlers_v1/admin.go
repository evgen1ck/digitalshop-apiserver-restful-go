package handlers_v1

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"test-server-go/internal/api_v1"
	"test-server-go/internal/storage"
	tl "test-server-go/internal/tools"
)

func (rs *Resolver) AdminGetVariants(w http.ResponseWriter, r *http.Request) {
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
	activeFirst := r.FormValue("active_first")

	products, err := storage.GetAdminVariants(r.Context(), rs.App.Postgres, rs.App.Config.App.Service.Url.Server, id, searchText, sortBy, sortType, activeFirst)
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
	serviceName := r.FormValue("service_name")

	services, err := storage.AdminGetServices(r.Context(), rs.App.Postgres, rs.App.Config.App.Service.Url.Server, serviceName)
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
	typ := r.FormValue("type_name")

	types, err := storage.AdminGetTypes(r.Context(), rs.App.Postgres, typ)
	if err != nil {
		rs.App.Logger.NewWarn("error in get types", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	api_v1.RespondOK(w, types)
}

func (rs *Resolver) AdminGetSubtypes(w http.ResponseWriter, r *http.Request) {
	subtypeName := r.FormValue("subtype_name")

	typeName := r.FormValue("type_name")

	subtypes, err := storage.AdminGetSubtypes(r.Context(), rs.App.Postgres, typeName, subtypeName)
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
	ServiceName     string  `json:"service_name"`
	StateName       string  `json:"state_name"`
	SubtypeName     string  `json:"subtype_name"`
	ItemName        string  `json:"item_name"`
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
	if err := tl.Validate(data.ProductName, tl.IsNotBlank(true), tl.IsMinMaxLen(MinTextLength, MaxTextLength), tl.IsNotContainsConsecutiveSpaces(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Product name: "+err.Error())
		fmt.Println("1")
		return
	}
	if err := tl.Validate(data.VariantName, tl.IsNotBlank(true), tl.IsMinMaxLen(MinTextLength, MaxTextLength), tl.IsNotContainsConsecutiveSpaces(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Variant name: "+err.Error())
		fmt.Println("2")
		return
	}
	if err := tl.Validate(data.ServiceName, tl.IsNotBlank(true), tl.IsMinMaxLen(MinTextLength, MaxTextLength), tl.IsNotContainsConsecutiveSpaces(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "ServiceName name: "+err.Error())
		fmt.Println("3")
		return
	}
	if err := tl.Validate(data.StateName, tl.IsNotBlank(true), tl.IsMinMaxLen(MinTextLength, MaxTextLength), tl.IsNotContainsConsecutiveSpaces(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "StateName name: "+err.Error())
		fmt.Println("4")
		return
	}
	if err := tl.Validate(data.SubtypeName, tl.IsNotBlank(true), tl.IsMinMaxLen(MinTextLength, MaxTextLength), tl.IsNotContainsConsecutiveSpaces(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "SubtypeName name: "+err.Error())
		fmt.Println("5")
		return
	}
	if err := tl.Validate(data.ItemName, tl.IsNotBlank(true), tl.IsMinMaxLen(MinTextLength, MaxTextLength), tl.IsNotContainsConsecutiveSpaces(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "ItemName name: "+err.Error())
		fmt.Println("6")
		return
	}
	if err := tl.Validate(data.Mask, tl.IsNotBlank(true), tl.IsMinMaxLen(MinTextLength, MaxTextLength), tl.IsNotContainsConsecutiveSpaces(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Mask: "+err.Error())
		fmt.Println("7")
		return
	}
	if err := tl.Validate(data.Price, tl.IsNotBlank(true), tl.IsMoney(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Price: "+err.Error())
		fmt.Println("8")
		return
	}
	if data.DiscountMoney != nil && *data.DiscountMoney != "" {
		if err := tl.Validate(*data.DiscountMoney, tl.IsMoney(), tl.IsTrimmedSpace()); err != nil {
			api_v1.RespondWithUnprocessableEntity(w, "Discount money: "+err.Error())
			fmt.Println("9")
			return
		}
	}
	if data.DiscountPercent != nil && *data.DiscountPercent != "" {
		if err := tl.Validate(*data.DiscountPercent, tl.IsValidInteger(false, true), tl.IsNotContainsConsecutiveSpaces(), tl.IsTrimmedSpace()); err != nil {
			api_v1.RespondWithUnprocessableEntity(w, "Discount percent: "+err.Error())
			fmt.Println("10")
			return
		}
	}

	_, jwtData, err := api_v1.ContextGetAuthenticated(r)
	if err != nil {
		rs.App.Logger.NewWarn("error in took jwt data", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	if err = storage.CreateAdminVariant(r.Context(), rs.App.Postgres, data.ProductName, data.VariantName, data.ServiceName, data.StateName, data.SubtypeName, data.ItemName, data.Mask, data.Price, *data.DiscountMoney, *data.DiscountPercent, jwtData.AccountUuid); err != nil {
		api_v1.RespondWithInternalServerError(w)
		rs.App.Logger.NewWarn("Error in create admin variant", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type AdminUpdateVariantData struct {
	VariantName     *string `json:"variant_name"`
	StateName       *string `json:"state_name"`
	ItemName        *string `json:"item_name"`
	Mask            *string `json:"mask"`
	Price           *string `json:"price"`
	DiscountMoney   *string `json:"discount_money"`
	DiscountPercent *string `json:"discount_percent"`
}

func (rs *Resolver) AdminUpdateVariant(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}

	id := r.FormValue("id")
	if id == "" {
		api_v1.RespondWithUnprocessableEntity(w, "Id: the parameter value is empty")
		return
	}
	if err = tl.Validate(id, tl.UuidFieldValidators(true)...); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Id: "+err.Error())
		return
	}

	var data AdminUpdateVariantData
	decodeErr := json.NewDecoder(r.Body).Decode(&data)
	if decodeErr != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}

	var updateData = make(map[string]interface{})
	if data.VariantName != nil {
		if err = tl.Validate(*data.VariantName, tl.TextFieldValidatorsWithSpaces()...); err != nil {
			api_v1.RespondWithUnprocessableEntity(w, "Variant name: "+err.Error())
			return
		}
		updateData["variant_name"] = *data.VariantName
	}
	if data.StateName != nil {
		if err = tl.Validate(*data.StateName, tl.TextFieldValidatorsWithSpaces()...); err != nil {
			api_v1.RespondWithUnprocessableEntity(w, "State name: "+err.Error())
			return
		}
		no, err := storage.GetStateNo(r.Context(), rs.App.Postgres, *data.StateName)
		if err != nil {
			rs.App.Logger.NewWarn("Error in get state no: ", err)
			api_v1.RespondWithBadRequest(w, "")
			return
		}
		updateData["variant_state"] = no
	}
	if data.ItemName != nil {
		if err = tl.Validate(*data.ItemName, tl.TextFieldValidatorsWithSpaces()...); err != nil {
			api_v1.RespondWithUnprocessableEntity(w, "Item name: "+err.Error())
			return
		}
		no, err := storage.GetItemNo(r.Context(), rs.App.Postgres, *data.ItemName)
		if err != nil {
			rs.App.Logger.NewWarn("Error in get item no: ", err)
			api_v1.RespondWithBadRequest(w, "")
			return
		}
		updateData["variant_item"] = no
	}
	if data.Mask != nil {
		if err = tl.Validate(*data.Mask, tl.TextFieldValidatorsWithSpaces()...); err != nil {
			api_v1.RespondWithUnprocessableEntity(w, "Mask: "+err.Error())
			return
		}
		updateData["mask"] = *data.Mask
	}
	if data.Price != nil {
		if err = tl.Validate(*data.Price, tl.IsNotBlank(true), tl.IsMoney(), tl.IsTrimmedSpace()); err != nil {
			api_v1.RespondWithUnprocessableEntity(w, "Price: "+err.Error())
			return
		}
		updateData["price"] = *data.Price
	}
	if data.DiscountMoney != nil {
		if err = tl.Validate(*data.DiscountMoney, tl.IsMoney(), tl.IsTrimmedSpace()); err != nil {
			api_v1.RespondWithUnprocessableEntity(w, "Discount money: "+err.Error())
			return
		}
		updateData["discount_money"] = *data.DiscountMoney
	}
	if data.DiscountPercent != nil {
		if err = tl.Validate(*data.DiscountPercent, tl.IsValidInteger(false, true), tl.IsNotContainsConsecutiveSpaces(), tl.IsTrimmedSpace()); err != nil {
			api_v1.RespondWithUnprocessableEntity(w, "Discount percent: "+err.Error())
			return
		}
		updateData["discount_percent"] = *data.DiscountPercent
	}

	if len(updateData) == 0 {
		api_v1.RespondWithUnprocessableEntity(w, "No values")
		return
	}

	if err = storage.UpdateAdminVariant(r.Context(), rs.App.Postgres, id, updateData); err == storage.FailedUpdate {
		api_v1.RedRespond(w, http.StatusNotFound, "Not found", "Variant with this id not found")
		return
	} else if err != nil {
		api_v1.RespondWithInternalServerError(w)
		rs.App.Logger.NewWarn("Error in update admin variant", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rs *Resolver) AdminDeleteVariant(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}

	id := r.FormValue("id")
	if id == "" {
		api_v1.RespondWithUnprocessableEntity(w, "Id: the parameter value is empty")
		return
	}
	if err = tl.Validate(id, tl.UuidFieldValidators(true)...); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Id: "+err.Error())
		return
	}

	if err = storage.AdminDeleteVariant(r.Context(), rs.App.Postgres, id); err != nil {
		api_v1.RespondWithConflict(w, storage.PgErrorsHandle(err, "variant_name"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

type AdminUploadVariantData []struct {
	Data string `json:"data"`
}

func (rs *Resolver) AdminUploadVariant(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		api_v1.RespondWithUnprocessableEntity(w, "Id: the parameter value is empty")
		return
	}
	if err := tl.Validate(id, tl.UuidFieldValidators(true)...); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Id: "+err.Error())
		return
	}

	var data AdminUploadVariantData
	decodeErr := json.NewDecoder(r.Body).Decode(&data)
	if decodeErr != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}

	if len(data) == 0 {
		api_v1.RespondWithUnprocessableEntity(w, "No values")
		return
	}

	var dataList []string
	for i, obj := range data {
		if err := tl.Validate(obj.Data, tl.LongTextFieldValidatorsWithSpaces()...); err != nil {
			api_v1.RespondWithUnprocessableEntity(w, "Data["+strconv.Itoa(i+1)+"]: "+err.Error())
			return
		}
		dataList = append(dataList, obj.Data)
	}

	err := storage.CreateAdminContent(r.Context(), rs.App.Postgres, id, dataList)
	if err != nil {
		rs.App.Logger.NewWarn("error in create content", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rs *Resolver) AdminGetVariantUploads(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		api_v1.RespondWithUnprocessableEntity(w, "Id: the parameter value is empty")
		return
	}
	if err := tl.Validate(id, tl.UuidFieldValidators(true)...); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Id: "+err.Error())
		return
	}

	contents, err := storage.GetAdminContents(r.Context(), rs.App.Postgres, id)
	if err != nil {
		rs.App.Logger.NewWarn("error in get contents", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	api_v1.RespondOK(w, contents)
}

func (rs *Resolver) AdminDeleteVariantUpload(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		api_v1.RespondWithUnprocessableEntity(w, "Id: the parameter value is empty")
		return
	}
	if err := tl.Validate(id, tl.UuidFieldValidators(true)...); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Id: "+err.Error())
		return
	}

	if err := storage.DeleteAdminContent(r.Context(), rs.App.Postgres, id); err != nil {
		rs.App.Logger.NewWarn("error in delete content(s)", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rs *Resolver) AdminDeleteType(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		api_v1.RespondWithUnprocessableEntity(w, "Name: the parameter value is empty")
		return
	}

	if err := storage.DeleteAdminType(r.Context(), rs.App.Postgres, name); err != nil {
		api_v1.RespondWithConflict(w, storage.PgErrorsHandle(err, "Name"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rs *Resolver) AdminDeleteSubtype(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		api_v1.RespondWithUnprocessableEntity(w, "Name: the parameter value is empty")
		return
	}

	if err := storage.DeleteAdminSubtype(r.Context(), rs.App.Postgres, name); err != nil {
		api_v1.RespondWithConflict(w, storage.PgErrorsHandle(err, "Name"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rs *Resolver) AdminDeleteProduct(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		api_v1.RespondWithUnprocessableEntity(w, "Name: the parameter value is empty")
		return
	}

	if err := storage.DeleteAdminProduct(r.Context(), rs.App.Postgres, name); err != nil {
		api_v1.RespondWithConflict(w, storage.PgErrorsHandle(err, "Name"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rs *Resolver) AdminDeleteService(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		api_v1.RespondWithUnprocessableEntity(w, "Name: the parameter value is empty")
		return
	}

	if err := storage.DeleteAdminService(r.Context(), rs.App.Postgres, name); err != nil {
		api_v1.RespondWithConflict(w, storage.PgErrorsHandle(err, "Name"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rs *Resolver) AdminAddType(w http.ResponseWriter, r *http.Request) {
	var data struct {
		TypeName string `json:"type_name"`
	}
	decodeErr := json.NewDecoder(r.Body).Decode(&data)
	if decodeErr != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}
	if err := tl.Validate(data.TypeName, tl.TextFieldValidatorsWithSpaces()...); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Type_name: "+err.Error())
		return
	}

	if err := storage.CreateAdminType(r.Context(), rs.App.Postgres, data.TypeName); err != nil {
		api_v1.RespondWithConflict(w, storage.PgErrorsHandle(err, "Type_name"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rs *Resolver) AdminAddSubtype(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		api_v1.RespondWithUnprocessableEntity(w, "Name: the parameter value is empty")
		return
	}
	if err := tl.Validate(name, tl.TextFieldValidatorsWithSpaces()...); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Name: "+err.Error())
		return
	}

	var data struct {
		SubtypeName string `json:"subtype_name"`
	}
	decodeErr := json.NewDecoder(r.Body).Decode(&data)
	if decodeErr != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}
	if err := tl.Validate(data.SubtypeName, tl.TextFieldValidatorsWithSpaces()...); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "Subtype_name: "+err.Error())
		return
	}

	if err := storage.CreateAdminSubtype(r.Context(), rs.App.Postgres, data.SubtypeName, name); err != nil {
		api_v1.RespondWithConflict(w, storage.PgErrorsHandle(err, "Subtype_name"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rs *Resolver) AdminAddService(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	serviceName := r.FormValue("service_name")
	if serviceName == "" {
		api_v1.RespondWithBadRequest(w, "")
		return
	}

	if err := storage.CreateAdminService(r.Context(), rs.App.Postgres, serviceName); err != nil {
		api_v1.RespondWithConflict(w, storage.PgErrorsHandle(err, "Service_name"))
		return
	}

	file, handler, err := r.FormFile("file")
	serviceName = strings.ToLower(serviceName)
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	if !strings.HasSuffix(handler.Filename, ".svg") || handler.Header.Get("Content-Type") != "image/svg+xml" {
		api_v1.RespondWithBadRequest(w, "")
		return
	}

	path, err := tl.GetExecutablePath()
	dir := filepath.Join(path, "resources", "svg_files")
	fileName := strings.ReplaceAll(serviceName, " ", "-") + ".svg"
	fullPath := filepath.Join(dir, fileName)

	tempFile, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}
	defer tempFile.Close()
	io.Copy(tempFile, file)

	w.WriteHeader(http.StatusNoContent)
}

func (rs *Resolver) AdminAddProduct(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	productName := r.FormValue("product_name")
	tags := r.FormValue("tags")
	description := r.FormValue("description")
	if productName == "" {
		api_v1.RespondWithBadRequest(w, "")
		return
	}

	uuid, err := storage.CreateAdminProduct(r.Context(), rs.App.Postgres, productName, tags, description)
	if err != nil {
		api_v1.RespondWithConflict(w, storage.PgErrorsHandle(err, "Service_name"))
		return
	}

	file, handler, err := r.FormFile("file")
	productName = strings.ToLower(productName)
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	path, err := tl.GetExecutablePath()
	dir := filepath.Join(path, "resources", "product_images")
	fileName := uuid + filepath.Ext(handler.Filename)
	fullPath := filepath.Join(dir, fileName)

	tempFile, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}
	defer tempFile.Close()
	io.Copy(tempFile, file)

	w.WriteHeader(http.StatusNoContent)
}

func (rs *Resolver) AdminEditType(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("type_name")
	var data struct {
		TypeName string `json:"type_name"`
	}
	decodeErr := json.NewDecoder(r.Body).Decode(&data)
	if decodeErr != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}

	if err := storage.EditAdminType(r.Context(), rs.App.Postgres, name, data.TypeName); err != nil {
		api_v1.RespondWithConflict(w, storage.PgErrorsHandle(err, "Name"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rs *Resolver) AdminEditSubtype(w http.ResponseWriter, r *http.Request) {
	subtypeName := r.FormValue("subtype_name")
	var data struct {
		SubtypeName string `json:"subtype_name"`
	}
	decodeErr := json.NewDecoder(r.Body).Decode(&data)
	if decodeErr != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}

	if err := storage.EditAdminSubtype(r.Context(), rs.App.Postgres, subtypeName, data.SubtypeName); err != nil {
		api_v1.RespondWithConflict(w, storage.PgErrorsHandle(err, "Name"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (rs *Resolver) AdminEditService(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(r.RequestURI)
	if err != nil {
		panic(err)
	}
	values := u.Query()
	serviceNameInUrl := values.Get("service_name")

	r.ParseMultipartForm(10 << 20)

	serviceName := r.FormValue("service_name2")
	if serviceName != "" {
		if err := storage.EditAdminService(r.Context(), rs.App.Postgres, serviceNameInUrl, serviceName); err != nil {
			api_v1.RespondWithConflict(w, storage.PgErrorsHandle(err, "Service_name"))
			return
		}
	}

	file, handler, err := r.FormFile("file")
	if file != nil {
		serviceName = strings.ToLower(serviceName)
		if err != nil {
			fmt.Println("Error Retrieving the File")
			fmt.Println(err)
			return
		}
		defer file.Close()

		if !strings.HasSuffix(handler.Filename, ".svg") || handler.Header.Get("Content-Type") != "image/svg+xml" {
			api_v1.RespondWithBadRequest(w, "")
			return
		}

		path, err := tl.GetExecutablePath()
		dir := filepath.Join(path, "resources", "svg_files")
		if serviceName == "" {
			serviceName = serviceNameInUrl
		}
		fileName := strings.ReplaceAll(serviceName, " ", "-") + ".svg"
		fullPath := filepath.Join(dir, fileName)

		tempFile, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			api_v1.RespondWithBadRequest(w, "")
			return
		}
		defer tempFile.Close()
		io.Copy(tempFile, file)
	}

	if serviceName != "" {
		path, _ := tl.GetExecutablePath()
		oldFileName := strings.ReplaceAll(serviceNameInUrl, " ", "-") + ".svg"
		newFileName := strings.ReplaceAll(serviceName, " ", "-") + ".svg"
		oldPath := filepath.Join(path, "resources", "svg_files", oldFileName)
		newPath := filepath.Join(path, "resources", "svg_files", newFileName)

		if _, err := os.Stat(oldPath); err == nil {
			err = os.Rename(oldPath, newPath)
			if err != nil {
				fmt.Println("Error renaming the file: ", err)
				return
			}
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
