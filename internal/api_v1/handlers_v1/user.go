package handlers_v1

import (
	"encoding/json"
	"net/http"
	"test-server-go/internal/api_v1"
	freekassa "test-server-go/internal/freekassa"
	"test-server-go/internal/storage"
	tl "test-server-go/internal/tools"
)

func (rs *Resolver) UserProfileDelete(w http.ResponseWriter, r *http.Request) {
	//token, data, err := api_v1.ContextGetAuthenticated(r)
	//if err != nil {
	//	rs.Client.Logger.NewWarn("error in took jwt data", err)
	//	api_v1.RespondWithInternalServerError(w)
	//	return
	//}
}

func (rs *Resolver) UserNewPayment(w http.ResponseWriter, r *http.Request) {
	// Block 0 - decode data
	var data struct {
		VariantId string  `json:"variant_id"`
		Coupon    *string `json:"coupon"`
	}
	decodeErr := json.NewDecoder(r.Body).Decode(&data)
	if decodeErr != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}

	// Block 1 - data validation
	if err := tl.Validate(data.VariantId, tl.IsNotBlank(true), tl.IsLen(UUIDLength), tl.IsNotContainsSpace(), tl.IsValidUUID(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "VariantId: "+err.Error())
		return
	}
	var coupon string
	if data.Coupon != nil && *data.Coupon != "" {
		coupon = *data.Coupon
		if err := tl.Validate(coupon, tl.IsNotBlank(true), tl.IsMinMaxLen(MinCouponLength, MaxCouponLength), tl.IsNotContainsSpace(), tl.IsTrimmedSpace()); err != nil {
			api_v1.RespondWithUnprocessableEntity(w, "Email: "+err.Error())
			return
		}
	}

	_, jwtData, err := api_v1.ContextGetAuthenticated(r)
	if err != nil {
		rs.App.Logger.NewWarn("error in took jwt data", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	// Block 2 - create payment url and check on access
	orderId, variantName, finalPrice, err := storage.CreateOrder(r.Context(), rs.App.Postgres, jwtData.AccountUuid, data.VariantId)
	if err != nil {
		rs.App.Logger.NewWarn("error in get product variant for payment", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	url := freekassa.NewOrderUrl(rs.App.Freekassa, finalPrice, freekassa.CurrencyRUB, variantName+"_"+orderId)

	// Block 3 - send the result
	response := struct {
		PaymentUrl string `json:"payment_url"`
	}{
		PaymentUrl: url,
	}

	api_v1.RespondWithCreated(w, response)
}

func (rs *Resolver) UserProfileOrders(w http.ResponseWriter, r *http.Request) {
	_, jwtData, err := api_v1.ContextGetAuthenticated(r)
	if err != nil {
		rs.App.Logger.NewWarn("error in took jwt data", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	orders, err := storage.GetUserOrders(r.Context(), rs.App.Postgres, jwtData.AccountUuid)
	if err != nil {
		rs.App.Logger.NewWarn("error in get orders", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	api_v1.RespondOK(w, orders)
}

func (rs *Resolver) UserProfileDump(w http.ResponseWriter, r *http.Request)   {}
func (rs *Resolver) UserProfileUpdate(w http.ResponseWriter, r *http.Request) {}
