package material

type AbstractMessage struct {
	Messagetype  string `json:"messagetype"`
}

func NewAbstractMessage(messagetype string)AbstractMessage{
	return AbstractMessage{Messagetype: messagetype}
}

func (m AbstractMessage) Getmessagetype() string {
	return m.Messagetype
}


