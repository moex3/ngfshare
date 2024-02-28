package net

import (
    "fmt"
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"

    "git.fuwafuwa.moe/x3/ngfshare/config"
    "git.fuwafuwa.moe/x3/ngfshare/controller"
    "git.fuwafuwa.moe/x3/ngfshare/auth"
    "github.com/gorilla/mux"
)

func setupRoutes() *mux.Router {
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

    return r
}

func Start(conf config.Config) error {

    lstStr := fmt.Sprintf("%s:%d", conf.Address, conf.Port)
    srv := &http.Server{
        Addr: lstStr,
        Handler: setupRoutes(),
    }

    done := make(chan os.Signal, 1)
    signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

    go func() {
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            log.Printf("HTTP ListenAndServe error: %v\n", err)
        }
    }()
    log.Println("Listening on", lstStr)

    <-done
    fmt.Println("")
    log.Println("Stopping HTTP server")

    err := srv.Shutdown(context.Background())
    return err
}
