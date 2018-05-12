package entity

import (
	"encoding/json"
	"errors"
)

/*
Entity struct
 */
type Entity struct {
	Altitude float32 `json:"altitude"`
	Longitude float32 `json:"longitude"`
	Latitude float32 `json:"latitude"`
	Health float32 	`json:"health"`
	Mobile bool		`json:"mobile"`
}

/*
Validate an entities structure
 */
func (e Entity) valid () (error) {
	if (e.Health < 0) {
		return errors.New("Health must be greater than zero")
	}

	return nil
}

/*
Create new Entity from json string
 */
func FromJson(jsonStr string) (Entity, error){
	var entity Entity
	err := json.Unmarshal([]byte(jsonStr), &entity)
	if (err != nil) {
		return entity, err
	}

	err = entity.valid()
	return entity, err
}
