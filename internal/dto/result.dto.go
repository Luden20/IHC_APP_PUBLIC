package dto

type Result struct {
	status  string
	message string
	id      string
}

func ErrorResult(message string) Result {
	return Result{
		status:  "error",
		message: message,
		id:      "",
	}
}
func SucessResult(message string, id string) Result {
	return Result{
		status:  "success",
		message: message,
		id:      id,
	}
}
func FromErrorResult(err error) Result {
	return Result{
		status:  "error",
		message: err.Error(),
		id:      "",
	}
}
func ToMap(e *Result) map[string]interface{} {
	return map[string]interface{}{
		"status":  e.status,
		"message": e.message,
		"id":      e.id,
	}
}
