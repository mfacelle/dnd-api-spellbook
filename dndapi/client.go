package dndapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// overall dnd api path
const apiBaseURL = "https://www.dnd5eapi.co"

// dnd api path to get all spells (or by index by continuing path)
const allSpellsUrl = "/api/2024/spells"

// REST client class.
// able to make queries to DND API, unmarshal responses into spell data types,
// and handle any API-specific conversions
type DndApiClient struct {
	httpClient        http.Client    // HTTP client for REST queries
	spellSummaryCache []SpellSummary // cache of spell summaries
	// TODO consider some kind of cache timeout or something to make this more like a live application.
	// in reality, this API will not really change, but updating the cache would be good practice
	cacheValid bool
}

// constructor
func NewDndApiClient(httpClient http.Client) *DndApiClient {
	return &DndApiClient{
		httpClient:        httpClient,
		spellSummaryCache: []SpellSummary{},
		cacheValid:        false,
	}
}

// get a spell from the DND API
func (c *DndApiClient) GetSpell(name string) (Spell, error) {
	// convert basic name to index-like name and perform query.
	// it's likely the converted name will give us exactly the spell we want,
	// i.e. "Charm Person" would conver to charm-person, which is the exact index needed
	spellIndex := nameToIndex(name)
	spell, err := c.getSpellByIndex(spellIndex)

	// if we failed to query by index, check all spells to find the index and query via that
	if err != nil {
		// if cache isn't valid, query the API for all spells and update the cache
		// TODO redundant with GetAllSpells, consider refactoring
		if !c.cacheValid {
			allSpells, err := c.getAllSpellsFromApi()
			if err == nil {
				c.spellSummaryCache = allSpells
				c.cacheValid = true
			}
		}

		// now use the cache to find the spell we want by name
		spellSummary, found := c.findSpellByNameInCache(name)

		if found {
			// use spell's index to query the API
			spell, err = c.getSpellByIndex(spellSummary.Index)
		} else {
			return Spell{}, fmt.Errorf("Unable to find spell %q", name)
		}
	}

	return spell, err
}

// use an index to look up a specific spell
func (c *DndApiClient) getSpellByIndex(index string) (Spell, error) {
	spellUrl := apiBaseURL + allSpellsUrl + "/" + url.PathEscape(index)
	response, err := c.getJsonResponse(spellUrl)
	if err != nil {
		return Spell{}, fmt.Errorf("fetch spell details by index failed: %w", err)
	}

	// unmarshal into a Spell
	var spell Spell
	if err := json.Unmarshal(response, &spell); err != nil {
		return Spell{}, fmt.Errorf("error unmarshaling spell: %w", err)
	}

	return spell, nil
}

// retrieve all spells, either from the cache or the API.
// update the cache if it was retrieved from the API
func (c *DndApiClient) GetAllSpells() ([]SpellSummary, error) {
	if !c.cacheValid {
		allSpells, err := c.getAllSpellsFromApi()
		if err == nil {
			c.spellSummaryCache = allSpells
			c.cacheValid = true
		}
	}
	return c.spellSummaryCache, nil
}

// get a summary of all spells from the DND API
func (c *DndApiClient) getAllSpellsFromApi() ([]SpellSummary, error) {
	spellsUrl := apiBaseURL + allSpellsUrl
	response, err := c.getJsonResponse(spellsUrl)
	if err != nil {
		return nil, fmt.Errorf("fetch all spells failed: %w", err)
	}

	var results SpellSearchResponse
	if err := json.Unmarshal(response, &results); err != nil {
		return nil, fmt.Errorf("error unmarshaling all spells response: %w", err)
	}

	return results.Results, nil
}

// fetch raw JSON from API endpoint
func (c *DndApiClient) getJsonResponse(endpoint string) (json.RawMessage, error) {
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	request.Header.Set("Accept", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("API returned non-OK status %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %w", err)
	}
	return body, nil
}

// find a spell by exact name in the spell summary cache
func (c *DndApiClient) findSpellByNameInCache(name string) (SpellSummary, bool) {
	for _, spell := range c.spellSummaryCache {
		if strings.EqualFold(strings.TrimSpace(spell.Name), strings.TrimSpace(name)) {
			return spell, true
		}
	}
	return SpellSummary{}, false
}

// converts a human-readable spell name into a suitable index to use when looking up the spell via API.
// may not always result in an exact match to the index we need, so additional search options may be used
func nameToIndex(name string) (index string) {
	index = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(name), " ", "-"))
	return index
}
