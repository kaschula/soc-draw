((container) => {
    function ComponentLibrary() {
        function userFormComponent() {
            return `
                <form id="entry-form">
                    <label>Username</label>
                    <input name="username" type="text">
                    <button id="enter-form-submit">Enter</button> 
                </form>
            `
        }
        
        function roomSelectComponent(room) {
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
        
        // can delete this
        function roomWaitingComponent() {
            return (`
                <p id="room-waiting-message" >Request latest room state.... </p>
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
        this.roomSelectComponent = roomSelectComponent
        this.userComponent = userComponent
        this.messageComponent = messageComponent
        this.roomWaitingComponent = roomWaitingComponent
        this.roomStatusMessageComponent = roomStatusMessageComponent
        this.applicationWindowComponent = applicationWindowComponent
    }

    container.ComponentLibrary = ComponentLibrary
})(modules)