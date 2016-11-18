var cookieList = document.cookie.split("; ");

if (cookieList.includes("login=Failed")) {
    alert("incorrect username/password")
    document.cookie = "login=; expires=Thu, 01 Jan 1970 00:00:00 UTC";
}

if (cookieList.includes("signup=Failed")) {
    alert("username taken")
    document.cookie = "signup=; expires=Thu, 01 Jan 1970 00:00:00 UTC";
}

