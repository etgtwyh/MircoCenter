package execption

import (
	"encoding/json"
	"fmt"
)

func NewApiException(code int, format string, a ...any) *Exception {
	if i := cap(a); i == 0 {
		return &Exception{
			Code:     code,
			Messages: fmt.Sprintf(format),
		}
	} else {
		return &Exception{
			Code:     code,
			Messages: fmt.Sprintf(format, a),
		}
	}

}

var _ error = (*Exception)(nil)

type Exception struct {
	Code     int
	Messages string
}

func (e *Exception) Error() string {
	indent, _ := json.MarshalIndent(e, "", " ")
	return string(indent)
}
