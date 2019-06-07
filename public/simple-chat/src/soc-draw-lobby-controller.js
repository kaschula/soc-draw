((container) => {
    function SocDrawLobbyController(uiService, sockDrawClient, socketFactory, roomApplication) {
        // events
        function enterButtonHandler2(event) {
            event.preventDefault()
        
            username = uiService.getUsernameValue()
            sockDrawClient.setGlobalUsername(username)
        
            // app.setSocket(createConnection(app.getSocketUrl()))    
            sockDrawClient.setSocket(socketFactory.createConnection(sockDrawClient.getSocketUrl()))    
        }
    
        function roomSelectHandler2(event) {
            event.preventDefault()
            
            const socket = sockDrawClient.getSocket()
        
            if (!socket) {
                console.log('Cant send message, no socket created')
                return;
            }
        
            const roomId = $(this).attr('data-room-id')
        
            if (sockDrawClient.isCurrentRoom(roomId)) {
                console.log("RoomSelectEvent:: this is the current room")
                return 
            }
        
            if (sockDrawClient.hasRoomState(roomId)) {
                // change currentroom
                changeToRoom(roomId)
                return
            }
        
            if (sockDrawClient.hasUserJoinedRoom(roomId)) {
                console.log("RoomSelectEvent:: User has already joined room")
                // request state ??
                return 
            }
        
            try {
                sockDrawClient.send(
                    JSON.stringify({
                        messageType: "LOBBY_ROOM_REQUEST", 
                        payload: JSON.stringify({roomId})
                    })
                )
            } catch(e) {
                console.log('send message error: ', e.message)
            }
        }
    
        function changeToRoom(roomId) {
            sockDrawClient.setCurrentRoomId(roomId)
            roomApplication.update(sockDrawClient.getRoomState(roomId))
            uiService.currentRoomUpdated(roomId)
        }
    
    
        this.run = function() {
            // Add Components
            uiService.hideMessageBoard()
            uiService.hideRooms()
            uiService.hideCurrentRoom()
        
            // Page events
            uiService.addEnterFormClick(enterButtonHandler2)
            uiService.setRoomSelectHandler(roomSelectHandler2)
        }
    }

    container.SocDrawLobbyController = SocDrawLobbyController
})(modules)