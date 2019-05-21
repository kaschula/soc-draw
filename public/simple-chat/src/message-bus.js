
let MessageBus

(()=> {
    function BuildMessageBus(app, uiService) {
        // ^ change app to LobbyClient
        function userSocketCreated() {
            // Should probably do a null check on the app
            uiService.replaceUserInputWithLoading()
            uiService.hideMessageBoard()
            app.getSocket().send(joinLobbyRequest())
        }
        
        function receivedUserLobbyData(payload) {
            lobbyData = JSON.parse(payload)
        
            app.setUser(lobbyData.User)
            uiService.displayRooms(lobbyData.Rooms)
            uiService.showRooms()
        
            uiService.displayUser(app.getUser()) 
        }
        
        function receivedUserJoinedRoom(payload) {
            console.log("receivedUserJoinedRoom", payload)
        
            const roomId = JSON.parse(payload).RoomId
            console.log("receivedUserLobbyData():: roomId", roomId)
            if (!roomId) {
                return // handle error
            }
        
            app.setUserJoinedRoom(roomId)
            app.setCurrentRoomId(roomId)
        
            uiService.currentRoomUpdated(roomId)
            uiService.showCurrentRoom()
        
            
            if (app.roomIsInitialisedAndStateExits(roomId)) {
                uiService.initialiseApp(roomStateContainer[roomId])
                return
            }
        
            console.log('receivedUserJoinedRoom(). room not set up, requesting room state')
            uiService.currentRoomWaitingForStatus()
            app.send(latestRoomStateRequest(roomId))
        }
        
        function appRoomBroadcast(roomId, payload) {
            console.log('appRoomBroadcast():: roomId, payload', roomId, payload)
            if (!app.isRoomInitialised(roomId)) {
                console.log('appRoomBroadcast() ERROR: room not initialised')
                return // make socket request for inital event
            }
        
        
            console.log("appRoomBroadcast::roomId", roomId)
            if (!roomId) {
                return console.log("Room ID not given, Invalid Room Broadcast")
            }
        
            const broadcast = JSON.parse(payload)
        
            if (broadcast.running === "false") {
                console.log("Room is not running")
                uiService.showRoomMessage(broadcast.message)
            }
        
            const {state} = broadcast
            console.log("appRoomBroadcast: state", state)
        
            if (app.isCurrentRoom(roomId)) {
                // RoomApplication::SimpleChatApplication
                appUpdate(state)
            }
        
            app.setRoomState(roomId, state)
        }
        
        function appRoomBroadcastRoomStarted(roomId, payload) {
            console.log("appRoomBroadcastRoomStarted() ")
            const broadcast = JSON.parse(payload)
        
            console.log("setting up App Inital State in Room", roomId, broadcast)
        
            if (!broadcast.state) {
                console.log("Broadcast did not return state")
            }
        
            uiService.initialiseAppWindow()
            // This is RoomApplication
            appInitalise(broadcast.state)
            app.setRoomState(roomId, broadcast.state)
            app.setRoomInitialised(roomId)
        }
    
        function handleError(payload) {
            uiService.displayErrorMessage(payload)
            if (!socket) {
                socket = null
                username = null
                return 
            }
        
            app.getSocket().close()
            uiService.showUserForm()
        }
    
        // chnage msg -> message
        function dispatch(msg) {
            switch(msg.Type) {
                case "CREATED": 
                    userSocketCreated()
                    return 
                case "USER_LOBBY_DATA": 
                    receivedUserLobbyData(msg.Payload)
                    return 
                case "USER_JOINED_ROOM":
                    receivedUserJoinedRoom(msg.Payload)
                    return
                case "ROOM_BROADCAST":
                    appRoomBroadcast(msg.RoomID, msg.Payload)
                    return
                case "ROOM_BROADCAST_INIT":
                    appRoomBroadcastRoomStarted(msg.RoomID, msg.Payload)
                    return
                case "ERROR": 
                    handleError(msg.Payload)
                    console.log("The error: ", msg)
                    return
                default:
                    console.log("App message type not recongised: ", msg)
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
    
    MessageBus = BuildMessageBus
})()