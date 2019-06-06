package api

import "net/http"

type API struct {
}

func New() *API {
	return &API{}
}

func (a *API) root(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello, maco"))
}

func (a *API) Register(mux *http.ServeMux) {
	mux.Handle("/characters", http.HandlerFunc(a.root))
}
