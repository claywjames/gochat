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
    r.HandleFunc("/joinGroupPage", joinGroupPageHandler)
    r.HandleFunc("/joinGroup", groupJoinHandler).Methods("POST")

    r.HandleFunc("/chat/{group}/websocket", messagingHandler)

    r.HandleFunc("/chat/{group}", chatPageHandler)

    r.HandleFunc("/", landingPageHandler)
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

    log.SetOutput(os.Stdout)

    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }
    log.Println("Listening on port " + port)
    
    http.Handle("/", r)

    http.ListenAndServe(":" + port, nil)
}

func joinGroupPageHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/joingroup.html")
}

func groupJoinHandler(w http.ResponseWriter, r *http.Request) {
    groupCode := r.FormValue("groupCode")
    joiner := getUsername(r)
    group, err := getGroupFromCode(groupCode)
    if err != nil {
        log.Println(err)
        return
    }
    joinerAccount, _ := getAccount(joiner)
    err = group.addGroupMember(joinerAccount)
    if err != nil {
        log.Println(err)
        return
    }

    http.Redirect(w, r, "/chat/" + group.Name, 302)
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
            setFailedGroupCreationCookie(w, "member does not exist")
            http.Redirect(w, r, "/createGroupPage", 302)
            return
        }
        groupMembers = append(groupMembers, groupMemberAccount)
    }

    err := createGroup(groupName, groupMembers)
    if err != nil {
        setFailedGroupCreationCookie(w, "group name taken")
        http.Redirect(w, r, "/createGroupPage", 302)
        return
    }
    http.Redirect(w, r, "/chat/" + groupName, 302)
}

func landingPageHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "static/login.html")
}

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
    username, password := r.FormValue("username"), r.FormValue("password")
    redirectTarget := "/"
    if validateAccount(username, password) {
        account, _ := getAccount(username)
        if len(account.Groups) > 0 {
            redirectTarget = "/chat/" + account.Groups[0].Name
        } else {
            redirectTarget = "/createGroupPage"
        }
        setSession(username, w)
    } else {
        setBadLoginCookie(w)
    }
    http.Redirect(w, r, redirectTarget, 302)
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
        setFailedSignUpCookie(w)
        http.Redirect(w, r, "/", 302)
    }
}

func chatPageHandler(w http.ResponseWriter, r *http.Request) {
    validRequest := false
    account, _ := getAccount(getUsername(r))
    requestedGroup, _ := getGroup(mux.Vars(r)["group"])
    for _, group := range account.Groups {
        if group.Name == requestedGroup.Name {
            validRequest = true
        }
    }

    if validRequest {
        t := template.New("chat.html")
        t, err := t.ParseFiles("static/chat.html")
        if err != nil {
            log.Println(err)
        }

        templateInfo := struct {
            Groups []group
            ActiveGroup string
            ActiveGroupCode string
        } {
            account.Groups,
            requestedGroup.Name,
            requestedGroup.Code,
        }
        if err = t.Execute(w, templateInfo); err != nil {
            log.Println(err)
        }
    } else {
        http.Error(w, "This group either doesn't exist or you aren't a member.", 403)
    }
}