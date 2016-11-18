package main

import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/gorilla/mux"
    "gopkg.in/mgo.v2"
    "time"
    "os"
)

type msg struct {
    Message string
    Sender string
    TimeStamp string
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func messagingHandler(w http.ResponseWriter, r *http.Request) {
    sender := getUsername(r)
    activeGroup := mux.Vars(r)["group"]

    conn, err := upgrader.Upgrade(w, r, nil) 
    if err != nil {
        log.Println(err)
        return
    }

    go func() {
        for {
            message := msg{}
            if err := conn.ReadJSON(&message); err != nil {
                log.Println(err)
                break
            }
            message.Sender = sender
            message.TimeStamp = time.Now().Format(time.Stamp)
            saveMessage(message, activeGroup)
        }
    }()

    go func() {
        uri := os.Getenv("MONGODB_URI")
        if uri == "" {
            uri = "localhost"
        }
        session, err := mgo.Dial(uri)
        if err != nil {
            log.Println(err)
        }
        defer session.Close()

        c := session.DB("gochat").C(activeGroup)
        newMessage := c.Find(nil).Tail(-1)
        defer newMessage.Close()

        message := msg{};
        for {
            if !newMessage.Next(&message) {
                break
            }
            if err := conn.WriteJSON(message); err != nil {
                log.Println(err)
            }
        }
    }()
}

func saveMessage(message msg, group string) {
    uri := os.Getenv("MONGODB_URI")
    if uri == "" {
        uri = "localhost"
    }
    session, err := mgo.Dial(uri)
    if err != nil {
        log.Println(err)
    }
    defer session.Close()

    c := session.DB("gochat").C(group)
    err = c.Insert(&message)
    if err != nil {
        log.Println(err)
    }
}
