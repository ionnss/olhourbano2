package categories

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Subcategory represents a subcategory within a main category
type Subcategory struct {
	ID          string `yaml:"id" json:"id"`
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
}

// Category represents a main category with its subcategories
type Category struct {
	ID            string        `yaml:"id" json:"id"`
	Name          string        `yaml:"name" json:"name"`
	Icon          string        `yaml:"icon" json:"icon"`
	Description   string        `yaml:"description" json:"description"`
	Subcategories []Subcategory `yaml:"subcategories" json:"subcategories"`
}

// Settings holds configuration settings for categories
type Settings struct {
	RequireSubcategory    bool   `yaml:"require_subcategory" json:"require_subcategory"`
	AllowOtherSubcategory bool   `yaml:"allow_other_subcategory" json:"allow_other_subcategory"`
	MaxOtherDescription   int    `yaml:"max_other_description" json:"max_other_description"`
	DefaultCategory       string `yaml:"default_category" json:"default_category"`
}

// LocationRequirements defines which categories require, optionally use, or don't need location
type LocationRequirements struct {
	LocationRequired  []string `yaml:"location_required" json:"location_required"`
	LocationOptional  []string `yaml:"location_optional" json:"location_optional"`
	LocationNotNeeded []string `yaml:"location_not_needed" json:"location_not_needed"`
}

// FormConfiguration defines how forms should behave for different location types
type FormConfiguration struct {
	MapPicker            bool `yaml:"map_picker" json:"map_picker"`
	AddressRequired      bool `yaml:"address_required" json:"address_required"`
	CoordinatesRequired  bool `yaml:"coordinates_required" json:"coordinates_required"`
	ShowOnPublicMap      bool `yaml:"show_on_public_map" json:"show_on_public_map"`
	ShowLocationCheckbox bool `yaml:"show_location_checkbox" json:"show_location_checkbox"`
}

// FormConfigurations holds form configs for different location requirement types
type FormConfigurations struct {
	LocationRequired  FormConfiguration `yaml:"location_required" json:"location_required"`
	LocationOptional  FormConfiguration `yaml:"location_optional" json:"location_optional"`
	LocationNotNeeded FormConfiguration `yaml:"location_not_needed" json:"location_not_needed"`
}

// AnonymousReporting defines which categories allow anonymous reports
type AnonymousReporting struct {
	AnonymousAllowed       []string `yaml:"anonymous_allowed" json:"anonymous_allowed"`
	IdentificationRequired []string `yaml:"identification_required" json:"identification_required"`
}

// SensitivityLevels defines sensitivity levels for different categories
type SensitivityLevels struct {
	High   []string `yaml:"high" json:"high"`
	Medium []string `yaml:"medium" json:"medium"`
	Low    []string `yaml:"low" json:"low"`
}

// Config represents the complete categories configuration
type Config struct {
	Categories           []Category           `yaml:"categories" json:"categories"`
	Settings             Settings             `yaml:"settings" json:"settings"`
	LocationRequirements LocationRequirements `yaml:"location_requirements" json:"location_requirements"`
	FormConfigurations   FormConfigurations   `yaml:"form_configurations" json:"form_configurations"`
	AnonymousReporting   AnonymousReporting   `yaml:"anonymous_reporting" json:"anonymous_reporting"`
	SensitivityLevels    SensitivityLevels    `yaml:"sensitivity_levels" json:"sensitivity_levels"`
	OtherCategory        Category             `yaml:"other_category" json:"other_category"`
}

// LocationRequirement represents the location requirement type for a category
type LocationRequirement int

const (
	LocationRequired LocationRequirement = iota
	LocationOptional
	LocationNotNeeded
)

// SensitivityLevel represents the sensitivity level of a report category
type SensitivityLevel int

const (
	SensitivityLow SensitivityLevel = iota
	SensitivityMedium
	SensitivityHigh
)

// Manager handles category operations
type Manager struct {
	config *Config
}

// NewManager creates a new category manager
func NewManager(configPath string) (*Manager, error) {
	config, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load categories config: %w", err)
	}

	return &Manager{config: config}, nil
}

// LoadConfig loads the categories configuration from a YAML file
func LoadConfig(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &config, nil
}

// GetAllCategories returns all available categories
func (m *Manager) GetAllCategories() []Category {
	return m.config.Categories
}

// GetCategoryByID returns a category by its ID
func (m *Manager) GetCategoryByID(id string) (*Category, error) {
	for _, category := range m.config.Categories {
		if category.ID == id {
			return &category, nil
		}
	}

	// Check if it's the "other" category
	if id == m.config.OtherCategory.ID {
		return &m.config.OtherCategory, nil
	}

	return nil, fmt.Errorf("category with ID '%s' not found", id)
}

// GetSubcategoryByID returns a subcategory by its ID within a category
func (m *Manager) GetSubcategoryByID(categoryID, subcategoryID string) (*Subcategory, error) {
	category, err := m.GetCategoryByID(categoryID)
	if err != nil {
		return nil, err
	}

	for _, subcategory := range category.Subcategories {
		if subcategory.ID == subcategoryID {
			return &subcategory, nil
		}
	}

	return nil, fmt.Errorf("subcategory with ID '%s' not found in category '%s'", subcategoryID, categoryID)
}

// ValidateCategory validates if a category and subcategory combination is valid
func (m *Manager) ValidateCategory(categoryID, subcategoryID string) error {
	_, err := m.GetCategoryByID(categoryID)
	if err != nil {
		return err
	}

	// If subcategory is required but not provided
	if m.config.Settings.RequireSubcategory && subcategoryID == "" {
		return fmt.Errorf("subcategory is required for category '%s'", categoryID)
	}

	// If subcategory is provided, validate it exists
	if subcategoryID != "" {
		_, err := m.GetSubcategoryByID(categoryID, subcategoryID)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetLocationRequirement returns the location requirement for a category
func (m *Manager) GetLocationRequirement(categoryID string) LocationRequirement {
	// Check required
	for _, id := range m.config.LocationRequirements.LocationRequired {
		if id == categoryID {
			return LocationRequired
		}
	}

	// Check optional
	for _, id := range m.config.LocationRequirements.LocationOptional {
		if id == categoryID {
			return LocationOptional
		}
	}

	// Check not needed
	for _, id := range m.config.LocationRequirements.LocationNotNeeded {
		if id == categoryID {
			return LocationNotNeeded
		}
	}

	// Default to optional if not specified
	return LocationOptional
}

// GetFormConfiguration returns the form configuration for a category
func (m *Manager) GetFormConfiguration(categoryID string) FormConfiguration {
	requirement := m.GetLocationRequirement(categoryID)

	switch requirement {
	case LocationRequired:
		return m.config.FormConfigurations.LocationRequired
	case LocationOptional:
		return m.config.FormConfigurations.LocationOptional
	case LocationNotNeeded:
		return m.config.FormConfigurations.LocationNotNeeded
	default:
		return m.config.FormConfigurations.LocationOptional
	}
}

// IsAnonymousAllowed checks if anonymous reporting is allowed for a category
func (m *Manager) IsAnonymousAllowed(categoryID string) bool {
	for _, id := range m.config.AnonymousReporting.AnonymousAllowed {
		if id == categoryID {
			return true
		}
	}
	return false
}

// IsIdentificationRequired checks if identification is required for a category
func (m *Manager) IsIdentificationRequired(categoryID string) bool {
	for _, id := range m.config.AnonymousReporting.IdentificationRequired {
		if id == categoryID {
			return true
		}
	}
	return false
}

// GetSensitivityLevel returns the sensitivity level for a category
func (m *Manager) GetSensitivityLevel(categoryID string) SensitivityLevel {
	// Check high sensitivity
	for _, id := range m.config.SensitivityLevels.High {
		if id == categoryID {
			return SensitivityHigh
		}
	}

	// Check medium sensitivity
	for _, id := range m.config.SensitivityLevels.Medium {
		if id == categoryID {
			return SensitivityMedium
		}
	}

	// Default to low sensitivity
	return SensitivityLow
}

// GetCategorySelectOptions returns categories formatted for HTML select options
func (m *Manager) GetCategorySelectOptions() []SelectOption {
	var options []SelectOption

	for _, category := range m.config.Categories {
		options = append(options, SelectOption{
			Value:       category.ID,
			Label:       fmt.Sprintf("%s %s", category.Icon, category.Name),
			Description: category.Description,
		})
	}

	// Add "other" category if enabled
	if m.config.Settings.AllowOtherSubcategory {
		options = append(options, SelectOption{
			Value:       m.config.OtherCategory.ID,
			Label:       fmt.Sprintf("%s %s", m.config.OtherCategory.Icon, m.config.OtherCategory.Name),
			Description: m.config.OtherCategory.Description,
		})
	}

	return options
}

// GetSubcategorySelectOptions returns subcategories for a specific category
func (m *Manager) GetSubcategorySelectOptions(categoryID string) ([]SelectOption, error) {
	category, err := m.GetCategoryByID(categoryID)
	if err != nil {
		return nil, err
	}

	var options []SelectOption
	for _, subcategory := range category.Subcategories {
		options = append(options, SelectOption{
			Value:       subcategory.ID,
			Label:       subcategory.Name,
			Description: subcategory.Description,
		})
	}

	return options, nil
}

// SelectOption represents an option for HTML select elements
type SelectOption struct {
	Value       string `json:"value"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

// GetSettings returns the current settings
func (m *Manager) GetSettings() Settings {
	return m.config.Settings
}

// GetCategoryName returns the display name for a category
func (m *Manager) GetCategoryName(categoryID string) string {
	category, err := m.GetCategoryByID(categoryID)
	if err != nil {
		return "Categoria Desconhecida"
	}
	return category.Name
}

// GetSubcategoryName returns the display name for a subcategory
func (m *Manager) GetSubcategoryName(categoryID, subcategoryID string) string {
	subcategory, err := m.GetSubcategoryByID(categoryID, subcategoryID)
	if err != nil {
		return "Subcategoria Desconhecida"
	}
	return subcategory.Name
}

// GetFullCategoryPath returns the full path (Category > Subcategory)
func (m *Manager) GetFullCategoryPath(categoryID, subcategoryID string) string {
	categoryName := m.GetCategoryName(categoryID)

	if subcategoryID == "" {
		return categoryName
	}

	subcategoryName := m.GetSubcategoryName(categoryID, subcategoryID)
	return fmt.Sprintf("%s > %s", categoryName, subcategoryName)
}

// GetCategoriesByLocationType returns categories filtered by location requirement
func (m *Manager) GetCategoriesByLocationType(locationType LocationRequirement) []Category {
	var filteredCategories []Category

	for _, category := range m.config.Categories {
		if m.GetLocationRequirement(category.ID) == locationType {
			filteredCategories = append(filteredCategories, category)
		}
	}

	return filteredCategories
}

// GetCategoriesBySensitivity returns categories filtered by sensitivity level
func (m *Manager) GetCategoriesBySensitivity(sensitivity SensitivityLevel) []Category {
	var filteredCategories []Category

	for _, category := range m.config.Categories {
		if m.GetSensitivityLevel(category.ID) == sensitivity {
			filteredCategories = append(filteredCategories, category)
		}
	}

	return filteredCategories
}

// ShouldShowOnMap determines if a report should be displayed on public map
func (m *Manager) ShouldShowOnMap(categoryID string, hasLocation bool) bool {
	formConfig := m.GetFormConfiguration(categoryID)

	// If category doesn't support public map display
	if !formConfig.ShowOnPublicMap {
		return false
	}

	// For optional location categories, only show if location is provided
	if m.GetLocationRequirement(categoryID) == LocationOptional && !hasLocation {
		return false
	}

	// Don't show high sensitivity reports on public map
	if m.GetSensitivityLevel(categoryID) == SensitivityHigh {
		return false
	}

	return true
}

// ValidateLocationData validates location data based on category requirements
func (m *Manager) ValidateLocationData(categoryID string, hasLocation bool, latitude, longitude float64) error {
	requirement := m.GetLocationRequirement(categoryID)

	if requirement == LocationRequired && !hasLocation {
		return fmt.Errorf("location is required for category '%s'", categoryID)
	}

	if hasLocation {
		if latitude == 0 && longitude == 0 {
			return fmt.Errorf("invalid coordinates provided")
		}

		// Basic coordinate validation (adjust ranges as needed for your region)
		if latitude < -90 || latitude > 90 || longitude < -180 || longitude > 180 {
			return fmt.Errorf("coordinates out of valid range")
		}
	}

	return nil
}
