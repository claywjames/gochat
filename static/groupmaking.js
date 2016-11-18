var cookieList = document.cookie.split("; ");

if (cookieList.includes("creation=group name taken")) {
    alert("The group name you have chosen is already in use.")
    document.cookie = "creation=; expires=Thu, 01 Jan 1970 00:00:00 UTC";
}

if (cookieList.includes("creation=member does not exist")) {
    alert("A group member does not exist.")
    document.cookie = "creation=; expires=Thu, 01 Jan 1970 00:00:00 UTC";
}

function addMember() {
    var groupMembers = document.getElementById("groupMembers");
    var numMembers = groupMembers.childElementCount;
    groupMembers.insertAdjacentHTML("beforeend", "<input name='groupMember" + numMembers + "' placeholder='Enter group memeber' required>")
}