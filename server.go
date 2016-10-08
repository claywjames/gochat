package main

import (
    "log"
    "fmt"
    "os"
    "net/http"
    "github.com/gorilla/websocket"
)

var connections []group = make([]group, 10)
var id int = 0

type chatClient struct {
    name int
    conn *websocket.Conn
}

type group struct {
    name string
    members []chatClient
}

var mangos group = group{"mangos", []chatClient{}}

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
}

func websocketHandler(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil) 
    if err != nil {
        log.Println(err)
        return
    }
    log.Println("Client connected")

    id++
    newClient := chatClient{id, conn}
    mangos.members = append(mangos.members, newClient)

    for {
        messageType, data, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            break
        }
        if messageType == websocket.TextMessage {
            log.Println("Before broadcastMessage")
            mangos.broadcastMessage(newClient.name, data)
        }

    }
}

func (g * group) broadcastMessage(name int, message []byte) {
    id := fmt.Sprintf("%d", name)
    namedMessage := []byte(id + ": " + string(message))
    for _, member := range g.members {
        log.Println(member)
        if err := member.conn.WriteMessage(websocket.TextMessage, namedMessage); err != nil {
            log.Println(err)
        }
    }
}

