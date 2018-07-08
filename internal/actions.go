package internal

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/darbs/atlas/model"
	"github.com/darbs/barbatos-fwk/database"

	//"github.com/darbs/atlas/model"
	log "github.com/sirupsen/logrus"
)

const (
	ActionStatusSuccess = "success"
	ActionStatusError   = "error"
	OpenLocale          = "OPEN_LOCALE"
	CloseLocale         = "CLOSE_LOCALE"
)

var (
	activeLocales map[string]interface{}
)

// maybe timestamp?
// Response interface
type ActionResponse struct {
	Type string
	Data interface{}
}

// Action Error
type ActionError struct {
	Message string
}

/*
Notes:
Only allow entities manipulate locales they exist in
 */

// parse string map into local struct
func parseLocale(data interface{}) (model.Locale, error) {
	var locale model.Locale

	mData, err := json.Marshal(data)
	if err != nil {
		return locale, fmt.Errorf("failed to parse locale from data: %v", err)
	}

	err = json.Unmarshal(mData, &locale)
	if err != nil {
		return locale, fmt.Errorf("failed to un parse locale from data: %v", err)
	}

	return locale, nil
}

// open locale
func openLocale(data interface{}) (interface{}, error) {
	locale, err := parseLocale(data)
	if err != nil {
		return locale, err
	}

	if locale.Id == "" {
		locale.Id = database.GetNewObjectId()
	}

	if locale.Timestamp.IsZero() {
		locale.Timestamp = time.Now().UTC()
	}

	locale.Active = true
	log.Debugf("opening locale: %v", locale)

	err = locale.Save()
	if err == nil && activeLocales[locale.Id] == nil {
		//activeLocales[locale.Id] = make(chan string)

	}

	return locale, err
}

// close locale
func closeLocale(data interface{}) (interface{}, error) {
	locale, err := parseLocale(data)
	if err != nil {
		return locale, err
	}

	locale, err = model.GetLocaleByIdAndName(locale.Id, locale.Name)
	if err != nil {
		return locale, err
	}

	locale.Active = false
	log.Debugf("closing locale: %v", locale)

	err = locale.Save()
	if activeLocales[locale.Id] != nil {
		delete(activeLocales, locale.Id)
	}

	return locale, err
}

// action handler
func ActionHandler(action string, data interface{}) ActionResponse {
	var err error
	var response interface{}

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
