package params

type ValidationErrorResponse struct {
	Message     string            `json:"message"`
	FieldErrors map[string]string `json:"field_errors,omitempty"`
}
