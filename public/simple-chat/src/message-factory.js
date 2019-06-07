((container) => {
    function MessageFactory() {
        this.unmarshallAppMessageJson = function (messagePayload) {
            const m = JSON.parse(messagePayload)
        
            return m;
        }
        this.unmarshallRoomMessageJson = function (messagePayload) {
            // console.log("unmarshallRoomMessageJson: ")
            const m = JSON.parse(messagePayload)
        
            if ( !m.username || !m.room || !m.message) {
                throw new Error("room message cant be resolved from payload: ", messagePayload)
            }
        
            return m;
        }
        
        this.joinLobbyRequest = function (username) {
            return JSON.stringify({
                messageType:"LOBBY_USER_JOIN_REQUEST", 
                payload: JSON.stringify({user: username}) 
            })
        }
        
        this.latestRoomStateRequest = function (roomId) {
            return JSON.stringify({
                messageType:"ROOM_REQUEST", 
                payload: JSON.stringify({roomId, requestType: "STATE"}) 
            })
        }
    }

    container.MessageFactory = MessageFactory
})(modules)