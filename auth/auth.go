package auth

import (
    "log"
    "net/http"

    "git.fuwafuwa.moe/x3/ngfshare/db"
)

func GetAuthCookie(r *http.Request) string {
    cookies := r.Cookies()
    for i, _ := range(cookies) {
        if cookies[i].Name == "auth" {
            return cookies[i].Value
        }
    }
    return ""
}

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        auth := r.Header.Get("Authorization")
        if auth == "" {
            auth = GetAuthCookie(r)
        }

        if auth != "" && db.Db.IsAuthKeyExists(auth) {
            log.Printf("Accepted auth header: '%s'\n", auth)
            next.ServeHTTP(w, r)
        } else {
            log.Printf("Rejected auth header: '%s'\n", auth)
            http.Error(w, "Forbidden", http.StatusForbidden)
        }
    })
}
