package handlers_v1

import (
	"net/http"
	"test-server-go/internal/api_v1"
)

func (rs *Resolver) FreekassaNotification(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		api_v1.RespondWithBadRequest(w, "")
		return
	}

	merchantID := r.FormValue("MERCHANT_ID")
	//amount := r.FormValue("AMOUNT")
	intid := r.FormValue("intid")
	//merchantOrderID := r.FormValue("MERCHANT_ORDER_ID")
	//pEmail := r.FormValue("P_EMAIL")
	//pPhone := r.FormValue("P_PHONE")
	//curID := r.FormValue("CUR_ID")
	//sign := r.FormValue("SIGN")
	//payerAccount := r.FormValue("payer_account")
	//commission := r.FormValue("commission")

	//if sign != freekassa.CreateNotificationSignature(rs.App.Freekassa, amount, merchantOrderID) {
	//	rs.App.Logger.NewWarn("error signature verification", nil)
	//	//api_v1.RespondWithBadRequest(w, "")
	//	//return
	//}

	//res := fmt.Sprintf(
	//	"заказ с номером %s был успешно оплачен. Сумма: %s, email: %s, телефон: %s, ID электронной валюты: %s, номер счета/карты плательщика: %s, комиссия: %s, id магазина: %s, номер операции Free-Kassa: %s\n",
	//	merchantOrderID, amount, pEmail, pPhone, curID, payerAccount, commission, merchantID, intid)

	rs.App.Logger.NewInfo("URI: " + r.RequestURI + " IP: " + r.RemoteAddr + " MERCHANT_ID: " + merchantID + " intid: " + intid)

	w.WriteHeader(http.StatusNoContent)
}
