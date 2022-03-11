package exit

import (
	"fmt"
)

type CliExit struct {
	Err        error  `json:"err,omitempty"`
	ReturnCode int    `json:"return_code"`
	Message    string `json:"message,omitempty"`
}

var theExit *CliExit = nil

func SetExitWithError(err error, message string) {
	theExit = NewExitWithError(err, message)
}

func Clear() {
	theExit = nil
}

func Get() *CliExit {
	return theExit
}

func NewExit(err error, message string, returnCode int) *CliExit {
	return &CliExit{
		Err:        err,
		ReturnCode: returnCode,
		Message:    message,
	}
}

func NewExitWithError(err error, message string) *CliExit {
	return NewExit(err, message, 1)
}

func (e *CliExit) String() string {
	s := ""
	if e.Message == "" {
		s = fmt.Sprintf("Error: [err = %s]", e.Err)
	} else {
		if e.Err == nil {
			s = fmt.Sprintf("Error: %s", e.Message)
		} else {
			s = fmt.Sprintf("Error: %s [Err = %s ]", e.Message, e.Err)
		}
	}

	return s
}
