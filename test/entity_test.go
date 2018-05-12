package test

import (
	"testing"
	"github.com/darbs/atlas/entity"
)

func TestEntityParser (t *testing.T){
	_, err := entity.FromJson("{}")
	if (err != nil) {
		t.Errorf("Failed to parse empty Object")
	}
}