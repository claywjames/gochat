package main

import (
    "log"
    "os"
    "net/http"
    "github.com/gorilla/websocket"
)

var connections map[int]*websocket.Conn = make(map[int]*websocket.Conn)
var id int = 0

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
    connections[id] = conn;

    for {
        messageType, data, err := conn.ReadMessage()
        if err != nil {
            log.Println(err)
            break
        }
        if messageType == websocket.TextMessage {
            broadcastMessage(data)
        }

    }
}

func broadcastMessage(message []byte) {
    for _, conn := range connections {
        if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
            log.Println(err)
        }
    }
}

