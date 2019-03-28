package msg

type Logout struct {
	Messagetype  string `json:"messagetype"`
	Topics  []string `json:"topics"`
}

func NewLogout(messagetype string, topics []string)Logout{
	return Logout{Messagetype: messagetype, Topics: topics}
}

func (m Logout) Getmessagetype() string {
	return m.Messagetype
}

func (t Logout) Gettopics() []string {
	return t.Topics
}


