package handlers_v1

import (
	"net/http"
)

func (rs *Resolver) UserProfileDelete(w http.ResponseWriter, r *http.Request) {
	//token, data, err := api_v1.ContextGetAuthenticated(r)
	//if err != nil {
	//	rs.Client.Logger.NewWarn("error in took jwt data", err)
	//	api_v1.RespondWithInternalServerError(w)
	//	return
	//}

}

func (rs *Resolver) UserProfileDump(w http.ResponseWriter, r *http.Request)   {}
func (rs *Resolver) UserProfileUpdate(w http.ResponseWriter, r *http.Request) {}
func (rs *Resolver) UserProfileOrders(w http.ResponseWriter, r *http.Request) {}
