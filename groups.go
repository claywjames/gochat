package main

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "errors"
    "log"
)

type group struct {
    Name string
    Members []clientAccount
}

func createGroup(name string, members []clientAccount) error {
    session, err := mgo.Dial("localhost")
    if err != nil {
        log.Println(err)
    }
    defer session.Close()

    c := session.DB("gochat").C("groups")
    err = c.Find(bson.M{"name" : name}).One(nil)
    if err == nil {
        return errors.New("Group Name Taken")
    }

    newGroup := group{name, members}
    err = c.Insert(&newGroup)
    if err != nil {
        return err
    }

    groupMessageCollection := session.DB("gochat").C(name)
    collectionInfo := &mgo.CollectionInfo{
        Capped: true,
        MaxBytes: 10000000000,
        MaxDocs: 1000,
    }

    err = groupMessageCollection.Create(collectionInfo)
    if err != nil {
        return err
    }

    c = session.DB("gochat").C("accounts")
    for _, member := range members {
        member.Groups = append(member.Groups, newGroup)
        err = c.Update(bson.M{"username": member.Username}, bson.M{"$set": bson.M{"groups": member.Groups}})
        if err != nil {
            return err
        }
    }
    return nil
}
