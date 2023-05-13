package handlers_v1

import (
	"github.com/google/uuid"
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
	////amount := r.FormValue("AMOUNT")
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

	rs.App.Logger.NewInfo("abc")

	splitID := strings.Split(merchantID, "_")

	splitID[0] = strings.ReplaceAll(splitID[0], "-", " ")
	rs.App.Logger.NewInfo("abc1")

	email, err := storage.GetDataForFreekassa(r.Context(), rs.App.Postgres, splitID[1])
	if err != nil {
		rs.App.Logger.NewWarn("error in get user email for freekassa", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}
	rs.App.Logger.NewInfo("abc2")

	uuidd, _ := uuid.NewRandom()

	err = rs.App.Mailer.SendOrderContent(email, splitID[0], uuidd.String(), rs.App.Config.App.Service.Name, rs.App.Config.App.Service.Url.Client)
	if err != nil {
		rs.App.Logger.NewWarn("error in send order content", err)
		api_v1.RespondWithInternalServerError(w)
		return
	}
	rs.App.Logger.NewInfo("abc3")

	rs.App.Logger.NewInfo("URI: " + r.RequestURI + " IP: " + r.RemoteAddr + " MERCHANT_ID: " + merchantID + " intid: " + intid)

	rs.App.Logger.NewInfo("abc4")
	w.WriteHeader(http.StatusNoContent)
}

func (rs *Resolver) FreekassaOK(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, rs.App.Config.App.Service.Url.Client+"/finish", http.StatusSeeOther)
	w.WriteHeader(http.StatusNoContent)
}
