package requests

type ValidationResult struct {
	Validated bool
	Message   string
}

func (v *ValidationResult) SetFailed(m string) {
	v.Validated = false
	if v.Message == "" {
		v.Message = m
	} else {
		v.Message += ", " + m
	}
}
