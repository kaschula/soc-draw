((container) => {
    function ComponentLibrary() {
        function userFormComponent() {
            return `
            <div>
                <h2 class="text-center">Please Login</h2>
                <form id="entry-form">
                    <div class="form-group">
                        <label class="">Username: </label>
                        <input class="form-control" name="username" type="text">
                    </div>
                    <div class="form-group">
                        <button id="enter-form-submit" class="btn btn-primary col-sm-12" type="submit">Enter</button> 
                    </div>
                </form>
            </div>      
            `
        }
        
        function roomSelectComponent(room) {
            return `
                <div class="room list-group-item w-100">
                    <p class="text-left sd-display-inline" >${room.Name}</p>
                    <button class="room-select-btn btn btn-secondary btn-sm text-right sd-display-inline" data-room-id="${room.ID}" class="room-select-btn">Enter</button>
                </div>
            `
        }
        
        function userComponent(user) {
            return `
            <div class="row">  
                <div class="card w-100">
                    <div class="card-body"> Welcome ${user.ID} </div>
                </div>
            </div>
            `
        }
        
        function messageComponent(username, messageType, message) {
            return (
                `<div class="${messageType}">
                    <span class="user"><i>${username}</i>: </span>
                    <span class="message-value">${message}</span>
                </div>`
            )
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
        this.roomStatusMessageComponent = roomStatusMessageComponent
        this.applicationWindowComponent = applicationWindowComponent
    }

    container.ComponentLibrary = ComponentLibrary
})(modules)