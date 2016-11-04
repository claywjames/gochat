package main

import (
    "log"
    "net/http"
    "github.com/gorilla/websocket"
    "gopkg.in/mgo.v2"
)

type msg struct {
    Message string
    Sender string
}

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

func messagingHandler(w http.ResponseWriter, r *http.Request) {
    sender := getUsername(r)
    account, _ := getAccount(sender)
    if len(account.Groups) > 0 {
        conn, err := upgrader.Upgrade(w, r, nil) 
        if err != nil {
            log.Println(err)
            return
        }

        activeGroup := account.Groups[0]

        go func() {
            for {
                message := msg{"", sender}
                if err := conn.ReadJSON(&message); err != nil {
                    log.Println(err)
                    break
                }
                activeGroup.saveMessage(message)
            }
        }()

        go func() {
            session, err := mgo.Dial("localhost")
            if err != nil {
                log.Println(err)
            }
            defer session.Close()

            c := session.DB("gochat").C(activeGroup.Name)
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
}

func (g * group) saveMessage(message msg) {
    session, err := mgo.Dial("localhost")
    if err != nil {
        log.Println(err)
    }
    defer session.Close()

    c := session.DB("gochat").C(g.Name)
    err = c.Insert(&message)
    if err != nil {
        log.Println(err)
    }
}
