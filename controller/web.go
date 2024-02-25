package controller

import (
    "net/http"
    "log"

    "git.fuwafuwa.moe/x3/ngfshare/view"
    "git.fuwafuwa.moe/x3/ngfshare/db"
    "git.fuwafuwa.moe/x3/ngfshare/auth"
)

func webGetList(w http.ResponseWriter, r *http.Request, auth string) {
    log.Println("In webget file list")

    files, err := db.Db.GetFilesByAuthKey(auth)
    if err != nil {
        log.Println("Failed to get files by auth key", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    err = view.Execute("files", w, files)
    if err != nil {
        log.Println("Error in webGetList:", err)
    }
}

func webGetLogin(w http.ResponseWriter, r *http.Request) {
    log.Println("In webget login")
    err := view.Execute("login", w, "")
    if err != nil {
        log.Println("Error in webGetLogin:", err)
    }
}

func WebGet(w http.ResponseWriter, r *http.Request) {
    auth := auth.GetAuthCookie(r)
    authOk := auth != "" && db.Db.IsAuthKeyExists(auth)

    if authOk {
        webGetList(w, r, auth)
    } else {
        webGetLogin(w, r)
    }
}

func WebLogin(w http.ResponseWriter, r *http.Request) {
    auth := r.FormValue("auth")
    ok := db.Db.IsAuthKeyExists(auth)
    if !ok {
        http.Error(w, "Authentication failure", http.StatusForbidden)
        return
    }

    log.Println("Successfull web auth for key", auth)
    cookie := http.Cookie{
        Name: "auth",
        Value: auth,
    }
    http.SetCookie(w, &cookie)
    http.Redirect(w, r, "/", http.StatusFound)
}

func WebLogout(w http.ResponseWriter, r *http.Request) {
    log.Println("In logout")
    cookie := http.Cookie{
        Name: "auth",
        Value: "",
        MaxAge: -1,
    }
    http.SetCookie(w, &cookie)
    http.Redirect(w, r, "/", http.StatusFound)
}
