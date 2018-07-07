package test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/darbs/atlas/model"
	"github.com/darbs/barbatos-fwk/config"
	"github.com/darbs/barbatos-fwk/database"
)

var (
	localeTable = "Locale"
	localeId    = "LOCALE_ABC_123"
	localeJson = "{\"Id\": \"" + localeId + "\", \"name\": \"some locale\"}"
)

func init () {
	log.Printf("Emptying Locale db for testing model.Locale package")
	conf := config.GetConfig()
	database.Configure(conf.DbEndpoint, conf.DbName)
	db := database.GetDatabase()
	db.Table(localeTable).Empty()
}

/////////
// Unit
/////////

func TestLocaleParserEmpty(t *testing.T) {
	t.Parallel()
	_, err := model.LocaleFromJson("{}")
	if err != nil {
		t.Errorf("Failed to parse empty Object")
	}
}

func TestLocaleParserValid(t *testing.T) {
	t.Parallel()
	locale, err := model.LocaleFromJson(localeJson)
	if err != nil {
		fmt.Printf("Resulting locale: %v error: %v", locale, err)
		t.Errorf("Failed to parse valid Object")
	}
}

func TestLocaleParserHealthInvalid(t *testing.T) {
	t.Parallel()
	locale, err := model.LocaleFromJson("{\"health\": -1}")

	err = locale.Valid()
	if err == nil {
		fmt.Printf("Resulting locale: %v error: %v", locale, err)
		t.Errorf("Failed to catch negative health")
	}
}

func TestLocaleParserIdInvalid(t *testing.T) {
	t.Parallel()
	locale, err := model.LocaleFromJson(localeJson)
	locale.Id = ""

	err = locale.Valid()
	if err == nil {
		fmt.Printf("Resulting locale: %v error: %v", locale, err)
		t.Errorf("Failed to catch missing id")
	}
}

/////////////////
// Integration //
/////////////////

func TestLocaleSaveIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestLocaleSaveIntegration")
	}

	locale := model.Locale{
		Id:   localeId,
		Name: "some locale",
		Timestamp: time.Now().UTC(),
	}

	err := locale.Save()
	if err != nil {
		log.Printf("%v", err)
		t.Errorf("Failed to save valid Locale")
	}
}

func TestGetLocaleByIdIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestLocaleFindByIdIntegration")
	}

	id := localeId + "-1"
	newLocale := model.Locale{
		Id:   id,
		Name: "mt. fuji",
		Timestamp: time.Now().UTC(),
	}
	err := newLocale.Save()

	locale, err := model.GetLocaleById(id)
	if err != nil {
		t.Errorf("Failed to query locale for entities")
	}

	if locale.Id != id {
		t.Errorf("Failed to retrieve locale")
	}
}