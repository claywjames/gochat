package main

import (
    "log"
    "os"
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/gorilla/mux"
)

var connections []group = make([]group, 10)

type msg struct {
    Message string
    Sender string
}

type chatClient struct {
    Account clientAccount
    Conn *websocket.Conn
}

type group struct {
    Name string
    Members []chatClient
}

var mangos group = group{"mangos", []chatClient{}}

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/login", loginPageHandler).Methods("POST")
    r.HandleFunc("/logout", logoutPageHandler).Methods("POST")
    r.HandleFunc("/signup", signupPageHandler).Methods("POST")

    r.HandleFunc("/websocket", websocketHandler)

    r.HandleFunc("/chat", chatHandler)

    r.HandleFunc("/", landingPageHandler)
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

    log.SetOutput(os.Stdout)
    log.Println("Listening on port 3000...")
    
    http.Handle("/", r)

    http.ListenAndServe(":3000", nil)
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    Subprotocols: []string{"name"},
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

func chatHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/chat.html")
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil) 
    if err != nil {
        log.Println(err)
        return
    }

    sender := getUsername(r)
    account, _ := getAccount(sender)
    newClient := chatClient{
        Account: account,
        Conn: conn,
    }

    mangos.Members = append(mangos.Members, newClient)

    for {
        message := msg{"", sender}
        if err := conn.ReadJSON(&message); err != nil {
            log.Println(err)
            break
        }
        mangos.broadcastMessage(message)
        //mangos.saveMessage(message msg)
    }
}

func (g * group) broadcastMessage(message msg) {
    for _, member := range g.Members {
        if err := member.Conn.WriteJSON(message); err != nil {
            log.Println(err)
        }
    }
}

// func (g * group) saveMessage(message msg) {
//     session, err := mgo.Dial("localhost")
//     if err != nil {
//         log.Println(err)
//     }
//     defer session.Close()

//     message = message.Sender + ": " + message.Message
//     c := session.DB("gochat").C("groups")
//}
