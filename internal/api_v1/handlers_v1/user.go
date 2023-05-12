package handlers_v1

import (
	"encoding/json"
	"net/http"
	"test-server-go/internal/api_v1"
	freekassa2 "test-server-go/internal/freekassa"
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
	if err := tl.Validate(data.VariantId, tl.IsNotBlank(), tl.IsLen(UUIDLength), tl.IsNotContainsSpace(), tl.IsValidUUID(), tl.IsTrimmedSpace()); err != nil {
		api_v1.RespondWithUnprocessableEntity(w, "VariantId: "+err.Error())
		return
	}
	var coupon string
	if data.Coupon != nil && *data.Coupon != "" {
		coupon = *data.Coupon
		if err := tl.Validate(coupon, tl.IsNotBlank(), tl.IsMinMaxLen(MinCouponLength, MaxCouponLength), tl.IsNotContainsSpace(), tl.IsTrimmedSpace()); err != nil {
			api_v1.RespondWithUnprocessableEntity(w, "Email: "+err.Error())
			return
		}
	}

	// Block 2 - create payment url and check on access
	productId, variantName, variantState, quantityCurrent, finalPrice, err := storage.GetProductVariantForPayment(r.Context(), rs.App.Postgres, data.VariantId)
	if err != nil {
		rs.App.Logger.NewWarn("error in get product variant for payment", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}
	if productId == "" || variantState == storage.ProductStateInvisible {
		api_v1.RedRespond(w, http.StatusNotFound, "Not found", "This product variant not found")
		return
	}

	if variantState == storage.ProductStateDeleted {
		api_v1.RedRespond(w, http.StatusForbidden, "Forbidden", "This product variant has been deleted")
		return
	} else if variantState == storage.ProductStateUnavailableWithoutPrice {
		api_v1.RedRespond(w, http.StatusForbidden, "Forbidden", "This product variant unavailable without price")
		return
	} else if variantState == storage.ProductStateUnavailableWithPrice {
		api_v1.RedRespond(w, http.StatusForbidden, "Forbidden", "This product variant unavailable with price")
		return
	}

	if quantityCurrent == 0 {
		api_v1.RedRespond(w, http.StatusForbidden, "Forbidden", "This product variant sold out")
		return
	}

	_, jwtData, err := api_v1.ContextGetAuthenticated(r)
	if err != nil {
		rs.App.Logger.NewWarn("error in took jwt data", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	orderId, err := storage.CreateUserOrder(r.Context(), rs.App.Postgres, data.VariantId, jwtData.AccountUuid, finalPrice)
	if err != nil {
		rs.App.Logger.NewWarn("error in create user order", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	url := freekassa2.NewOrderUrl(rs.App.Freekassa, finalPrice, freekassa2.CurrencyRUB, variantName+"_"+orderId)

	// Block 3 - send the result
	response := struct {
		PaymentUrl string `json:"payment_url"`
	}{
		PaymentUrl: url,
	}

	api_v1.RespondWithCreated(w, response)
}

func (rs *Resolver) UserProfileDump(w http.ResponseWriter, r *http.Request)   {}
func (rs *Resolver) UserProfileUpdate(w http.ResponseWriter, r *http.Request) {}
func (rs *Resolver) UserProfileOrders(w http.ResponseWriter, r *http.Request) {}
