package syntax

import "fmt"

type Error struct {
	Msg      string
	Position Position
}

// Error implements error.
func (e *Error) Error() string {
	return fmt.Sprintf("%v %s", e.Position, e.Msg)
}

func NewError(pos Position, msg string) *Error {
	return &Error{Msg: msg, Position: pos}
}

var (
	_ error = (*Error)(nil)
)
