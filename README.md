# Trade wars server

## Introduction

This go applications acts as the server for coordinating gameplay events between clients. Every action takes place through a series of messages exchanged via the websocket connection initiated at the start of gameplay.

## Game state organization

Game state is stored using a system common to game development called an Entity Component System. Using an ECS allows for the server to keep track of many seperate entities with their attributes, in an extendable way. See the [wikipedia](https://en.wikipedia.org/wiki/Entity_component_system) page for more information, however the general idea of an ECS is that state is broken up into three seperate data types, entities, components, and systems. An entity is simply an ID used to represent an object in the game world. Components hold state information and can be assigned to entities, for example a PhysicsComponent might hold x, y, z coordinates, plus velocity in their respective directions. This PhysicsComponent can be attached to any entity that needs to respond to physics. The final piece of the puzzle is the system, which is also slightly abused in this implementation. Systems contain the game logic that is applicable to their situation, such as a PhysicsSystem calculating gravity. Every system has an update function and a reference to every entity that contains the components the system requires. The update function is ideally called every frame, during which it iterates through all entities avilable to it and applies the require state changes. However in this implementation there is no update tick, so the systems hold functions to modify state that can be called upon messages being recieved.

## Important files

- internal/tradewars/gameEndpoints.go holds most of everything that has to do with websocket connections. Look to decodeMessage for a general overview
- internal/tradewars/map.go MapSystem that handles player position and functions to change state.

## Other notes

- This application has A LOT of leftover components and systems from when we were planning to implement many more features then what ended up getting implemented. It's also extremely unorganized due to the constant debugging trying to get things working, sorry :(

- We did end up kind of using an eventbus, although only in one direction due to pointer problems. The subscriptions can be found in world.go and the publish events in gameEndpoints.go.

## Routes

The only route used is /gameServer. Initiate a websocket connection with this route to gain a connection and enable the rest of the server functionality.

## Messages

The following table documents every websocket command that can be used with your connection. Every message is formatted in a standard JSON based format, with a key for
the command name, and an arbitrary number of fields following to represent the arguments

Example

{"command":"setOwnPosition","x":4,"y":2}

| Command        | Description                                                                                                                               | Arguments                  | Example Response                                                                                           |
| -------------- | ----------------------------------------------------------------------------------------------------------------------------------------- | -------------------------- | ---------------------------------------------------------------------------------------------------------- |
| ping           | Ensures that heroku doesn't go to sleep while clients are connected. Must send every 10 seconds or else server will kill your connection! | None                       | { "EventType": "pong", "Target": "NULL", "EventParams": {} }                                               |
| getCallsign    | Assigns callsign to current websocket connection. MUST set callsign before using any of the below commands                                | callsign: desired callsign | { "EventType": "callsignChange", "Target": "Picard", "EventParams": { "old": "NULL" } }                    |
| setOwnPosition | Sets the players own location to the given coordinates                                                                                    | x: xPos, y: yPos           | { "EventType": "positionUpdate", "Target": "Picard", "EventParams": { "x": 3, "y": 2 } }                   |
| getOwnPosition | Returns a players own position                                                                                                            | None                       | { "EventType": "positionUpdate", "Target": "Picard", "EventParams": { "x": 3, "y": 2 } }                   |
| getAllPosition | Returns a series of events representing all connected clients' locations                                                                  | None                       | { "EventType": "positionUpdate", "Target": "Picard", "EventParams": { "x": 3, "y": 2 } }                   |
| chatMessage    | Broadcasts a chat event to all connected clients                                                                                          | message: chat message text | { "EventType": "chatMessage", "Target": "Picard", "EventParams": { "message": "This is a test message" } } |
