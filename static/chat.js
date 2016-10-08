var socket = new WebSocket("ws://localhost:3000/websocket");

socket.onopen = function() {console.log("connected")}
socket.onclose = function(e) {console.log("connection closed: " + e.code)}

socket.onmessage = function(e) {
    var messageBox = document.getElementById("messageBox");
    messageBox.innerHTML += e.data + "<br>";
}

function sendMessage() {
    var message = document.getElementById("messageInput").value;
    socket.send(message);
    return false;
}
