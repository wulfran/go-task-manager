package requests

import "fmt"

type ValidationResult struct {
	Validated bool
	Message   string
}

func (v *ValidationResult) SetFailed(m string) {
	v.Validated = false
	v.Message = fmt.Sprintf(v.Message, " ", m)
}
