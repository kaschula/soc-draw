((container) => {
    SimpleChatApp.prototype = container.roomApplicationInterface
    function SimpleChatApp(uiService, lobby) {
        function sendMessageHandler(event) {
            event.preventDefault()

            const message = getMessage()
            
            if (!message || message == '') {
                return
            }
        
            try {
                let payload = createMessageJson(
                    lobby.getGlobalUsername(), 
                    lobby.getCurrentRoomId(), 
                    message
                );
                console.log('message json: ', payload)
                lobby.send(payload)
            } catch(e) {
                console.log('send message error: ', e.message)
            }
        }
        
        function getMessage() {
            // All UI service stuff must just use Query
            const message =  uiService.getMessageValue()
            uiService.clearMessageValue()
        
            return message
        }
        
        function hydrateApplicationWindow(state) {
            console.log("hydrateApplicationWindow()", state)
            const messages = uiService.query("#current-room #application-window .messages")
            console.log("messages div", messages)
        
            if (!messages.length) {
                console.log("Error: App messages div could not be found")
                return 
            }

            if (!state.messages) {
                console.log("Error: messages could not be found")
                return  
            }

            messages.empty()

            state.messages.forEach((message) => {
                messages.append(messageComponent(username, message))
            })


        }
        
        function messagingAppComponent() {
            return (`
                <div id="simple-chat-message" class="row">
                    <div id="message-board" class="col-md-12">
                        <div id="messages-row" class="row">
                            <div class="messages" class="col-md-12">
        
                            </div>
                        </div>
                        <div id="message-form" class="row">
                            <form class="col-md-12">
                                <input name="message" type="textbox">
                                <button id="chat-submit">Send</button>
                            </form>
                        </div>
                    </div>
                </div>
            `)
        }
        
        function messageComponent(username, message) {
            let className = null
            
            if (message.sender === "room") {
                className = "message-room";
            } else if (message.sender === username) {
                className = "message-sent";
            } else {
                className = "message-received"; 
            }
        
            return (`
                <div class="${className}"><span class="user">${message.sender}: </span><span class="message-value">${message.message}</span></div>
            `)
        }
        
        
        // -------------- App.Message stuff
        function createMessageJson(username, roomId, message) {
            const eventMessage = {
                messageType: "ROOM_REQUEST",
                payLoad: JSON.stringify({username, roomId, message, requestType: "ROOM_EVENT"})
            }
        
            // Send new message type
            // Need to change recieve Message as well so it canbe added to board
            return JSON.stringify(eventMessage);
        }

        this.initialise = function(state) {
            console.log('roomApplication.initialise() called')
            const appWindow = uiService.query("#application-window")
            const x = messagingAppComponent();
            appWindow.append(x)
        
            // uiService can only Query
            uiService.addChatMessageSubmit(sendMessageHandler)
        
            hydrateApplicationWindow(state)
        }
        
        this.update = function(state) {
            hydrateApplicationWindow(state)
        }
    }

    container.buildRoomApplication = (uiService, lobby) => new SimpleChatApp(uiService, lobby)
})(modules)