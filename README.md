# [Change Project Name] Socket Server 

This application will provide means of authenticating and creating web-sockets from client requests. It provides functionality for allowing sockets to enter rooms with where by communication with other sockets is facilatated. The Intent of the application is to provide a framework in which to build real time socket appliucation around. This server will pass messages onto a an application that will handle the specific application payloads. IOn theory this server will be able to be used for simple chat application to multiplayer games.

## Motivation
This project begun as away of experimenting with GO Lang's concurrency mechanisms and exploration the langauge further. .... write more about the actual project and gaming side.


## Build status
Build status of continus integration i.e. travis, appveyor etc. Ex. - 

[![Build Status](https://travis-ci.org/akashnimare/foco.svg?branch=master)](https://travis-ci.org/akashnimare/foco)
[![Windows Build Status](https://ci.appveyor.com/api/projects/status/github/akashnimare/foco?branch=master&svg=true)](https://ci.appveyor.com/project/akashnimare/foco/branch/master)

## Code style
If you're using any code style like xo, standard etc. That will help others while contributing to your project. Ex. -

[![js-standard-style](https://img.shields.io/badge/code%20style-standard-brightgreen.svg?style=flat)](https://github.com/feross/standard)
 
## Screenshots
Include logo/demo screenshot etc.

## Tech/framework used
Ex. -

<b>Built with</b>
- [Electron](https://electron.atom.io)

## Features
What makes your project stand out?

## Domain 
Explain the relation ship with diagrams of Rooms, Clients, Users

## Preformance

Test results

# Code Example

## Example 1 - GoLang Chat Server
A simple Go chat server. In this example the Socket server will be used a module with a go application. The actual application that will be handling the payloads is passed to the application before starting.

This is a simple struct that implements the "something something" interface. The Application on the backend will simply decode the payload and count the number of words used in the each message for each user, this will be sent back to the frontend for the frontend to display.

The frontend client will be the simple-chat JS application provided by this project that allows for the user to enter the application, enter a room or conversation with someone. We will update this in the example to use the updated payload by the new go-lang application.

......

## Example 2 - Simple Node JS Pong Game

In this example our socket server will be set up to point to a local NodeJS server, This will compuncate of http. The Socket Server will handle all the client request as usual as well as maintain a connection to the node server. This will also be a websocket server. The payload will be passed the node js server which will manage its state.

This game server will use the room ID to maintain previous state for that particlar game.

The frontend client will be the pong game js server.

## Installation
Provide step by step series of examples and explanations about how to get a development env running.

## API Reference

Depending on the size of the project, if it is small and simple enough the reference docs can be added to the README. For medium size to larger projects it is important to at least provide a link to where the API reference docs live.

## Tests
Describe and show how to run the tests with code examples.

## How to use?
If people like your project they’ll want to learn how they can use it. To do so include step by step guide to use your project.

## Contribute

Let people know how they can contribute into your project. A [contributing guideline](https://github.com/zulip/zulip-electron/blob/master/CONTRIBUTING.md) will be a big plus.

## Credits
Give proper credits. This could be a link to any repo which inspired you to build this project, any blogposts or links to people who contrbuted in this project. 

#### Anything else that seems useful

## License
A short snippet describing the license (MIT, Apache etc)

MIT © [Yourname]()
