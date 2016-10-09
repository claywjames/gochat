var socket = new WebSocket("ws://localhost:3000/websocket");

socket.onopen = function() {console.log("connected")}
socket.onclose = function(e) {console.log("connection closed: " + e.code)}

socket.onmessage = function(e) {
    var messageBox = document.getElementById("messageBox");
    messageBox.innerHTML += e.data + "<br>";
}

function sendMessage() {
    var messageBox = document.getElementById("messageInput");
    var message = messageBox.value;
    messageBox.value = "";
    socket.send(message);
    return false;
}

document.getElementById("chatForm").addEventListener("submit", function(event) {
    sendMessage();
    event.preventDefault();
}, false)