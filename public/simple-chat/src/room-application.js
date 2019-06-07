((container) => {
    roomApplicationInterface = {
        initialise: () => {throw new Error("Must Implement initialise method")},
        update: () => {throw new Error("Must Implement update method")},
    }
    
    container.roomApplicationInterface = roomApplicationInterface
})(modules)