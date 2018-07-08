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
	entityTable = "Entity"
)

// Entity struct
// Note: 	when querying the attributes defined on the map
// 			must equal the json parsed attributes
type Entity struct {
	Id        string    `json:"id"`
	LocaleId  string    `json:"localeId"`
	Friendly  bool      `json:"ally"`
	Altitude  float32   `json:"altitude"`
	Longitude float32   `json:"longitude"`
	Latitude  float32   `json:"latitude"`
	Health    float32   `json:"health"`
	Mobile    bool      `json:"mobile"`
	Timestamp time.Time `json:"timestamp"`
}

// Initialize Atlas entity
func init() {
	log.Info("Initializing Atlas - Entity")

	conf := config.GetConfig()
	database.Configure(conf.DbEndpoint, conf.DbName)

	// Index
	index := database.Index{
		Key:        []string{"Id", "LocaleId"},
		Unique:     true,
		Dups:       false,
		Background: true,
		Sparse:     true,
	}

	table := database.GetDatabase().Table(entityTable)
	err := table.Index(index)
	if err != nil {
		panic(err)
	}
}

// Validate an entities structure
func (e Entity) Valid() error {
	if e.Health < 0 {
		return fmt.Errorf("health must be greater than zero")
	}

	if e.Id == "" {
		return fmt.Errorf("entity must have an Id")
	}

	if e.LocaleId == "" {
		return fmt.Errorf("entity must have a locale")
	}

	if e.Timestamp.IsZero() {
		return fmt.Errorf("entity must have an Timestamp")
	}

	return nil
}

// Save entity
func (e Entity) Save() error {
	err := e.Valid()
	if err != nil {
		return err
	}

	table := database.GetDatabase().Table(entityTable)
	err = table.Insert(e)
	if err != nil {
		return err
	}

	return nil
}

// Get entities at a corresponding entity locale
func (e Entity) GetLocalEntities() ([]Entity, error) {
	return GetEntitiesAtLocale(e.LocaleId)
}

// Get all entities at a locale
func GetEntitiesAtLocale(localeId string) ([]Entity, error) {
	table := database.GetDatabase().Table(entityTable)

	var result []Entity
	err := table.Find(database.Query{"localeid": localeId}, &result, -1)
	return result, err
}

// Get Entity by Id
func GetEntityById(id string) (Entity, error) {
	table := database.GetDatabase().Table(entityTable)

	var result []Entity
	err := table.Find(database.Query{"id": id}, &result, -1)
	if len(result) == 1 {
		return result[0], err
	}

	return Entity{}, err
}

// Create new Entity from json string
func EntityFromJson(jsonStr string) (Entity, error) {
	var entity Entity
	err := json.Unmarshal([]byte(jsonStr), &entity)
	if err != nil {
		return entity, err
	}

	return entity, err
}
