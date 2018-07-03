package model

import (
	"encoding/json"
	"log"

	"github.com/darbs/barbatos-fwk/config"
	"github.com/darbs/barbatos-fwk/database"
)

var (
	localeTable = "Locale"
)

// Locale struct
type Locale struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Area int32  `json:"area"` // todo maybe large area?
}

// Validate a locale structure
func (e Locale) Valid() error {
	return nil
}

// Save locale
func (e Locale) Save() error {
	err := e.Valid()
	if err != nil {
		return err
	}

	table := database.GetDatabase().Table(localeTable)
	err = table.Insert(e)
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

// Initialize Atlas locale
func init() {
	log.Println("Initializing Atlas - Locale")

	conf := config.GetConfig()
	database.Configure(conf.DbEndpoint, conf.DbName)

	// Index
	index := database.Index{
		Key:        []string{"name"},
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
