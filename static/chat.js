var socket = new WebSocket("ws://" + location.hostname+(location.port ? ':'+location.port: '') + location.pathname + "/websocket");

socket.onopen = function() {console.log("connected")}
socket.onclose = function(e) {console.log("connection closed: " + e.code)}

socket.onmessage = function(e) {
    var messages = document.getElementById("messages"),
        messageObj = JSON.parse(e.data),
        messageDiv = document.createElement("div"),
        messageContent = document.createElement("div"),
        messageSender = document.createElement("div"),
        messageTime = document.createElement("div");

    messageDiv.className = "messageDiv";

    messageContent.className = "messageContent";
    messageContent.innerHTML = messageObj.Message;

    messageSender.className = "messageSender";
    messageSender.innerHTML = messageObj.Sender;

    messageTime.className = "messageTime";
    messageTime.innerHTML = messageObj.TimeStamp;

    messageDiv.appendChild(messageSender);
    messageDiv.appendChild(messageContent);
    messageDiv.appendChild(messageTime);
    messages.appendChild(messageDiv);
    
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