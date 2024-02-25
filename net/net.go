package net

import (
    "fmt"
    "net/http"

    "git.fuwafuwa.moe/x3/ngfshare/config"
    "git.fuwafuwa.moe/x3/ngfshare/controller"
    "git.fuwafuwa.moe/x3/ngfshare/auth"
    "github.com/gorilla/mux"
)

func Start(conf config.Config) error {
    r := mux.NewRouter()

    authedR := r.PathPrefix("/api").Methods("POST").Subrouter()
    authedR.Use(auth.AuthMiddleware)
    authedR.HandleFunc("/upload", controller.Upload)
    authedR.HandleFunc("/delete/{id}", controller.Delete)

    r.HandleFunc("/-{id}", controller.Download).Methods("GET")
    r.HandleFunc("/-{id}/{filename}", controller.Download).Methods("GET")
    r.HandleFunc("/-{id}/", controller.Download).Methods("GET")

    r.HandleFunc("/", controller.WebGet).Methods("GET")
    r.HandleFunc("/login", controller.WebLogin).Methods("POST")
    r.HandleFunc("/logout", controller.WebLogout).Methods("POST")

    http.Handle("/", r)

    lstStr := fmt.Sprintf("%s:%d", conf.Address, conf.Port)
    fmt.Println("Listening on", lstStr)
    http.ListenAndServe(lstStr, nil)

    return nil
}
