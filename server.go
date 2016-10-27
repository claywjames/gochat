package main

import (
    "log"
    "os"
    "net/http"
    "github.com/gorilla/mux"
)

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/login", loginPageHandler).Methods("POST")
    r.HandleFunc("/logout", logoutPageHandler).Methods("POST")
    r.HandleFunc("/signup", signupPageHandler).Methods("POST")
    r.HandleFunc("/createGroupPage", createGroupPageHandler)
    r.HandleFunc("/createGroup", groupCreationHandler).Methods("POST")

    r.HandleFunc("/websocket", messagingHandler)

    r.HandleFunc("/chat", chatPageHandler)

    r.HandleFunc("/", landingPageHandler)
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

    log.SetOutput(os.Stdout)
    log.Println("Listening on port 3000...")
    
    http.Handle("/", r)

    http.ListenAndServe(":3000", nil)
}

func createGroupPageHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/groupmaking.html")
}

func groupCreationHandler(w http.ResponseWriter, r *http.Request) {
    defer http.Redirect(w, r, "/chat", 302)
    groupName, groupMember := r.FormValue("groupName"), r.FormValue("groupMember")
    creator := getUsername(r)
    groupMemberAccount, err := getAccount(groupMember)
    if err != nil {
        log.Println(err)
        return
    }
    creatorAccount, _ := getAccount(creator)

    err = createGroup(groupName, []clientAccount{creatorAccount, groupMemberAccount})
    if err != nil {
        log.Println(err)
        return
    }
}

func landingPageHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/first.html")
}

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
    username, password := r.FormValue("username"), r.FormValue("password")
    redirectTarget := "/"
    if username != "" && password != "" {
        if validateAccount(username, password) {
            redirectTarget = "/chat"
            setSession(username, w)
        }
        http.Redirect(w, r, redirectTarget, 302)
    }
}

func logoutPageHandler(w http.ResponseWriter, r *http.Request) {
    clearSession(w)
    http.Redirect(w, r, "/", 302)
}

func signupPageHandler(w http.ResponseWriter, r *http.Request) {
    username, password := r.FormValue("username"), r.FormValue("password")
    err := createAccount(username, password)
    if err == nil {
        setSession(username, w)
        http.Redirect(w, r, "/chat", 302)
    } else {
        http.Redirect(w, r, "/", 302)
    }
}

func chatPageHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/chat.html")
}