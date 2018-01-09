package connectionServices

import (
	"goms-server/src/domain/material"
)

type MessageHandler interface {
	Handle(message material.Message)
}
