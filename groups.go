package main

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/dchest/uniuri"
    "errors"
    "log"
    "time"
    "os"
)

type group struct {
    Name string
    Members []clientAccount
    Code string
}

func getGroup(name string) (g group, err error) {
    uri := os.Getenv("MONGODB_URI")
    if uri == "" {
        uri = "localhost"
    }
    session, err := mgo.Dial(uri)
    if err != nil {
        log.Println(err)
    }
    defer session.Close()

    c := session.DB("heroku_jhn2m29z").C("groups")
    err = c.Find(bson.M{"name" : name}).One(&g)
    
    return
}

func getGroupFromCode(code string) (g group, err error) {
    uri := os.Getenv("MONGODB_URI")
    if uri == "" {
        uri = "localhost"
    }
    session, err := mgo.Dial(uri)
    if err != nil {
        log.Println(err)
    }
    defer session.Close()

    c := session.DB("heroku_jhn2m29z").C("groups")
    err = c.Find(bson.M{"code" : code}).One(&g)   
    
    return
}

func createGroup(name string, members []clientAccount) error {
    uri := os.Getenv("MONGODB_URI")
    if uri == "" {
        uri = "localhost"
    }

    _, err := getGroup(name)
    if err == nil {
        return errors.New("Group Name Taken")
    }

    session, err := mgo.Dial(uri)
    if err != nil {
        log.Println(err)
    }
    defer session.Close()

    newGroup := group{name, members, uniuri.NewLen(8)}
    c := session.DB("heroku_jhn2m29z").C("groups")
    err = c.Insert(&newGroup)
    if err != nil {
        return err
    }

    groupMessageCollection := session.DB("heroku_jhn2m29z").C(name)
    collectionInfo := &mgo.CollectionInfo{
        Capped: true,
        MaxBytes: 512000,
        MaxDocs: 100,
    }
    err = groupMessageCollection.Create(collectionInfo)
    if err != nil {
        return err
    }

    c = session.DB("heroku_jhn2m29z").C(name)
    err = c.Insert(&msg{"I have created " + name, members[0].Username, time.Now().Format(time.Stamp)})

    c = session.DB("heroku_jhn2m29z").C("accounts")
    for _, member := range members {
        member.Groups = append(member.Groups, newGroup)
        err = c.Update(bson.M{"username": member.Username}, bson.M{"$set": bson.M{"groups": member.Groups}})
        if err != nil {
            return err
        }
    }
    return nil
}

func (g * group) addGroupMember(joiner clientAccount) (err error) {
    uri := os.Getenv("MONGODB_URI")
    if uri == "" {
        uri = "localhost"
    }
    session, err := mgo.Dial(uri)
    if err != nil {
        log.Println(err)
    }
    defer session.Close()

    g.Members = append(g.Members, joiner)
    c := session.DB("heroku_jhn2m29z").C("groups")
    err = c.Update(bson.M{"name" : g.Name}, bson.M{"$set": bson.M{"members": g.Members}})
    if err != nil {
        return
    }

    return
}
