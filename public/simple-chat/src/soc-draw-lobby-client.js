let SocDrawLobbyClient;

(() => {
    function SocDrawLobbyClientBuild() {
        // private
        // const uiService = newUIService($);
        const roomStateContainer = {};
        const roomsInitialised = {};
        const roomsJoined = {}
        const url = 'localhost:8089/ws' // Dep
    
        let globalUsername; // is this being used
        let globalUser;
        let globalCurrentRoomId;
        let socket;
    
        // public
    
        this.getSocketUrl = function() {
            return url;
        }
        
        this.setCurrentRoomId = function(id) {
            globalCurrentRoomId = id
        }
    
        this.getCurrentRoomId = function() {
            return globalCurrentRoomId
        }
        
        this.setUser = function(user) {
            this.setGlobalUsername(user.ID)
            globalUser = user
        }
        
        this.getUser = function() {
            return globalUser
        }
        
        this.getGlobalUsername = function() {
            return globalUsername;
        }
        
        this.setGlobalUsername = function(name) {
            globalUsername = name;
        }
        
        this.getRoomState = function(roomId) {
            return !!roomStateContainer[roomId] ? roomStateContainer[roomId] : {}
        }
        
        this.setUserJoinedRoom = function(roomId) {
            roomsJoined[roomId] = true
        }
        
        this.hasUserJoinedRoom = function(roomId) {
            return !!roomsJoined[roomId]
        }
        
        this.setRoomState = function (roomId, state) {
            roomStateContainer[roomId] = state
        }
        
        this.hasRoomState = function(roomId) {
            return !!roomStateContainer[roomId]
        }
        
        this.roomIsInitialisedAndStateExits = function(roomId) {
            return roomsInitialised[roomId] && !!roomStateContainer[roomId];
        }
        
        this.isCurrentRoom = function(roomId) {
            return roomId == this.getCurrentRoomId()
        }
    
        this.getSocket = function() {
            return socket
        }
    
        this.setSocket = function(connection) {
            socket = connection
        }
    
        this.send = function(payload) {
            socket.send(payload)
        }
    
        this.isRoomInitialised = function (roomId) {
            return !!roomsInitialised[roomId]
        }
    
        this.setRoomInitialised = function (roomId) {
            roomsInitialised[roomId] = true;
        } 

        this.processSocketData = function(data) {
            if (!data.Type || !data.Payload) {
                console.log("Message incomplete: ")
                return
            }

            messageBus.dispatch(data)
        }
    }

    SocDrawLobbyClient = SocDrawLobbyClientBuild
})();