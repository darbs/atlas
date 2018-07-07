package model

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/darbs/barbatos-fwk/config"
	"github.com/darbs/barbatos-fwk/database"
	log "github.com/sirupsen/logrus"
)

var (
	localeTable = "Locale"
)

// Locale struct
type Locale struct {
	Id        string    `json:"id"`
	Name      string    `json:"name"`
	Area      int32     `json:"area"` // todo maybe large area? and should this be updateble?
	Active    bool      `json:"active"`
	Timestamp time.Time `json:"timestamp"`
}

// Validate a locale structure
func (l Locale) Valid() error {
	if l.Id == "" {
		return fmt.Errorf("locale must have an Id")
	}

	if l.Name == "" {
		return fmt.Errorf("locale must have an Name")
	}

	// restrict the size to something reasonable
	if l.Area < 0 {
		return fmt.Errorf("locale must have an Area")
	}

	if l.Timestamp.IsZero() {
		return fmt.Errorf("locale must have an Timestamp")
	}

	return nil
}

// Save locale
func (l Locale) Save() error {
	err := l.Valid()
	if err != nil {
		return err
	}

	table := database.GetDatabase().Table(localeTable)

	// todo maybe return the changelog?
	_, err = table.Upsert(database.Query{"id":l.Id, "name": l.Name}, l)
	if err != nil {
		return err
	}

	return nil
}

// Create new Locale from json string
func LocaleFromJson(jsonStr string) (Locale, error) {
	var locale Locale
	err := json.Unmarshal([]byte(jsonStr), &locale)
	if err != nil {
		return locale, err
	}

	return locale, err
}

// Get locale by Id
func GetLocaleById(localeId string) (Locale, error) {
	table := database.GetDatabase().Table(localeTable)

	var result []Locale
	err := table.Find(database.Query{"id": localeId}, &result, -1)
	if len(result) == 1 {
		return result[0], err
	}

	return Locale{}, err
}

// Get locale by id and Name
func GetLocaleByIdAndName(localeId string, localeName string) (Locale, error) {
	table := database.GetDatabase().Table(localeTable)

	var result []Locale
	err := table.Find(database.Query{"id": localeId, "name": localeName}, &result, -1)
	if len(result) == 1 {
		return result[0], err
	}

	return Locale{}, err
}

// Initialize Atlas locale
func init() {
	log.Info("Initializing Atlas - Locale")

	conf := config.GetConfig()
	database.Configure(conf.DbEndpoint, conf.DbName)

	// Index
	index := database.Index{
		Key:        []string{"id", "name"},
		Unique:     true,
		Dups:       false,
		Background: true,
		Sparse:     true,
	}

	table := database.GetDatabase().Table(localeTable)
	err := table.Index(index)
	if err != nil {
		panic(err)
	}
}
