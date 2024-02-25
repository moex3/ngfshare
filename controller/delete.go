package controller

import (
    "net/http"
    "fmt"
    "log"
    "os"

    "git.fuwafuwa.moe/x3/ngfshare/db"
    "git.fuwafuwa.moe/x3/ngfshare/config"
    sauth "git.fuwafuwa.moe/x3/ngfshare/auth"
    "github.com/gorilla/mux"
)

func deleteFile(sha1sum string) error {
    d1 := sha1sum[0:2]
    d2 := sha1sum[2:4]
    path := fmt.Sprintf("%s/%s/%s/%s", config.Conf.StoreDir, d1, d2, sha1sum)

    // This leaves empty directories, but whatever
    err := os.Remove(path)
    return err
}

func Delete(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    auth := r.Header.Get("Authorization")
    if auth == "" {
        auth = sauth.GetAuthCookie(r)
    }
    w.Header().Add("Content-Type", "application/json")

    file, err := db.Db.GetFileById(id)
    if err != nil {
        log.Println("Delete: File by id", id, "not found in databse", err)
        http.Error(w, `{"status":"Not found"}`, http.StatusNotFound)
        return
    }

    if file.UploadKey != auth {
        log.Println("Trying to delete a file uploaded by a different key", id)
        http.Error(w, "", http.StatusForbidden)
        return
    }

    ex, err := db.Db.DeleteFile(id)
    if err != nil {
        log.Println("Failed to delete file with id", id, err)
        http.Error(w, "", http.StatusInternalServerError)
        return
    }

    if !ex {
        log.Println("Tried to delete file that doens't exists", id)
        http.NotFound(w, r)
        return
    }

    err = deleteFile(file.Sha1Sum)
    if err != nil {
        log.Println("Cannot delete file with sum", file.Sha1Sum, err)
    }

    log.Println("Deleted file", id)

    /*
    if auth.GetAuthCookie() != "" {
        // Called from the webpage, return redirect
    }
    */

    fmt.Fprintln(w, `{"status":"OK"}`)
}
