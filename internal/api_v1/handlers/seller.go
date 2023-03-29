package handlers

import "net/http"

func (rs *Resolver) SellerGetProducts(w http.ResponseWriter, r *http.Request)    {}
func (rs *Resolver) SellerCreateProduct(w http.ResponseWriter, r *http.Request)  {}
func (rs *Resolver) SellerProductsUpdate(w http.ResponseWriter, r *http.Request) {}
func (rs *Resolver) SellerProductsDelete(w http.ResponseWriter, r *http.Request) {}
