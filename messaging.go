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

    quit := make(chan int)
    go func(quit chan int) {
        for {
            select {
            case <- quit:
                return
            default:
                message := msg{}
                if err := conn.ReadJSON(&message); err != nil {
                    log.Println(err)
                    conn.Close()
                    quit <- 0
                    break
                }
                message.Sender = sender
                message.TimeStamp = time.Now().UTC().Format(time.RFC3339)
                saveMessage(message, activeGroup)
            }
        }
    }(quit)

    go func(quit chan int) {
        uri := os.Getenv("MONGODB_URI")
        if uri == "" {
            uri = "localhost"
        }
        session, err := mgo.Dial(uri)
        if err != nil {
            log.Println(err)
        }
        defer session.Close()

        c := session.DB("heroku_jhn2m29z").C(activeGroup)
        newMessage := c.Find(nil).Tail(-1)
        defer newMessage.Close()

        message := msg{};
        for {
            select {
            case <- quit:
                return
            default:
                if !newMessage.Next(&message) {
                break
                }
                if err := conn.WriteJSON(message); err != nil {
                    log.Println(err)
                }
            }
        }
    }(quit)

    go func(quit chan int) {
        for range time.Tick(time.Second * 30){
            select {
            case <- quit:
                return
            default:
                conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(15 * time.Second))
            }
        }
    }(quit)
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

    c := session.DB("heroku_jhn2m29z").C(group)
    err = c.Insert(&message)
    if err != nil {
        log.Println(err)
    }
}
