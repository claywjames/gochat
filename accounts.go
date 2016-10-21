package main

import (
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "errors"
    "log"
    "io"
    "crypto/rand"
    "golang.org/x/crypto/scrypt"
)

type clientAccount struct {
    Username string
    Password []byte
    HashSalt []byte
    Groups []group
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

    salt := make([]byte, 32)
    _, err = io.ReadFull(rand.Reader, salt)
    if err != nil {
        log.Fatal(err)
    }
    hashedPassword := hashPassword(password, salt)

    c := session.DB("gochat").C("accounts")
    groups := []group{}
    err = c.Insert(&clientAccount{username, hashedPassword, salt, groups})
    if err != nil {
        log.Println(err)
    }
    return nil
}

func hashPassword(plaintext string, salt []byte) []byte {
    hash, err := scrypt.Key([]byte(plaintext), salt, 1<<14, 8, 1, 64)
    if err != nil {
        log.Fatal(err)
    }

    return hash
}

func getAccount(username string) (account clientAccount, err error) {
    session, err := mgo.Dial("localhost")
    if err != nil {
        log.Println(err)
    }
    defer session.Close()

    c := session.DB("gochat").C("accounts")

    err = c.Find(bson.M{"username":username}).One(&account)
    return 
}

func validateAccount(username string, password string) bool {
    account, err := getAccount(username)
    if err != nil {
        return false
    }
    if string(account.Password) != string(hashPassword(password, account.HashSalt)) {
        return false
    }
    return true
}