package handlers

import "net/http"

func (rs *Resolver) ServerPostgresBackup(w http.ResponseWriter, r *http.Request) {}
func (rs *Resolver) AdminGamesAdd(w http.ResponseWriter, r *http.Request)        {}
func (rs *Resolver) AdminGamesUpdate(w http.ResponseWriter, r *http.Request)     {}
func (rs *Resolver) AdminGamesDelete(w http.ResponseWriter, r *http.Request)     {}
func (rs *Resolver) AdminOthersAdd(w http.ResponseWriter, r *http.Request)       {}
func (rs *Resolver) AdminOthersUpdate(w http.ResponseWriter, r *http.Request)    {}
func (rs *Resolver) AdminOthersDelete(w http.ResponseWriter, r *http.Request)    {}
