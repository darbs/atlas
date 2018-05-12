package entity

import (
	"encoding/json"
)

type Entity struct {
	Location string `json:"location"`
	Health float32 	`json:"health"`
	Mobile bool		`json:"mobile"`
}

// todo handle empty objects or missing attributes
func FromJson(jsonStr string) (Entity, error){
	var entity Entity
	err := json.Unmarshal([]byte(jsonStr), &entity)
	return entity, err
}
