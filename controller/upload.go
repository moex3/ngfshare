package controller

import (
    "fmt"
    "encoding/json"
    "io"
    "log"
    "net/http"
    "crypto/sha1"
    "os"
    "mime/multipart"

    "git.fuwafuwa.moe/x3/ngfshare/config"
    sauth "git.fuwafuwa.moe/x3/ngfshare/auth"
    "git.fuwafuwa.moe/x3/ngfshare/db"
)

type responseStruct struct {
    Id string `json:"id"`
    Filename string `json:"filename"`
    Url string `json:"url"`
    UrlShort string `json:"url_short"`
    DeleteUrl string `json:"delete_url"`
}
func responseWithFile(id, filename string, w http.ResponseWriter) {
    urlShort := fmt.Sprintf("%s/-%s", config.Conf.UrlPrefix, id)
    rsp := responseStruct{
        Id:         id,
        Filename:   filename,
        Url:        fmt.Sprintf("%s/%s", urlShort, filename),
        UrlShort:   urlShort,
        DeleteUrl:  fmt.Sprintf("%s/api/delete/%s", config.Conf.UrlPrefix, id),
    }
    jsb, err := json.Marshal(rsp)
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintln(w, `{"error":"Failed response json marshal"}`)
        return
    }
    w.Write(jsb)
}

func copyFileToStorage(srcFile multipart.File, sha1sum string) error {
    d1 := sha1sum[0:2]
    d2 := sha1sum[2:4]
    saveDir := fmt.Sprintf("%s/%s/%s", config.Conf.StoreDir, d1, d2)
    savePath := fmt.Sprintf("%s/%s", saveDir, sha1sum)
    
    err := os.MkdirAll(saveDir, 0750)
    if err != nil {
        log.Println("Cannot create directory for savedir:", saveDir, err)
        return err
    }

    dstFile, err := os.Create(savePath)
    if err != nil {
        log.Println("Cannot create file for savepath:", savePath, err)
        return err
    }
    defer dstFile.Close()

    _, err = io.Copy(dstFile, srcFile)
    if err != nil {
        log.Println("Cannot copy file", err)
        return err
    }

    log.Println("File copied")
    return nil
}

func Upload(w http.ResponseWriter, r *http.Request) {
    auth := r.Header.Get("Authorization")
    cookieAuth := false
    if auth == "" {
        auth = sauth.GetAuthCookie(r)
        cookieAuth = true
    }
    r.ParseMultipartForm(40*1024*1024)

    w.Header().Add("Content-Type", "application/json")

    file, fheader, err := r.FormFile("file")
    if err != nil {
        log.Println("No file POST field in upload, or some other error", err)
        http.Error(w, "Bad Request", http.StatusBadRequest)
        return
    }
    defer file.Close()

    sha := sha1.New()
    hlen, err := io.Copy(sha, file)
    if err != nil || hlen != fheader.Size {
        log.Println("Hash copy failed or not all file got hashed", err)
        http.Error(w, "Bad Request", http.StatusBadRequest)
    }
    sum := fmt.Sprintf("%x", sha.Sum(nil))

    dbFile, exists := db.Db.GetFileBySha1(sum)
    if exists {
        log.Println("Dupe upload detected on id", dbFile.Id)
        // Don't check if it was uploaded by the same api key or not because i don't care.
        if !cookieAuth {
            responseWithFile(dbFile.Id, dbFile.Filename, w)
        } else {
            // If upload from web, redirect to the final url
            http.Redirect(w, r, fmt.Sprintf("/-%s/%s", dbFile.Id, dbFile.Filename), http.StatusFound)
        }
        return
    }
    file.Seek(0, 0)
    b := make([]byte, 512)
    file.Read(b)
    fType := http.DetectContentType(b)
    if fType == "application/octet-stream" {
        hType := fheader.Header.Get("Content-Type")
        if hType != "" {
            fType = hType
        }
    }

    resId, tx, err := db.Db.InsertFile(fheader.Filename, fheader.Size, fType, sum, auth)
    if err != nil {
        log.Println("Failed to insert into database", err)
        w.WriteHeader(http.StatusInternalServerError)
        fmt.Fprintln(w, `{"error":"Failed to insert into databse"}`)
        return
    }

    file.Seek(0, 0)
    err = copyFileToStorage(file, sum)
    if err != nil {
        tx.Rollback()
        fmt.Fprintln(w, `{"error":"Failed to copy file"}`)
        return
    }
    tx.Commit()

    log.Printf("Added file with id: '%s' name: '%s' sha1sum: '%s'", resId, fheader.Filename, sum)
    if !cookieAuth {
        responseWithFile(resId, fheader.Filename, w)
    } else {
        // If upload from web, redirect to the final url
        http.Redirect(w, r, fmt.Sprintf("/-%s/%s", resId, fheader.Filename), http.StatusFound)
    }
}
