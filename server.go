package main

import (
    "log"
    "os"
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/gorilla/mux"
    "github.com/gorilla/securecookie"
    "errors"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

var connections []group = make([]group, 10)

var cookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))

type msg struct {
    Message string
}

type clientAccount struct {
    Username string
    Password string
    Groups []group
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

func setSession(username string, w http.ResponseWriter) {
    value := map[string]string{
        "username": username,
    }
    if encoded, err := cookieHandler.Encode("session", value); err == nil {
        cookie := &http.Cookie{
            Name:  "session",
            Value: encoded,
            Path:  "/",
        }
        http.SetCookie(w, cookie)
    }
}

func logoutPageHandler(w http.ResponseWriter, r *http.Request) {
    clearSession(w)
    http.Redirect(w, r, "/", 302)
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

func createAccount(username string, password string) error {
    if _, err := getAccount(username); err == nil {
        return errors.New("username taken")
    }
    session, err := mgo.Dial("localhost")
    if err != nil {
        log.Println(err)
    }
    defer session.Close()

    c := session.DB("plaintext").C("accounts")
    groups := []group{}
    err = c.Insert(&clientAccount{username, password, groups})
    if err != nil {
        log.Println(err)
    }
    return nil
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
        message := msg{}
        if err := conn.ReadJSON(&message); err != nil {
            log.Println(err)
            break
        }
        mangos.broadcastMessage(sender, []byte(message.Message))

    }
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

func (g * group) broadcastMessage(name string, message []byte) {
    message = []byte(name + ": " + string(message))
    for _, member := range g.Members {
        if err := member.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
            log.Println(err)
        }
    }
}

func validateAccount(username string, password string) bool {
    account, err := getAccount(username)
    if err != nil {
        return false
    }
    if account.Password != password {
        return false
    }
    return true
}

func getAccount(username string) (account clientAccount, err error) {
    session, err := mgo.Dial("localhost")
    if err != nil {
        log.Println(err)
    }
    defer session.Close()

    c := session.DB("plaintext").C("accounts")

    err = c.Find(bson.M{"username":username}).One(&account)
    return 
}
