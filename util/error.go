package util

import "encoding/json"

var (
	NoError = NewError(0, "")
)

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type Reply struct {
	Error
	Content interface{} `json:"content"`
}

func (e *Error) String() string {
	b, err := json.Marshal(e)
	if err != nil {
		return err.Error()
	}
	return String(b)
}

func (e *Error) Error() string {
	return e.String()
}

func NewError(code int, msg string) Error {
	return Error{code, msg}
}

func NewReply(body interface{}) Reply {
	return Reply{Error: NoError, Content: body}
}
