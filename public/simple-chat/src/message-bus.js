((container)=> {
    function MessageBus(lobby, uiService, messageFactory, roomApplication) {
        function userSocketCreated() {
            uiService.replaceUserInputWithLoading()
            uiService.hideMessageBoard()
            lobby.getSocket().send(messageFactory.joinLobbyRequest(lobby.getGlobalUsername()))
        }
        
        function receivedUserLobbyData(payload) {
            lobbyData = JSON.parse(payload)
        
            lobby.setUser(lobbyData.User)
            uiService.displayRooms(lobbyData.Rooms)
            uiService.showRooms()
        
            uiService.displayUser(lobby.getUser()) 
        }
        
        function receivedUserJoinedRoom(payload) {
            const roomId = JSON.parse(payload).RoomId
            if (!roomId) {
                return // handle error
            }
        
            lobby.setUserJoinedRoom(roomId)
            lobby.setCurrentRoomId(roomId)
        
            uiService.currentRoomUpdated(roomId)
            uiService.showCurrentRoom()
        
            if (lobby.roomIsInitialisedAndStateExists(roomId)) {
                roomApplication.update(lobby.getRoomState(roomId))
                return
            }
        
            uiService.currentRoomWaitingForStatus()
            lobby.send(messageFactory.latestRoomStateRequest(roomId))
        }
        
        function appRoomMessage(roomId, payload) {
            const broadcast = JSON.parse(payload)

            if(!lobby.isCurrentRoom(roomId)) {
                return
            }

            uiService.showRoomMessage(broadcast.message)
        }

        function appRoomBroadcast(roomId, payload) {
            if (!lobby.isRoomInitialised(roomId)) {
                console.log('appRoomBroadcast() ERROR: room not initialised')
                return // make socket request for inital event
            }
        
            if (!roomId) {
                return console.log("Room ID not given, Invalid Room Broadcast")
            }
        
            const broadcast = JSON.parse(payload)
        
            if (broadcast.running === "false") {
                console.log("Room is not running")
                uiService.showRoomMessage(broadcast.message, true)
            }
        
            const {state} = broadcast
        
            if (lobby.isCurrentRoom(roomId)) {
                roomApplication.update(state)
            }

            lobby.setRoomState(roomId, state)
        }
        
        function appRoomBroadcastInit(roomId, payload) {
            const broadcast = JSON.parse(payload)
        
            if (!broadcast.state) {
                console.log("Broadcast did not return state")
                return 
            }

            if (lobby.roomIsInitialisedAndStateExists(roomId)) {
                return
            }
        
            uiService.initialiseAppWindow()
            roomApplication.initialise(broadcast.state)
            lobby.setRoomState(roomId, broadcast.state)
            lobby.setRoomInitialised(roomId)
            uiService.removeRoomMessages()
        }
    
        function handleError(payload) {
            uiService.displayErrorMessage(payload)
            if (!socket) {
                socket = null
                username = null
                return 
            }
        
            lobby.getSocket().close()
            uiService.showUserForm()
        }
    
        // change msg -> message
        function dispatch(message) {
            switch(message.Type) {
                case "CREATED": 
                    userSocketCreated()
                    return 
                case "USER_LOBBY_DATA": 
                    receivedUserLobbyData(message.Payload)
                    return 
                case "USER_JOINED_ROOM":
                    receivedUserJoinedRoom(message.Payload)
                    return
                case "ROOM_BROADCAST":
                    appRoomBroadcast(message.RoomID, message.Payload)
                    return
                case "ROOM_BROADCAST_MESSAGE":
                    appRoomMessage(message.RoomID, message.Payload)
                    return
                case "ROOM_BROADCAST_INIT":
                    appRoomBroadcastInit(message.RoomID, message.Payload)
                    return
                case "ERROR": 
                    handleError(message.Payload)
                    console.log("The error: ", message)
                    return
                default:
                    console.log("App message type not recongised: ", message)
                    return
            }
        }
    
        this.processSocketData = function(data) {
            if (!data.Type || !data.Payload) {
                console.log("Message incomplete: ")
                return
            }
    
            dispatch(data)
        }
    }
    
    container.MessageBus = MessageBus
})(modules)