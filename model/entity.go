package entity

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/darbs/barbatos-fwk/config"
	"github.com/darbs/barbatos-fwk/database"
)

var (
	tableName  = "Entity"
	conf    config.Configuration
)

/*
Entity struct
*/
type Entity struct {
	Locale    int64   `json:"locale"`
	Altitude  float32 `json:"altitude"`
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
	Health    float32 `json:"health"`
	Mobile    bool    `json:"mobile"`
}

/*
Validate an entities structure
*/
func (e Entity) valid() error {
	if e.Health < 0 {
		return errors.New("Health must be greater than zero")
	}

	return nil
}

func (e Entity) Save() error {
	err := e.valid()
	if err != nil {
		return err
	}

	err = database.Database().Insert(tableName, &e)
	if err != nil {
		return err
	}

	return nil
}

/*
Create new Entity from json string
*/
func FromJson(jsonStr string) (Entity, error) {
	var entity Entity
	err := json.Unmarshal([]byte(jsonStr), &entity)
	if err != nil {
		return entity, err
	}

	err = entity.valid()
	return entity, err
}

func init() {
	log.Println("Initializing Atlas - Entity")

	conf = config.GetConfig()
	database.Configure(conf.DbEndpoint, conf.DbName)

	// Index
	index := database.Index{
		Key:        []string{"name", "phone"},
		Unique:     true,
		Dups:       false,
		Background: true,
		Sparse:     true,
	}

	err := database.Database().Index(tableName, index)
	if err != nil {
		panic(err)
	}
}
