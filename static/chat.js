var socket = new WebSocket("ws://localhost:3000/websocket/" + window.location.href.split("/").pop());

socket.onopen = function() {console.log("connected")}
socket.onclose = function(e) {console.log("connection closed: " + e.code)}

socket.onmessage = function(e) {
    var messages = document.getElementById("messages");
    var messageObj = JSON.parse(e.data);
    messages.innerHTML += messageObj.Sender + ": " + messageObj.Message + "<br>";
    messages.scrollTop = messageBox.scrollHeight;
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