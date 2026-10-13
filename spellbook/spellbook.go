package spellbook

import (
	"dnd-api-spellbook/dndapi"
	"fmt"
	"strings"

	"github.com/mitchellh/go-wordwrap"
)

// max output width for lines displayed in spellbook
const spellOutputWidth = 80

// spellbook holds and can output information about various spells.
// essentially acts as a cache of spells searched for.
// TODO consider some way to timeout or invalidate the cache, to ensure all queries are up-to-date.
// in reality, this is unnecessary because this API won't change, but it would be good practice to learn
type Spellbook struct {
	Spells    map[string]dndapi.Spell // list of known spells (the cache)
	apiClient *dndapi.DndApiClient    // client used to get spells we don't know yet

}

func NewSpellbook(apiClient *dndapi.DndApiClient) *Spellbook {
	return &Spellbook{
		Spells:    make(map[string]dndapi.Spell),
		apiClient: apiClient,
	}
}

// retrieves a spell from this spellbook by name.
// returns spell (or nil if not found) and boolean indicating if it was found.
// If the spell doesn't exist in local storage, will perform REST query to get it from the API
func (s *Spellbook) GetSpell(name string) (dndapi.Spell, bool) {
	// check the local "book" (cache) of spells first
	keyName := spellNameToKeyName(name)
	fmt.Println("Checking our spellbook for \"" + name + "\"...")
	spell, found := s.Spells[keyName]

	// if not found, check the API for it
	if !found {
		fmt.Println("Spell not found in our spellbook, consulting with the API...")
		apiSpell, err := s.apiClient.GetSpell(name)
		if err == nil {
			spell = apiSpell
			found = true
			// store the spell locally for future queries
			s.Spells[keyName] = spell
		} else {
			spell = dndapi.Spell{}
			found = false
		}
	}

	return spell, found
}

// formats a Spell for human-readable output.
// we don't just want raw JSON, some fields (like index and url) are unnecessary.
// not outputting all fields for now, just to keep this brief
func FormatSpellDetails(spell dndapi.Spell) string {
	// descriptions can be long and have multiple lines,
	// so first perform wordwrap to prevent long lines
	descriptionString := wordwrap.WrapString(spell.Description, spellOutputWidth)
	// noticed that newlines inserted via the api response have two spaces after, which
	// makes the formatting a bit weird. fix that first.
	descriptionString = strings.ReplaceAll(descriptionString, "\n  ", "\n")
	// replace newlines with indentations for formatting
	descriptionString = strings.ReplaceAll(descriptionString, "\n", "\n      ")

	spellDetails := spell.Name + "\n" +
		"   Level: " + fmt.Sprintf("%d", spell.Level) + "\n" +
		"   School: " + spell.School.Name + "\n" +
		"   Casting Time: " + spell.CastingTime + "\n" +
		"   Range: " + spell.Range + "\n" +
		"   Components: " + strings.Join(spell.Components, ", ") + "\n" +
		"   Duration: " + spell.Duration + "\n" +
		"   Description: " + descriptionString

	return spellDetails
}

// convert spell name to a consistent format for keys in the spellbook's map
func spellNameToKeyName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// display the names of only the spells currently written (cached) in this spellbook
func (s *Spellbook) FormatMySpells() string {
	mySpells := ""
	for _, spell := range s.Spells {
		mySpells += "   " + spell.Name + "\n"
	}
	return mySpells
}

// display the names of all possible spells from the DND API
func (s *Spellbook) FormatAllSpells() string {
	allSpells, err := s.apiClient.GetAllSpells()
	if err != nil {
		fmt.Println("Error retrieving all spells: ", err)
		return ""
	}
	allSpellsStr := ""
	for _, spell := range allSpells {
		allSpellsStr += "   " + spell.Name + "\n"
	}
	return allSpellsStr
}
