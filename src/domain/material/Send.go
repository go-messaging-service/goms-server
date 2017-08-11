package material

type Send struct {
	Messagetype  string `json:"messagetype"`
	Topics  []string `json:"topics"`
	Data  string `json:"data"`
}

func NewSend(messagetype string, topics []string, data string)Send{
	return Send{Messagetype: messagetype, Topics: topics, Data: data}
}

func (m Send) Getmessagetype() string {
	return m.Messagetype
}

func (t Send) Gettopics() []string {
	return t.Topics
}

func (d Send) Getdata() string {
	return d.Data
}


