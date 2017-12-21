# Preamble
This messaging Service is pretty similar to the ActiveMQ broker from Apache. This one is written in go and is very simple structured and therefore a fast and easy to use messaging service.

# Basic stuff
## Topics
This service is based in the so called `topics` that are familiar since ActiveMQ used them too.
A topics is just a key on which you can send or receive messages to/from.

Clients can register them self to a topic and is allowed to receive and send message to it.

All topic names are __lowercase__ letters! A registration with some __capital__ letters (e.g. camel-case like "aCoolTopicName") will be __changed into lowercase__ names.

## Messages
A message is the whole data send by the server or the client.
The data in the message is sometimes also calles _message_, because it's the most important thing, but is in the `goMS` world just the message-data.

The message-data can be everything and of (currently) any size. Normally it would be some XML, JSON or *insert-your-favorite-file-format-here* data and thats fine.

In the JSON-message (see below) will everything be escaped, but you can send everything you like.

## Message types
There're different types of messages, namely `register`, `send`, `message`, `logout`, `close` and `error`.

### register
_Client &#8594; Server_

Need the field `topics` with a list of topics. A client can only connect to topics, that are configured in the server config.

### send
_Client &#8594; Server_

Needs the fields `topics` and `data`. The field `topics` contains all topics the data goes to. The field `data` contains all data that will be send as message to all registered clients.

### message
_Server &#8594; Client_

Only contains the `data` field from the `send` message and also the `topics` the data belongs to.

### logout
_Client &#8594; Server_

Needs the field `topics` which is a list of all topics the client wants to unregister himself of.

### close
_Client &#8594; Server_

Needs no fields.

### error
_Server &#8594; Client_

Contains the field `error-code` which is normally a number (like the HTTP status codes) and the field `error` which contains some data belonging to the error (meaning of error codes below).

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

### General
```json
{
  "topic-config":"path/to/topics.json",
  "connectors":[
    {
      "protocoll":"tcp",
      "ip":"127.0.0.1",
      "port":55545
    }
  ]
}
```
The following fields need to be specified:

| Field | Type | Description |
|:---|:---|:---|
| topic-config | string | Path to the `topics.json` file |
| connectors | list (s. below) | List of connectors, on which the server is available |

### The Connector
A connector is a specification of a listener, the server is listening on.

| Field | Type | Description |
|:---|:---|:---|
| protocoll | string | The protocoll of the connector. This can be `tcp` or `udp`. |
| ip | string | The IP address the server is listening on. |
| port | number | The port, which is `55545` as default.

# Connect with server
The process of connecting and notifying is described below.

It's very important to do these steps in the __given order__, otheriwse your request will be ignored.
For example: If you send a message to a topic before register yourself to it, your request will be ignored.

### 1.) Connecting to server
1. Client creates normal TCP-connection on the port specified in the connectors of the server configuration (default is `55545`)

### 2.) Register to a topic
1. Client sends the topics as JSON-list to the server:
```json
{
  "messagetype": "register",
  "topics": [
    "some",
    "topics"
  ]
}
```
2. Server saves the client in his internal map

Maybe there'll be an acknowledgement from the server (__not implemented yet__):
3. Server sends notification that everything is ready:
```json
{
  "messagetype": "reg-ack"
}
```
It's also possible (for error-correction) to send the list of topics within the acknowledgement (but as mentioned above: __not yet implemented__):
```json
{
  "messagetype": "reg-ack",
  "topics": [
    "some",
    "topics"
  ]
}
```

### 3.) Distribute messages
1. Client sends a `send` request of the following form to the server:
```json
{
  "messagetype": "send",
  "topics": [
    "some",
    "topics"
  ],
  "data": "This will be the data to send."
}
```

2. The server will then look into his map to determine all connections to clients which registeres themselves, and then send the message. The message will have this simple format:
```json
{
  "messagetype": "message",
  "topics": [
    "some",
    "topics"
  ],
  "data": "This is the sent message."
}
```
There will be no acknowledgement from the client about the messages. We trust in the TCP protocoll and the client/server implmenetation.

### 4.) Logout from a topic
Just send the `logout`-message:
1. Client removes himself from some topics.
```json
{
  "messagetype": "logout",
  "topics": [
    "some",
    "topics"
  ]
}
```
The `topics` list is __optional__, leaving it out will logout the client from __all__ topics.

### 5.) Close connection
If you want to be kind to the server, you can use the `close`-message:
1. Client closes connection, this will remove him from all topics (of course).
```json
{
  "messagetype": "close"
}
```

# Errors
Now some words about the `error` message (s. above).

This list of numbers and codes may not be up-to-date and also may change very quickly, so don't wonder about differeces.

## Categories
To structure the whole thing, each message has its own category.

| Error code | Category |
|:---:|:---|
| 000xxx | General Server error |
| 001xxx | `register` error |
| 002xxx | `send` error |
| 003xxx | `logout` error |
| 004xxx | `close` error |

## Error-code list
### 000
### 001
| Error code   | Describtion | The field `Error` contains ... |
|:---:|:---|:---|
| 001001 | Registration not allowed, the topic doesn't exist in the server config. | ... a list (normal string separated by comma `,`) of all topics the client was not able to register to. |
| 001002 | The client was already registered in on of the topics. | ... a list (normal string separated by comma `,`) of all topics the client was already registered to.

### 002
| Error code   | Describtion | The field `Error` contains ... |
|:---:|:---|:---|
| 002001 | Error sending the message | ... the error message the runtime gives to the server. |
### 003
### 004

# Planned things
* Create users (with optional password) for server, to allow multiple topics with the same name within a single server. A user can create an account at the service provider (e.g. a website that runs this server) and then setup his/her own message service.
* ~~Do not allow the creation of new topics (only via config-file)~~ _implemented_
* Every user has it's own directory with own config file
* Cache huge files (e.g. images or huge XML files) and just send a reference. The client then can download it when needed. Therefore new message-types must be created (`reference` and `download`). The files are cahced withing the usrers directory
