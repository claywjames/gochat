package main

import (
    "net/http"
    "github.com/gorilla/securecookie"
)

var cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))

func setSession(username string, w http.ResponseWriter) {
    value := map[string]string{
        "username": username,
    }
    if encoded, err := cookieHandler.Encode("session", value); err == nil {
        cookie := &http.Cookie{
            Name:  "session",
            Value: encoded,
            Path:  "/",
           // Secure: true,
        }
        http.SetCookie(w, cookie)
    }
}

func clearSession(w http.ResponseWriter) {
    cookie := &http.Cookie{
        Name: "session",
        Value: "",
        Path: "/",
        MaxAge: -1,
    }
    http.SetCookie(w, cookie)
}

func getUsername(r *http.Request) (username string) {
    if cookie, err := r.Cookie("session"); err == nil {
        cookieValue := make(map[string]string)
        if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
            username = cookieValue["username"]
        }
    }
    return username
}

func setBadLoginCookie(w http.ResponseWriter) {
    http.SetCookie(w, &http.Cookie{Name: "login", Value: "Failed", Path: "/",})
}

func setFailedSignUpCookie(w http.ResponseWriter) {
    http.SetCookie(w, &http.Cookie{Name: "signup", Value: "Failed", Path: "/",})
}

func setFailedGroupCreationCookie(w http.ResponseWriter, reason string) {
    http.SetCookie(w, &http.Cookie{Name: "creation", Value: reason, Path: "/",})
}