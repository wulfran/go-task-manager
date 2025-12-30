package requests

type ValidationResult struct {
	Validated bool
	Message   string
}

func (v *ValidationResult) SetFailed(m string) {
	v.Validated = false
	v.Message += m + ", "
}
