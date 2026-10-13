package dndapi

// contains data types with JSON annotations for the DND 5e API

// generic API reference for DND 5e API
type APIReference struct {
	Index string `json:"index"`
	Name  string `json:"name"`
	URL   string `json:"url"`
}

// individual spell response when searching for all spells
type SpellSummary struct {
	Index string `json:"index"`
	Name  string `json:"name"`
	Level int    `json:"level"`
	Url   string `json:"url"`
}

// API response when searching for all spells
type SpellSearchResponse struct {
	Results []SpellSummary `json:"results"`
}

// represents all information for a single spell
type Spell struct {
	Index           string             `json:"index"`
	Name            string             `json:"name"`
	Level           int                `json:"level"`
	School          APIReference       `json:"school"`
	Classes         []APIReference     `json:"classes"`
	Subclasses      []APIReference     `json:"subclasses"`
	URL             string             `json:"url"`
	UpdatedAt       string             `json:"updated_at"`
	Description     string             `json:"description"`
	HigherLevel     string             `json:"higher_level,omitempty"`
	Range           string             `json:"range"`
	Components      []string           `json:"components"`
	Material        string             `json:"material,omitempty"`
	Ritual          bool               `json:"ritual"`
	Duration        string             `json:"duration"`
	Concentration   bool               `json:"concentration"`
	CastingTime     string             `json:"casting_time"`
	AttackType      string             `json:"attack_type,omitempty"`
	Damage          *SpellDamage       `json:"damage,omitempty"`
	DC              *SpellDC           `json:"dc,omitempty"`
	AreaOfEffect    *SpellAreaOfEffect `json:"area_of_effect,omitempty"`
	HealAtSlotLevel map[string]string  `json:"heal_at_slot_level,omitempty"`
}

// spell damage type
type SpellDamage struct {
	DamageType             APIReference      `json:"damage_type"`
	DamageAtSlotLevel      map[string]string `json:"damage_at_slot_level,omitempty"`
	DamageAtCharacterLevel map[string]string `json:"damage_at_character_level,omitempty"`
}

// spell DC information
type SpellDC struct {
	DCType    APIReference `json:"dc_type"`
	DCSuccess string       `json:"dc_success"`
	Desc      string       `json:"desc,omitempty"`
}

// spell area of effect information
type SpellAreaOfEffect struct {
	Type string `json:"type"`
	Size int    `json:"size"`
}
