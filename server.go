package main

import (
    "log"
    "os"
    "net/http"
    "github.com/gorilla/mux"
    "html/template"
    "strconv"
)

func main() {
    r := mux.NewRouter()

    r.HandleFunc("/login", loginPageHandler).Methods("POST")
    r.HandleFunc("/logout", logoutPageHandler).Methods("POST")
    r.HandleFunc("/signup", signupPageHandler).Methods("POST")
    r.HandleFunc("/createGroupPage", createGroupPageHandler)
    r.HandleFunc("/createGroup", groupCreationHandler).Methods("POST")

    r.HandleFunc("/websocket/{group}", messagingHandler)

    r.HandleFunc("/chat/{group}", chatPageHandler)

    r.HandleFunc("/", landingPageHandler)
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

    log.SetOutput(os.Stdout)
    log.Println("Listening on port 3000...")
    
    http.Handle("/", r)

    http.ListenAndServe(":3000", nil)
}

func createGroupPageHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/groupmaking.html")
}

func groupCreationHandler(w http.ResponseWriter, r *http.Request) {
    groupName := r.FormValue("groupName")
    creator := getUsername(r)

    groupMemberNames := make([]string, 0)
    groupMemberNames = append(groupMemberNames, creator);
    for i := 1; r.FormValue("groupMember" + strconv.Itoa(i)) != ""; i++ {
        groupMemberNames = append(groupMemberNames, r.FormValue("groupMember" + strconv.Itoa(i)))
    }

    groupMembers := make([]clientAccount, 0)
    for _, member := range groupMemberNames {
        groupMemberAccount, err := getAccount(member)
        if err != nil {
            log.Println(err)
            return
        }
        groupMembers = append(groupMembers, groupMemberAccount)
    }

    err := createGroup(groupName, groupMembers)
    if err != nil {
        log.Println(err)
        return
    }
    http.Redirect(w, r, "/chat/" + groupName, 302)
}

func landingPageHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/first.html")
}

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
    username, password := r.FormValue("username"), r.FormValue("password")
    redirectTarget := "/"
    if username != "" && password != "" {
        if validateAccount(username, password) {
            account, _ := getAccount(username)
            if len(account.Groups) > 0 {
                redirectTarget = "/chat/" + account.Groups[0].Name
            } else {
                redirectTarget = "/createGroupPage"
            }
            setSession(username, w)
        }
        http.Redirect(w, r, redirectTarget, 302)
    }
}

func logoutPageHandler(w http.ResponseWriter, r *http.Request) {
    clearSession(w)
    http.Redirect(w, r, "/", 302)
}

func signupPageHandler(w http.ResponseWriter, r *http.Request) {
    username, password := r.FormValue("username"), r.FormValue("password")
    err := createAccount(username, password)
    if err == nil {
        setSession(username, w)
        http.Redirect(w, r, "/createGroupPage", 302)
    } else {
        http.Redirect(w, r, "/", 302)
    }
}

func chatPageHandler(w http.ResponseWriter, r *http.Request) {
    t := template.New("chat.html")
    t, err := t.ParseFiles("static/chat.html")
    if err != nil {
        log.Println(err)
    }
    user := getUsername(r)
    account, _ := getAccount(user)

    templateInfo := struct {
        Groups []group
        ActiveGroup string
    } {
        account.Groups,
        mux.Vars(r)["group"],
    }
    if err = t.Execute(w, templateInfo); err != nil {
        log.Println(err)
    }
}