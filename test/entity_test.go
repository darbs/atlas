package test

import (
	"testing"
	"fmt"
	"github.com/darbs/atlas/entity"
)

func TestEntityParserEmpty (t *testing.T){
	_, err := entity.FromJson("{}")
	if (err != nil) {
		t.Errorf("Failed to parse empty Object")
	}
}

func TestEntityParserValid (t *testing.T){
	entity, err := entity.FromJson("{\"latitude\": 30.307182, \"longitude\": -97.755996, \"altitude\": 489}")
	if (err != nil) {
		fmt.Printf("Resulting entity: %v error: %v", entity, err)
		t.Errorf("Failed to parse valid Object")
	}
}

func TestEntityParserInvalid (t *testing.T){
	entity, err := entity.FromJson("{\"health\": -1}")
	if (err == nil) {
		fmt.Printf("Resulting entity: %v error: %v", entity, err)
		t.Errorf("Failed to catch negative health")
	}
}