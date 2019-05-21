let ComponentLibrary

(() => {
    function ComponentLibraryBuilder() {
        function userFormComponent() {
            return `
                <form id="entry-form">
                    <label>Username</label>
                    <input name="username" type="text">
                    <button id="enter-form-submit">Enter</button> 
                </form>
            `
        }
        
        function roomComponent(room) {
            return `
                <div class="room">
                    <span>${room.Name}</span>
                    <button data-room-id="${room.ID}" class="room-select-btn">Enter</button>
                </div>
            `
        }
        
        function userComponent(user) {
            return `<div> Welcome ${user.ID}</div>`
        }
        
        function messageComponent(username, messageType, message) {
            return (
                `<div class="${messageType}">
                    <span class="user"><i>${username}</i>: </span>
                    <span class="message-value">${message}</span>
                </div>`
            )
        }
        
        function roomWaitingComponent() {
            return (`
                <p>Request latest room state.... </p>
            `)
        }
        
        
        function roomStatusMessageComponent(message) {
            return (
                `<p>Room Status: ${message}</p>`
            )
        }
        
        function applicationWindowComponent() {
            return (
                `<div id="application-window"></div>`
            )
        }
    
        this.userFormComponent = userFormComponent
        this.roomComponent = roomComponent
        this.userComponent = userComponent
        this.messageComponent = messageComponent
        this.roomWaitingComponent = roomWaitingComponent
        this.roomStatusMessageComponent = roomStatusMessageComponent
        this.applicationWindowComponent = applicationWindowComponent
    }

    ComponentLibrary = ComponentLibraryBuilder
})()