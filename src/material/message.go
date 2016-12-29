package material

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
	MtRegister = "register"
	MtSend     = "send"
	MtLogout   = "logout"
	MtClose    = "close"
	// server -> client
	MtMessage = "message"
)

type Message struct {
	MessageType string   `json:"type,omitempty"`
	Data        string   `json:"data,omitempty"`
	Topics      []string `json:"topics,omitempty"`
}
