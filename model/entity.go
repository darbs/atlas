package model

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/darbs/barbatos-fwk/config"
	"github.com/darbs/barbatos-fwk/database"
)

var (
	entityTable = "Entity"
)

// Entity struct
// Note: 	when querying the attributes defined on the map
// 			must equal the json parsed attributes
type Entity struct {
	Id        string    `json:"id"`
	Locale    string    `json:"locale"`
	Ally      bool      `json:"ally"`
	Altitude  float32   `json:"altitude"`
	Longitude float32   `json:"longitude"`
	Latitude  float32   `json:"latitude"`
	Health    float32   `json:"health"`
	Mobile    bool      `json:"mobile"`
	Timestamp time.Time `json:"timestamp"`
}

// Validate an entities structure
func (e Entity) Valid() error {
	if e.Health < 0 {
		return errors.New("health must be greater than zero")
	}

	if e.Id == "" {
		return errors.New("entity must have an Id")
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
	return GetEntitiesAtLocale(e.Locale)
}

// Get all entities at a locale
func GetEntitiesAtLocale(locale string) ([]Entity, error) {
	table := database.GetDatabase().Table(entityTable)

	var result []Entity
	err := table.Find(database.Query{"locale": locale}, &result, -1)
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

	// TODO invalid entity from json
	return entity, err
}

// Initialize Atlas entity
func init() {
	log.Println("Initializing Atlas - Entity")

	conf := config.GetConfig()
	database.Configure(conf.DbEndpoint, conf.DbName)

	// Index
	index := database.Index{
		Key:        []string{"Locale"},
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
