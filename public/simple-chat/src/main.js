console.log("main.js Loaded", SocDrawLobbyClient)

// ----- App SetUp
if (!$) {
    throw new Error('jQuery not found');
}

// App Page Data
const components = new ComponentLibrary()
const uiService = new UIService($, components);
const app = new SocDrawLobbyClient();
const messageBus = new MessageBus(app, uiService);
const socketFactory = new SocketFactory(messageBus);

init()

function init() {
    // initalise UI Page

    // Add Components
    uiService.hideMessageBoard()
    uiService.hideRooms()
    uiService.hideCurrentRoom()

    // Page events
    uiService.addEnterFormClick(enterButtonHandler)
}

// ----- SocDraw Code

// Event CTA handlers
function enterButtonHandler(event) {
    event.preventDefault()

    username = getUsernameInput()
    app.setGlobalUsername(username)

    // app.setSocket(createConnection(app.getSocketUrl()))    
    app.setSocket(socketFactory.createConnection(app.getSocketUrl()))    
}


// UI CTA Event Handler
function roomSelectHandler(event) {
    event.preventDefault()
    console.log('room event: ');
    
    const socket = app.getSocket()

    if (!socket) {
        console.log('Cant send message, no socket created')
        return;
    }

    const roomId = $(this).attr('data-room-id')

    if (app.isCurrentRoom(roomId)) {
        console.log("RoomSelectEvent:: this is the current room")
        return 
    }

    if (app.hasRoomState(roomId)) {
        // change currentroom
        changeToRoom(roomId)
        return
    }

    if (app.hasUserJoinedRoom(roomId)) {
        console.log("RoomSelectEvent:: User has already joined room")
        // request state ??
        return 
    }

    try {
        app.send(
            JSON.stringify({
                messageType: "LOBBY_ROOM_REQUEST", 
                payload: JSON.stringify({roomId})
            })
        )
    } catch(e) {
        console.log('send message error: ', e.message)
    }
}


// App - Controller Actions
function changeToRoom(roomId) {
    // Check if initialised ??
    app.setCurrentRoomId(roomId)
    appUpdate(app.getRoomState(roomId))
    uiService.currentRoomUpdated(roomId)
}

// Move into SocDrawClient
function getUsernameInput() {
    return uiService.getUsernameValue();
}


function unmarshallAppMessageJson(messagePayload) {
    console.log("unmarshallAppMessageJson: ")
    const m = JSON.parse(messagePayload)

    return m;
}
function unmarshallRoomMessageJson(messagePayload) {
    console.log("unmarshallRoomMessageJson: ")
    const m = JSON.parse(messagePayload)

    if ( !m.username || !m.room || !m.message) {
        throw new Error("room message cant be resolved from payload: ", messagePayload)
    }

    return m;
}


function joinLobbyRequest() {
    return JSON.stringify({
        messageType:"LOBBY_USER_JOIN_REQUEST", 
        payload: JSON.stringify({user: username}) 
    })
}

function latestRoomStateRequest(roomId) {
    return JSON.stringify({
        messageType:"ROOM_REQUEST", 
        payload: JSON.stringify({roomId, requestType: "STATE"}) 
    })
}
// ---------------------------------------------------------------------------


// ------------------ SimpleChatApp code -----------------
// All UI service stuff should be removed into functions that use UI service query only but are contained here as part of the SimpleChatApp

function sendMessageHandler(event) {
    event.preventDefault()

    // if (!appsocket) {
    //     console.log('Cant send message, no socket created')
    //     return;
    // }

    const message = getMessage()
    
    if (!message || message == '') {
        return
    }



    try {
        let payload = createMessageJson(
            app.getGlobalUsername(), 
            app.getCurrentRoomId(), 
            message
        );
        console.log('message json: ', payload)
        app.send(payload)
    } catch(e) {
        console.log('send message error: ', e.message)
    }
}

function getMessage() {
    const message =  uiService.getMessageValue()
    uiService.clearMessageValue()

    return message
}

// RoomApplication::SimpleChatApplication
function appInitalise(state) {
    console.log('uiService.initialiseApp() called')

    const appWindow = uiService.query("#application-window")
    appWindow.append(messagingAppComponent())

    // uiService can only Query
    uiService.addChatMessageSubmit(sendMessageHandler)

    hydrateApplicationWindow(state)
}

// RoomApplication::SimpleChatApplication
function appUpdate(state) {
    hydrateApplicationWindow(state)
}

// RoomApplication::SimpleChatApplication
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

// RoomApplication::SimpleChatApplication Components
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