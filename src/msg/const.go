package msg

/*

The types of messages are:

register(topics) 	: Will register the client to a topic
send(topics, data)	: Sends the data to all subscribers of the topics
message(data)		: The receiving message a client gets
logout(topics)		: The clients unsubscribes from the given topics
close()				: The client closes the connection and unsubscribes from all topics

The optionel reg-ack() is left out, because it's not clear wether this will be implemented

*/

// All message types (mt)
const (
	// client -> server
	MT_REGISTER = "register"
	MT_SEND     = "send"
	MT_LOGOUT   = "logout"
	MT_CLOSE    = "close"
	// server -> client
	MT_MESSAGE = "message"
	MT_ERROR   = "error"
)

// All error codes (err)
const (
	//000
	//001
	ERR_REG_INVALID_TOPIC      = "001001" // registration on topics forbidden
	ERR_REG_ALREADY_REGISTERED = "001002" // already registered in (some) topics
	//002
	ERR_SEND_FAILED = "002001" // internal server error sending message
	//003
	//004
)
