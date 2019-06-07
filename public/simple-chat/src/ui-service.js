((container)=>{
    function UIService(jQueryLib, components) {
        const _$ = (query) => jQueryLib(query)
        let roomSelectHandler = () => {}
        
        // events
        // these should use handler passed in
        const addEnterFormClick = (handler) => _$('#enter-form-submit').click(handler)
        // This should be moved into RoomApplication
        const addChatMessageSubmit = (handler) => _$('#chat-submit').click(handler)
        const addRoomEventHandler = (handler) => _$('.room-select-btn').click(handler)
        
        // inputs (These should return just the values)
        const getUsernameValue = () => _$('input[name=username]').val()
        const getRoomValue = () => _$('input[name=room]').val()
        const getMessageValue = () => _$('input[name=message]').val()
        const clearMessageValue = () => _$('input[name=message]').val('')
    
        // messages
        const getMessagesContainer = () => _$('#messages')
        const clearMessages = () => getMessagesContainer().empty()
        const appendMessage = (username, messageType, message) => {
            getMessagesContainer().append(components.messageComponent(username, messageType, message))
        }
        const focusOnSendButton = () => _$('#chat-submit').focus()
        const replaceUserInputWithLoading = () => {
            const container = _$("#entry-form-container")
            container.empty()
            container.append(`<p>Welcome. Fetching details..... </p>`)
    
        }
    
        const showUserForm = () => {
            const container = _$("#entry-form-container")
            container.empty()
            container.append(components.userFormComponent())
            uiService.addEnterFormClick(enterButtonHandler)
        }
    
        const displayErrorMessage = message => {
            const container = _$("#error-container")
            container.empty()
            container.append(`<p>Error: ${message}</p>`)
            setTimeout(() => container.empty(), 3000)
        }
    
        const displayRooms = (rooms) => {
            const container = _$("#rooms-container")
            // container.empty()
            rooms.forEach(room => {
                container.append(components.roomSelectComponent(room))
            });
            addRoomEventHandler(roomSelectHandler)
        }
    
        const hideRooms = () => _$("#rooms-container").hide()
        const showRooms = () => _$("#rooms-container").show()
    
        const hideCurrentRoom = () => _$("#current-room").hide()
        const showCurrentRoom = () => _$("#current-room").show()
        const showRoomMessage = (message, persist) => { 
            const container = _$("#current-room-log")
            container.empty()
            container.append(components.roomStatusMessageComponent(message))

            if (!persist) {
                setTimeout(() => container.empty, 5000)
            }
        }
        const removeRoomMessages = () => {
            const container = _$("#current-room-log")
            container.empty()
        }
    
        const currentRoomWaitingForStatus = () => {
            // const currentRoom = _$("#current-room")
            // currentRoom.empty()
            showRoomMessage("Room waiting to start", true)
            // currentRoom.append(components.roomWaitingComponent())
        } 
        const displayUser = (user) => {
            const container = _$("#entry-form-container")
            container.empty()
            container.append(components.userComponent(user))
        }
    
        const hideMessageBoard = () =>  _$("#message-board").hide()
    
        // ChatApp code should be moved out of UI service
        const initialiseAppWindow = () => {
            const currentRoom = _$("#current-room")
            // currentRoom.empty()
            // need to remove room waiting status of present
            currentRoom.append(components.applicationWindowComponent())
        }
    
        // UI Events
        const currentRoomUpdated = (roomName) => {
            const roomSpan = _$("#current-room-name")
            roomSpan.html(`Your in room <i>${roomName}</i>`)
        }
    
        const setRoomSelectHandler = (handler) => roomSelectHandler = handler
    
        return {
            query: _$,
            addEnterFormClick,
            addChatMessageSubmit,
            getMessageValue,
            clearMessageValue,
            getUsernameValue,
            getRoomValue,
            clearMessages,
            appendMessage,
            focusOnSendButton,
            replaceUserInputWithLoading,
            displayErrorMessage,
            showUserForm,
            displayRooms,
            hideMessageBoard,
            displayUser,
            hideRooms,
            showRooms,
            hideCurrentRoom,
            showCurrentRoom,
            currentRoomWaitingForStatus,
            showRoomMessage,
            initialiseAppWindow,
            currentRoomUpdated,
            setRoomSelectHandler,
            removeRoomMessages,
            // hydrateApplicationWindow,
        }
    }

    container.UIService = UIService
})(modules)