((container) => {
    function SocDrawLobbyController(uiService, sockDrawClient, socketFactory, roomApplication) {
        function enterButtonHandler(event) {
            event.preventDefault()
        
            username = uiService.getUsernameValue()
            sockDrawClient.setGlobalUsername(username)
        
            sockDrawClient.setSocket(socketFactory.createConnection(sockDrawClient.getSocketUrl()))    
        }
    
        function roomSelectHandler(event) {
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
                changeToRoom(roomId)
                return
            }
        
            if (sockDrawClient.hasUserJoinedRoom(roomId)) {
                console.log("RoomSelectEvent:: User has already joined room")
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
            const state = sockDrawClient.getRoomState(roomId);

            uiService.destroyApp()
            roomApplication.initialise(state)
            roomApplication.update(state)

            uiService.currentRoomUpdated(roomId)
        }
    
    
        this.run = function() {
            // Add Components
            uiService.hideMessageBoard()
            uiService.hideRooms()
            uiService.hideCurrentRoom()
        
            // Page events
            uiService.addEnterFormClick(enterButtonHandler)
            uiService.setRoomSelectHandler(roomSelectHandler)
        }
    }

    container.SocDrawLobbyController = SocDrawLobbyController
})(modules)