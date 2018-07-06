package actions

import "fmt"

const (
	ActionStatusSuccess = "success"
	ActionStatusError = "error"
	OpenLocale = "OPEN_LOCALE"
	CloseLocale = "CLOSE_LOCALE"
)

// maybe timestamp?
// Response interface
type ActionResponse struct {
	Type string
	Data interface {}
}

type ActionError struct {
	Message string
}

/*
Notes:
Only allow entities manipulate locales they exist in
 */

// todo validate upsert locale
// cases:
// 	1. locale exists open/close
// 	2. locale does not exist
func openLocale (data interface{}) (interface{}, error) {
	return data, nil
}

// todo update close locale
// cases:
// 	1. locale exists open/close
// 	2. locale does not exist
func closeLocale (data interface{}) (interface{}, error) {
	return data, nil
}

func Handler (action string, data interface{}) ActionResponse {
	var err error
	var response interface {}

	switch action {
	case OpenLocale:
		response, err = openLocale(data)
	case CloseLocale:
		response, err = closeLocale(data)
	default:
		err = fmt.Errorf("no action defined for %v", action)
	}

	if err != nil {
		return ActionResponse{ActionStatusError, ActionError{err.Error()}}
	}

	return ActionResponse{ActionStatusSuccess, response}
}