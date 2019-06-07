console.log("main.js Loaded", modules)

// ----- App SetUp
if (!$) {
    throw new Error('jQuery not found');
}

function createApp() {
    const components = new modules.ComponentLibrary()   
    const uiService = new modules.UIService($, components);
    const lobby = new modules.SocDrawLobbyClient();
    
    const messageFactory = new modules.MessageFactory()
    const roomApp = modules.buildRoomApplication(uiService, lobby)
    const messageBus = new modules.MessageBus(lobby, uiService, messageFactory, roomApp);
    const socketFactory = new modules.SocketFactory(messageBus, lobby, messageFactory);

    return new modules.SocDrawLobbyController(uiService, lobby, socketFactory, roomApp)
}
const app = createApp()
app.run()
