package test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"github.com/darbs/atlas/model"
	"github.com/darbs/barbatos-fwk/config"
	"github.com/darbs/barbatos-fwk/database"
)

var (
	validJson = "{\"Id\": \"ABC123\", \"latitude\": 30.307182, \"longitude\": -97.755996, \"altitude\": 489}"
	tableName = "Entity"
)

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
	entity, err := model.EntityFromJson(validJson)
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
	entity, err := model.EntityFromJson(validJson)
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

var testId = "ID12345"

func TestEntitySaveIntegration(t * testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestEntitySaveIntegration")
	}

	entity := model.Entity{
		Id: testId,
		Locale:    "ABC123",
		Altitude:  4567,
		Longitude: 1234,
		Latitude:  1234,
		Health:    100,
		Mobile:    false,
	}

	err := entity.Save()

	if err != nil {
		t.Errorf("Failed to save valid Entity")
	}
}

func TestEntityFindByIdIntegration(t * testing.T) {
	if testing.Short() {
		t.Skip("Skipping TestEntitySaveIntegration")
	}

	entity, err := model.GetEntityById(testId)
	if err != nil {
		t.Errorf("Failed to retrieve valid test Entity")
	}

	if entity.Id != testId {
		t.Errorf("Failed to retrieve matching test Entity")
		log.Printf("ERROR: %v", entity)
	}
}

func setup () {
	conf := config.GetConfig()
	database.Configure(conf.DbEndpoint, conf.DbName)
	db := database.Database()
	table := db.Table(tableName)

	table.Empty()
}

func TestMain(m *testing.M) {
	setup()
	log.Println("TestMain")
	code := m.Run()

	//shutdown()
	os.Exit(code)
}