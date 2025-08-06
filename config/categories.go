package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Category represents a report category
type Category struct {
	ID          string `yaml:"id" json:"id"`
	Name        string `yaml:"name" json:"name"`
	Icon        string `yaml:"icon" json:"icon"`
	Description string `yaml:"description" json:"description"`
}

// CategorySettings holds category-related settings
type CategorySettings struct {
	RequireSubcategory bool `yaml:"require_subcategory" json:"require_subcategory"`
}

// LocationRequirements defines which categories require location
type LocationRequirements struct {
	LocationRequired []string `yaml:"location_required" json:"location_required"`
}

// FormConfigurations defines form behavior for location-required categories
type FormConfigurations struct {
	LocationRequired struct {
		MapPicker           bool `yaml:"map_picker" json:"map_picker"`
		AddressRequired     bool `yaml:"address_required" json:"address_required"`
		CoordinatesRequired bool `yaml:"coordinates_required" json:"coordinates_required"`
		ShowOnPublicMap     bool `yaml:"show_on_public_map" json:"show_on_public_map"`
	} `yaml:"location_required" json:"location_required"`
}

// TransportField represents a single transport field configuration
type TransportField struct {
	Name        string   `yaml:"name" json:"name"`
	Label       string   `yaml:"label" json:"label"`
	Type        string   `yaml:"type" json:"type"`
	Placeholder string   `yaml:"placeholder" json:"placeholder"`
	Required    bool     `yaml:"required" json:"required"`
	Options     []string `yaml:"options,omitempty" json:"options,omitempty"`
}

// TransportType represents a transport type configuration
type TransportType struct {
	Name   string           `yaml:"name" json:"name"`
	Icon   string           `yaml:"icon" json:"icon"`
	Fields []TransportField `yaml:"fields" json:"fields"`
}

// TransportConfigurations defines transport-specific configurations
type TransportConfigurations struct {
	TransportRequired []string                    `yaml:"transport_required" json:"transport_required"`
	TransportTypes    map[string]TransportType    `yaml:"transport_types" json:"transport_types"`
}

// CategoriesConfig holds the complete categories configuration
type CategoriesConfig struct {
	Categories              []Category              `yaml:"categories" json:"categories"`
	Settings                CategorySettings        `yaml:"settings" json:"settings"`
	LocationRequirements    LocationRequirements    `yaml:"location_requirements" json:"location_requirements"`
	FormConfigurations      FormConfigurations      `yaml:"form_configurations" json:"form_configurations"`
	TransportConfigurations TransportConfigurations `yaml:"transport_configurations" json:"transport_configurations"`
}

// Global variable to hold loaded categories
var CategoriesData *CategoriesConfig

// LoadCategories loads the categories configuration from YAML file
func LoadCategories() (*CategoriesConfig, error) {
	data, err := os.ReadFile("config/categories.yaml")
	if err != nil {
		return nil, fmt.Errorf("error reading categories.yaml: %w", err)
	}

	var config CategoriesConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("error parsing categories.yaml: %w", err)
	}

	// Store globally for easy access
	CategoriesData = &config

	return &config, nil
}

// GetCategoryByID returns a category by its ID
func (c *CategoriesConfig) GetCategoryByID(id string) *Category {
	for _, category := range c.Categories {
		if category.ID == id {
			return &category
		}
	}
	return nil
}

// IsLocationRequired checks if a category requires location information
func (c *CategoriesConfig) IsLocationRequired(categoryID string) bool {
	for _, id := range c.LocationRequirements.LocationRequired {
		if id == categoryID {
			return true
		}
	}
	return false
}

// IsTransportRequired checks if a category requires transport information
func (c *CategoriesConfig) IsTransportRequired(categoryID string) bool {
	for _, id := range c.TransportConfigurations.TransportRequired {
		if id == categoryID {
			return true
		}
	}
	return false
}

// GetTransportTypes returns all available transport types
func (c *CategoriesConfig) GetTransportTypes() map[string]TransportType {
	return c.TransportConfigurations.TransportTypes
}

// GetTransportType returns a specific transport type by key
func (c *CategoriesConfig) GetTransportType(transportType string) *TransportType {
	if transportType, exists := c.TransportConfigurations.TransportTypes[transportType]; exists {
		return &transportType
	}
	return nil
}

// GetCategories returns all available categories
func (c *CategoriesConfig) GetCategories() []Category {
	return c.Categories
}

// ValidateCategoryID checks if a category ID exists
func (c *CategoriesConfig) ValidateCategoryID(id string) bool {
	return c.GetCategoryByID(id) != nil
}

// Helper functions for easy access to global categories

// GetAllCategories returns all categories from the global config
func GetAllCategories() []Category {
	if CategoriesData == nil {
		return []Category{}
	}
	return CategoriesData.Categories
}

// GetCategory returns a specific category by ID from global config
func GetCategory(id string) *Category {
	if CategoriesData == nil {
		return nil
	}
	return CategoriesData.GetCategoryByID(id)
}

// IsLocationRequiredGlobal checks if location is required for a category
func IsLocationRequiredGlobal(categoryID string) bool {
	if CategoriesData == nil {
		return false
	}
	return CategoriesData.IsLocationRequired(categoryID)
}

// IsTransportRequiredGlobal checks if transport is required for a category
func IsTransportRequiredGlobal(categoryID string) bool {
	if CategoriesData == nil {
		return false
	}
	return CategoriesData.IsTransportRequired(categoryID)
}

// GetTransportTypesGlobal returns all transport types from global config
func GetTransportTypesGlobal() map[string]TransportType {
	if CategoriesData == nil {
		return map[string]TransportType{}
	}
	return CategoriesData.GetTransportTypes()
}

// GetTransportTypeGlobal returns a specific transport type from global config
func GetTransportTypeGlobal(transportType string) *TransportType {
	if CategoriesData == nil {
		return nil
	}
	return CategoriesData.GetTransportType(transportType)
}
