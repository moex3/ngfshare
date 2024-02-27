package controller

import (
    "net/http"
    "fmt"
    "log"

    "git.fuwafuwa.moe/x3/ngfshare/db"
    "github.com/gorilla/mux"
)

func Download(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    file, err := db.Db.GetFileById(id)
    if err != nil {
        log.Println("File by id", id, "not found in databse", err)
        http.NotFound(w, r)
        return
    }

    ifModifH := r.Header.Get("If-Modified-Since")
    if ifModifH != "" {
        // Is (probably) not modified, send back 304
        log.Printf("Download (%s): got If-Modified-Since, replying 304\n", id)
        w.WriteHeader(http.StatusNotModified)
        return
    }

    d1 := file.Sha1Sum[0:2]
    d2 := file.Sha1Sum[2:4]

    w.Header().Add("Content-Disposition", fmt.Sprintf("filename=\"%s\"", file.Filename))
    w.Header().Add("Content-Type", file.ContentType)
    w.Header().Add("X-Accel-Expires", "1800")
    w.Header().Add("X-Accel-Redirect", fmt.Sprintf("/store/%s/%s/%s", d1, d2, file.Sha1Sum))

    log.Printf("Served file with id: '%s'  sha1sum: '%s'", file.Id, file.Sha1Sum)
}
