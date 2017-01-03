# Preamble
This messaging Service is pretty similar to the ActiveMQ broker from Apache. This one is written in go and is very simple structured and therefore fast messaging service.

# Basic stuff
## Topics
This service is based in the so called "topics" that are familiar since ActiveMQ used them too.
A topics is just a key on which you can send or receive messages to/from.

A client can register itself to a topic and is allowed to receive and send message to it.

All topic names are __lowercase__ letters! A registration with some capital letters (e.g. camel-case like "aCoolTopicName") will be changed into lowercase names.

## Messages
A message can be everything and of (currently) any size. Normally it would be some XML, JSON or *insert-file-format-here* data and thats fine. In the JSON-message (see below) will everything be escaped, but you can send everything you like.

## Message types
There're different types of messages, namely `register`, `send`, `message`, `logout`, `close` and `error`.

### register
Goes from client to server.

Need the field `topics` with a list of topics. A client can only connect to topics, that are configured in the server config.

### send
Goes from client to server.

Needs the fields `topics` and `data`. The field `topics` contains all topics the data goes to. The field `data` contains all data that will be send as message to all registered clients.

### message
Goes from server to client.

Only contains the `data` field from the `send` message.

### logout
Goes from client to server.

Needs the field `topics` which is a list of all topics the client wants to unregister himself of.

### close
Goes from client to server.

Needs no fields.

### error
Goes from server to client.

Contains the field `error-code` which is normally a number (like the HTTP status codes) and the field `error` which is a message that describes the error.

# Server stuff
Here're some information about the server (usage, configuration and internals).

## Configuration
### Topics
To limit the amount of topics, a client is not able to create one. This is the responsibility of the server administrator. The file `/conf/topics.json` contains all available topics.
It's a simple json list like this one:
```json
{
  "topics":[
    "technology",
    "goms",
    "golang"
  ]
}
```

## Internals
### Data structure for connections and topics
The server uses a map from `topic (string)` to `connection (net.Conn)`. The reason is, that the normal situation would be a notification to all users of a topic. This need exactly this kind of mapping for a fast distribution of messages.
# Connect with server
The process of connecting and notifying is the following:

It's very important to do these steps in the given order, otheriwse your request will be ignored. For example: If you send a message to a topic before register yourself to it, your request will be ignored.

### 1.) Connecting to server
1. Client creates normal TCP-connection on port 55545

### 2.) Register to a topic
1. Client sends the topics as JSON-list to the server:
```json
{
  "type": "register",
  "topics": [
    "some",
    "topics"
  ]
}
```
2. Server saves the client in his internal map

Maybe there'll be an acknowledgement from the server:
3. Server sends notification that everything is ready:
```json
{
  "type": "reg-ack"
}
```
It's also possible (for error-correction) to send the list of topics within the acknowledgement:
```json
{
  "type": "reg-ack",
  "topics": [
    "some",
    "topics"
  ]
}
```

### 3.) Distribute messages
1. Client sends Request of the following form to the server:
```json
{
  "type": "send",
  "topics": [
    "some",
    "topics"
  ],
  "data": "This will be the message to send."
}
```
2. The server will then look into his map to determine all connections and send the message. The message will have this simple format:
```json
{
  "type": "message",
  "data": "This is the sent message."
}
```
There will be no acknowledgement from the client about the messages. We trust in the TCP protocoll and the client/server implmenetation.

### 4.) Logout from a topic
Just send the logout-message:
1. Client removes himself from some topics.
```json
{
  "type": "logout",
  "topics": [
    "some",
    "topics"
  ]
}
```
The `topics` list is __optional__, leaving it out will close the whole connection.

### 5.) Close connection
If you want to be kind to the server, you can use the close-message:
1. Client closes connection, this will remove him from all topics (of course).
```json
{
  "type": "close"
}
```

# Errors
Now some words about the `error` message (s. above).

This list of numbers and codes may not be up-to-date and also may change very quickly, so don't wonder about differeces.

## Categories
To structure the whole thing, each message has its own category.

| Error code | Category |
|-|:-|
| 000xxx | General Server error |
| 001xxx | `register` error |
| 002xxx | `send` error |
| 003xxx | `logout` error |
| 004xxx | `close` error |

## Error-code list
### 000
### 001
| Error code   | Message |
|-|:-|
| 001001 | Registration not allowed. Maybe the topic does'nt exist in the server config? |

### 002
### 003
### 004

# Planned things
* Create users (with optional passoword) for server, to allow multiple topics with the same name within a single server. A user can create an account at the service provider (e.g. a website that runs this server) and then setup his/her own message service.
* Do not allow the creation of new topics (only via config-file)
* Every user has it's own directory with own config file
* Cache huge files (e.g. images or huge XML files) and just send a reference. The client then can download it when needed. Therefore new message-types must be created (`reference` and `download`). The files are cahced withing the usrers directory
