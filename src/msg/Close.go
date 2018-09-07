package msg

type Close struct {
	Messagetype  string `json:"messagetype"`
}

func NewClose(messagetype string)Close{
	return Close{Messagetype: messagetype}
}

func (m Close) Getmessagetype() string {
	return m.Messagetype
}


