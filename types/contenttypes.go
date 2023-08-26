package types

type ContentType string
const (
	JSON ContentType = "application/json"
	FORM ContentType = "application/x-www-form-urlencoded"
)