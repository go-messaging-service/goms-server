package material

type Message struct {
	Messagetype  string `json:"messagetype"`
	Topics  []string `json:"topics"`
	Data  string `json:"data"`
}

func NewMessage(messagetype string, topics []string, data string)Message{
	return Message{Messagetype: messagetype, Topics: topics, Data: data}
}

func (m Message) Getmessagetype() string {
	return m.Messagetype
}

func (t Message) Gettopics() []string {
	return t.Topics
}

func (d Message) Getdata() string {
	return d.Data
}


