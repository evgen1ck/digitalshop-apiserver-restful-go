package handlers_v1

import (
	"net/http"
	"strings"
	"test-server-go/internal/api_v1"
	"test-server-go/internal/storage"
)

func (rs *Resolver) FreekassaNotification(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}

	merchantID := r.FormValue("MERCHANT_ORDER_ID")
	//amount := r.FormValue("AMOUNT")
	//intid := r.FormValue("intid")
	//merchantOrderID := r.FormValue("MERCHANT_ORDER_ID")
	//pEmail := r.FormValue("P_EMAIL")
	//pPhone := r.FormValue("P_PHONE")
	//curID := r.FormValue("CUR_ID")
	//sign := r.FormValue("SIGN")
	//payerAccount := r.FormValue("payer_account")
	//commission := r.FormValue("commission")

	splitID := strings.Split(merchantID, "_")
	if len(splitID) < 2 {
		rs.App.Logger.NewWarn("error in split merchantID", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	email, nickname, content, productName, variantName, serviceName, itemName, err := storage.GetDataForFreekassa(r.Context(), rs.App.Postgres, splitID[len(splitID)-1])
	if err != nil {
		rs.App.Logger.NewWarn("error in get data for freekassa", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	err = rs.App.Mailer.SendOrderContent(email, nickname, productName+" - "+variantName, serviceName, itemName, content, rs.App.Config.App.Service.Url.Client)
	if err != nil {
		rs.App.Logger.NewWarn("error in send order content", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
