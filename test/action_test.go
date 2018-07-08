package test

import (
	"fmt"
	"log"
	"testing"

	"github.com/darbs/atlas/internal"
	"github.com/darbs/atlas/model"
	"github.com/darbs/barbatos-fwk/config"
	"github.com/darbs/barbatos-fwk/database"
)

func init () {
	log.Printf("Emptying Locale db for testing action package")
	conf := config.GetConfig()
	database.Configure(conf.DbEndpoint, conf.DbName)
	db := database.GetDatabase()
	db.Table(localeTable).Empty()
}

/////////////////
// Integration //
/////////////////

func TestActionOpenLocaleNewIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestActionOpenLocaleIntegration")
	}

	locale := model.Locale{
		Name: "Test",
		Active: false,
		Area: 345,
	}
	response := internal.ActionHandler(internal.OpenLocale, locale)
	if response.Type != internal.ActionStatusSuccess {
		fmt.Printf("Resulting response: %v", response)
		t.Errorf("Failed to open valid locale")
	}

	actionLocale := response.Data.(model.Locale)
	if actionLocale.Id == "" {
		t.Errorf("Locale was not updated")
	}

	if actionLocale.Name != locale.Name {
		t.Errorf("Locale name was not updated")
	}

	if actionLocale.Area != locale.Area {
		t.Errorf("Locale area was not updated")
	}

	if actionLocale.Active == false {
		t.Errorf("Failed to open valid locale")
	}
}

func TestActionOpenLocaleValidIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestActionOpenLocaleIntegration")
	}

	locale := model.Locale{
		Id: "ABC-123-!@#",
		Name: "Test",
		Active: false,
		Area: 345,
	}
	response := internal.ActionHandler(internal.OpenLocale, locale)
	if response.Type != internal.ActionStatusSuccess {
		fmt.Printf("Resulting response: %v", response)
		t.Errorf("Failed to open valid locale")
	}

	actionLocale := response.Data.(model.Locale)
	if actionLocale.Id != locale.Id {
		t.Errorf("Locale was not updated")
	}

	dbLocale, _ := model.GetLocaleByIdAndName(locale.Id, locale.Name)
	if locale.Id != dbLocale.Id || locale.Name != dbLocale.Name {
		t.Errorf("Opened locale was not persisted in the database")
	}
}

func TestActionCloseLocaleIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestActionCloseLocaleIntegration")
	}

	locale := model.Locale{
		Id: "ABC-123-!@#",
		Name: "Test",
		Active: true,
		Area: 345,
	}
	response := internal.ActionHandler(internal.CloseLocale, locale)
	if response.Type != internal.ActionStatusSuccess {
		fmt.Printf("Resulting response: %v", response)
		t.Errorf("Failed to open valid locale")
	}

	actionLocale := response.Data.(model.Locale)
	if actionLocale.Id != locale.Id {
		t.Errorf("Locale was not updated")
	}

	if actionLocale.Name != locale.Name {
		t.Errorf("Locale name was not updated")
	}

	if actionLocale.Active == true {
		t.Errorf("Failed to open valid locale")
	}

	dbLocale, _ := model.GetLocaleByIdAndName(locale.Id, locale.Name)
	if locale.Id != dbLocale.Id || locale.Name != dbLocale.Name {
		t.Errorf("Opened locale was not persisted in the database")
	}
}

func TestActionCloseLocaleMissingIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestActionCloseLocaleIntegration")
	}

	locale := model.Locale{
		Id: "ABC-123-DOES_NOT_WORK",
		Name: "Test",
		Active: true,
		Area: 345,
	}
	response := internal.ActionHandler(internal.CloseLocale, locale)
	if response.Type == internal.ActionStatusSuccess {
		fmt.Printf("Resulting response: %v", response)
		t.Errorf("Failed to open valid locale")
	}

	dbLocale, _ := model.GetLocaleByIdAndName(locale.Id, locale.Name)
	if locale.Id == dbLocale.Id || locale.Name == dbLocale.Name {
		t.Errorf("Opened locale was not persisted in the database")
	}
}