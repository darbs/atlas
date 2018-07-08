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
	entityJson  = "{\"Id\": \"ABC123\", \"latitude\": 30.307182, \"longitude\": -97.755996, \"altitude\": 489}"
	entityTable = "Entity"
	entityId    = "ID12345"
	testLocale  = "ABC123"
)

func init() {
	log.Printf("Emptying Entity db for testing model.Entity package")
	conf := config.GetConfig()
	database.Configure(conf.DbEndpoint, conf.DbName)
	db := database.GetDatabase()
	db.Table(entityTable).Empty()
}

/////////
// Unit
/////////

func TestEntityParserEmpty(t *testing.T) {
	t.Parallel()
	_, err := model.EntityFromJson("{}")
	if err != nil {
		t.Errorf("Failed to parse empty Object")
	}
}

func TestEntityParserValid(t *testing.T) {
	t.Parallel()
	entity, err := model.EntityFromJson(entityJson)
	if err != nil {
		fmt.Printf("Resulting entity: %v error: %v", entity, err)
		t.Errorf("Failed to parse valid Object")
	}
}

func TestEntityParserHealthInvalid(t *testing.T) {
	t.Parallel()
	entity, err := model.EntityFromJson("{\"health\": -1}")

	err = entity.Valid()
	if err == nil {
		fmt.Printf("Resulting entity: %v error: %v", entity, err)
		t.Errorf("Failed to catch negative health")
	}
}

func TestEntityParserIdInvalid(t *testing.T) {
	t.Parallel()
	entity, err := model.EntityFromJson(entityJson)
	entity.Id = ""

	err = entity.Valid()
	if err == nil {
		fmt.Printf("Resulting entity: %v error: %v", entity, err)
		t.Errorf("Failed to catch missing id")
	}
}

/////////////////
// Integration //
/////////////////

func TestEntitySaveIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestEntitySaveIntegration")
	}

	entity := model.Entity{
		Id:        entityId,
		LocaleId:  testLocale,
		Altitude:  4567,
		Longitude: 1234,
		Latitude:  1234,
		Health:    100,
		Mobile:    false,
		Timestamp: time.Now().UTC(),
	}

	err := entity.Save()

	if err != nil {
		t.Errorf("Failed to save valid Entity")
	}
}

func TestEntityFindByIdIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestEntitySaveIntegration")
	}

	entity, err := model.GetEntityById(entityId)
	if err != nil {
		t.Errorf("Failed to retrieve valid test Entity")
	}

	if entity.Id != entityId {
		t.Errorf("Failed to retrieve matching test Entity")
	}
}

func TestGetEntitiesAtLocaleIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestEntityFindByIdIntegration")
	}

	entity := model.Entity{
		Id:        entityId + "-1",
		LocaleId:  testLocale,
		Altitude:  4567,
		Longitude: 1234,
		Latitude:  1234,
		Health:    100,
		Mobile:    false,
		Timestamp: time.Now().UTC(),
	}
	err := entity.Save()

	entities, err := model.GetEntitiesAtLocale(testLocale)
	if err != nil {
		t.Errorf("Failed to query locale for entities")
	}

	if len(entities) != 2 {
		t.Errorf("Failed to retrieve entities at local")
	}

	if entities[0].LocaleId != testLocale {
		t.Errorf("Failed to retrieve entities at correct local")
	}
}
