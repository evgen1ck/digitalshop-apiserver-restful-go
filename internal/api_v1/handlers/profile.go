package handlers

import "net/http"

func (rs *Resolver) ProfileData(w http.ResponseWriter, r *http.Request)   {}
func (rs *Resolver) ProfileDump(w http.ResponseWriter, r *http.Request)   {}
func (rs *Resolver) ProfileUpdate(w http.ResponseWriter, r *http.Request) {}
func (rs *Resolver) ProfileDelete(w http.ResponseWriter, r *http.Request) {}

func (rs *Resolver) ProfileOrders(w http.ResponseWriter, r *http.Request) {}
func (rs *Resolver) ProfileOrder(w http.ResponseWriter, r *http.Request)  {}
