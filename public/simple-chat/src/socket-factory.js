((container) => {
    function SocketFactory(messageBus, lobbyClient, messageFactory) {
        this.createConnection = function (url) {
            const ws = new WebSocket(`ws://${url}`)
            attachSocketHandlers(ws)
        
            return ws
        }
    
        function attachSocketHandlers(ws) {
            ws.onopen = onOpen
            ws.onclose = onClose
            ws.onmessage = onMessage
            ws.onerror = onError
        }
        
        function onOpen(data) {
            console.log('onOpen', data)
        }
        
        function onClose(data) {
            console.log('onClose', data)
            lobbyClient.setSocket(undefined)
        }
        
        function onMessage(messageEvent) {
            console.log("socket on message")
            const msg = messageFactory.unmarshallAppMessageJson(messageEvent.data);
            
            try {
                messageBus.processSocketData(msg)
            } catch (e) {
                console.log("Error processing App message: ", e.message, e.stack)
            }
        }
        
        function onError(data) {
            console.log('onError', data)
        }
    }

    container.SocketFactory = SocketFactory 
})(modules)