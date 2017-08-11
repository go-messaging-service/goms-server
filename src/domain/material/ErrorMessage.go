package material

type ErrorMessage struct {
	Messagetype  string `json:"messagetype"`
	Errorcode  string `json:"errorcode"`
	Error  string `json:"error"`
}

func NewErrorMessage(messagetype string, errorcode string, error string)ErrorMessage{
	return ErrorMessage{Messagetype: messagetype, Errorcode: errorcode, Error: error}
}

func (m ErrorMessage) Getmessagetype() string {
	return m.Messagetype
}

func (e ErrorMessage) Geterrorcode() string {
	return e.Errorcode
}

func (e ErrorMessage) Geterror() string {
	return e.Error
}


