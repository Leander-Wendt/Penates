package dto_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Leander-Wendt/Penates/backend/internal/dto"
	"github.com/Leander-Wendt/Penates/backend/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var restrictedKeys = []string{
	"manufacturer",
	"serialNumber",
	"lastTechnicalInspectionDate",
	"lastElectricalInspectionDate",
	"dateOfPurchase",
	"price",
	"resolutionNumber",
}

func sampleItem() models.Item {
	now := time.Now()
	return models.Item{
		InventoryNumber:              "123456",
		Name:                         "Drill",
		Category:                     "Tools",
		Amount:                       3,
		Manufacturer:                 "Bosch",
		SerialNumber:                 "SN-1",
		LastTechnicalInspectionDate:  &now,
		LastElectricalInspectionDate: &now,
		DateOfPurchase:               &now,
		Price:                        99.99,
		ResolutionNumber:             "RES-1",
	}
}

func toJSONMap(t *testing.T, v interface{}) map[string]interface{} {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(b, &m))
	return m
}

func TestNewItemDTO_StudentRedactsRestrictedFields(t *testing.T) {
	item := sampleItem()
	result := dto.NewItemDTO(item, models.RoleStudent)
	m := toJSONMap(t, result)
	for _, key := range restrictedKeys {
		_, present := m[key]
		assert.Falsef(t, present, "expected key %q to be absent for student role", key)
	}
	assert.Equal(t, "123456", m["inventoryNumber"])
	assert.Equal(t, "Drill", m["name"])
}

func TestNewItemDTO_AdminIncludesRestrictedFields(t *testing.T) {
	item := sampleItem()
	result := dto.NewItemDTO(item, models.RoleAdmin)
	m := toJSONMap(t, result)
	for _, key := range restrictedKeys {
		_, present := m[key]
		assert.Truef(t, present, "expected key %q to be present for admin role", key)
	}
}

func TestNewItemDTO_LogisticsIncludesRestrictedFields(t *testing.T) {
	item := sampleItem()
	result := dto.NewItemDTO(item, models.RoleLogistics)
	m := toJSONMap(t, result)
	for _, key := range restrictedKeys {
		_, present := m[key]
		assert.Truef(t, present, "expected key %q to be present for logistics role", key)
	}
}
