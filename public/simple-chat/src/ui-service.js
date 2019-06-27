((container)=>{
    function UIService(jQueryLib, components) {
        const _$ = (query) => jQueryLib(query)
        let roomSelectHandler = () => {}
        
        // events
        const addEnterFormClick = (handler) => _$('#enter-form-submit').click(handler)
        const addRoomEventHandler = (handler) => _$('.room-select-btn').click(handler)
        
        // inputs (These should return just the values)
        const getUsernameValue = () => _$('input[name=username]').val()

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
            const applicationWindow = _$("#application-window")
            applicationWindow.empty()

            showRoomMessage("Room waiting to start", true)
        } 
        const displayUser = (user) => {
            const container = _$("#entry-form-container")
            container.empty()
            container.append(components.userComponent(user))
        }
    
        const hideMessageBoard = () =>  _$("#message-board").hide()
    
        const initialiseAppWindow = () => {
            const applicationWindow = _$("#application-window")
            if (applicationWindow.length) {
                return
            }

            const currentRoom = _$("#current-room")
            currentRoom.append(components.applicationWindowComponent())
        }
    
        // UI Events
        const currentRoomUpdated = (roomName) => {
            const roomSpan = _$("#current-room-name")
            roomSpan.html(`Your in room <i>${roomName}</i>`)
        }
    
        const setRoomSelectHandler = (handler) => roomSelectHandler = handler

        destroyApp = () => {
            const appWindow = _$("#application-window")
            appWindow.empty()
        }
    
        return {
            query: _$,
            addEnterFormClick,
            getUsernameValue,
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
            destroyApp,
        }
    }

    container.UIService = UIService
})(modules)