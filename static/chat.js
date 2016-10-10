function setName(socket) {
    var username = prompt("What is your username?");
    var password = prompt("What is your password?");
    var loginMessage = {
        msgType: "LOGIN",
        username: username,
        password: password
    }
    socket.send(JSON.stringify(loginMessage));
}



var socket = new WebSocket("ws://localhost:3000/websocket");


socket.onopen = function() {
    console.log("connected");
    setName(socket);
}
socket.onclose = function(e) {console.log("connection closed: " + e.code)}

socket.onmessage = function(e) {
    var messageBox = document.getElementById("messageBox");
    messageBox.innerHTML += e.data + "<br>";
}

function sendMessage() {
    var messageBox = document.getElementById("messageInput");
    var message = {
        Message: messageBox.value,
    }
    messageBox.value = "";
    socket.send(JSON.stringify(message));
    return false;
}

document.getElementById("chatForm").addEventListener("submit", function(event) {
    sendMessage();
    event.preventDefault();
}, false)