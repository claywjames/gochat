package main

import (
    "log"
    "os"
    "net/http"
    "github.com/gorilla/websocket"
    "errors"
)

var connections []group = make([]group, 10)

type msg struct {
    MsgType string
    Username string
    Password string
    Message string
}

type clientAccount struct {
    username string
    password string
    groups []group
}

type chatClient struct {
    account clientAccount
    conn *websocket.Conn
}

type group struct {
    name string
    members []chatClient
}

var mangos group = group{"mangos", []chatClient{}}
var clay clientAccount = clientAccount{
    username: "clay",
    password: "1234",
    groups: []group{mangos},
}
var users []clientAccount = []clientAccount{clay}

func main() {
    fileServer := http.FileServer(http.Dir("static"))
    http.Handle("/", fileServer)

    http.HandleFunc("/websocket", websocketHandler)

    log.SetOutput(os.Stdout)
    log.Println("Listening on port 3000...")
    
    http.ListenAndServe(":3000", nil)
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    Subprotocols: []string{"name"},
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
    var newClient chatClient;
    conn, err := upgrader.Upgrade(w, r, nil) 
    if err != nil {
        log.Println(err)
        return
    }

    message := msg{}
    if err := conn.ReadJSON(&message); err != nil {
        log.Println(err)
        return
    }

    if message.MsgType == "LOGIN" {
        if account, err := getAccount(message.Username, message.Password); err != nil {
            log.Println(err)
        } else {
            newClient = chatClient{account, conn}
            log.Println(message.Username + " connected")
        }
    }

    mangos.members = append(mangos.members, newClient)

    for {
        text := msg{}
        if err := conn.ReadJSON(&text); err != nil {
            log.Println(err)
            break
        }
        mangos.broadcastMessage(newClient.account.username, []byte(text.Message))

    }
}

func (g * group) broadcastMessage(name string, message []byte) {
    message = []byte(name + ": " + string(message))
    for _, member := range g.members {
        if err := member.conn.WriteMessage(websocket.TextMessage, message); err != nil {
            log.Println(err)
        }
    }
}

func getAccount(username string, password string) (clientAccount, error) {
    for _, user := range users {
        if username == user.username && password == user.password {
            return user, nil
        }
    }
    return clientAccount{"", "", []group{}}, errors.New("No Account Found")
}
