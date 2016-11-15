function addMember() {
    console.log("called");
    var groupMembers = document.getElementById("groupMembers");
    var numMembers = groupMembers.childElementCount;
    groupMembers.insertAdjacentHTML("beforeend", "<input name='groupMember" + numMembers + "' placeholder='Enter group memeber' required>")
}