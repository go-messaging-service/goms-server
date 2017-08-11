package material

type Register struct {
	Messagetype  string `json:"messagetype"`
	Topics  []string `json:"topics"`
}

func NewRegister(messagetype string, topics []string)Register{
	return Register{Messagetype: messagetype, Topics: topics}
}

func (m Register) Getmessagetype() string {
	return m.Messagetype
}

func (t Register) Gettopics() []string {
	return t.Topics
}


