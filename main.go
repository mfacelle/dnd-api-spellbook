package main

import (
	"bufio"
	"dnd-api-spellbook/dndapi"
	"dnd-api-spellbook/spellbook"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// global reader and writer that will be initialized differently depending on the build target.
// allows different io handling for cli build vs wasm (which requires javascript stuff)
var (
	InputReader  io.Reader
	OutputWriter io.Writer
)

func main() {
	// initialize IO based on what was included in build
	InitIO()

	apiClient := dndapi.NewDndApiClient(http.Client{
		Timeout: 10 * time.Second,
	})
	mySpellbook := spellbook.NewSpellbook(apiClient)

	fmt.Fprintln(OutputWriter, "")
	fmt.Fprintln(OutputWriter, "Welcome to Mike's D&D Spellbook!")

	// loop until user types "exit"
	for {
		// TODO kind of a clunky interface, but I guess it works for now
		fmt.Fprintln(OutputWriter, "")
		fmt.Fprintln(OutputWriter, "Special Commands:")
		fmt.Fprintln(OutputWriter, "   mine: Display only the spells currently in my spellbook")
		fmt.Fprintln(OutputWriter, "   all:  Display all possible spells")
		fmt.Fprintln(OutputWriter, "   exit: Close the book (exit the program)")
		fmt.Fprint(OutputWriter, "Enter a spell name or special command: ")
		name, err := processSpellName(InputReader)

		// handle exit command
		if strings.EqualFold(formatInput(name), "exit") {
			fmt.Fprintln(OutputWriter, "")
			fmt.Fprintln(OutputWriter, "Exiting...")
			return
		}

		// handle non-exit special commands
		if strings.EqualFold(formatInput(name), "mine") {
			fmt.Fprintln(OutputWriter, "")
			fmt.Fprintln(OutputWriter, "My Spells:")
			fmt.Fprintln(OutputWriter, mySpellbook.FormatMySpells())
			continue
		}
		if strings.EqualFold(formatInput(name), "all") {
			fmt.Fprintln(OutputWriter, "")
			fmt.Fprintln(OutputWriter, "All Spells:")
			fmt.Fprintln(OutputWriter, mySpellbook.FormatAllSpells())
			continue
		}

		// normal processing for spell lookup
		if err == nil {
			err = lookupSpell(mySpellbook, name)
		} else {
			fmt.Fprintln(OutputWriter, "Error: ", err)
		}
	}
}

// process input from user to get a spell name to use for searching
func processSpellName(input io.Reader) (string, error) {
	line, err := bufio.NewReader(input).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("Error reading spell name %w", err)
	}
	name := strings.TrimSpace(line)

	if name == "" {
		return "", errors.New("Spell name cannot be empty")
	}
	return name, nil
}

// look up a spell in the spellbook and print its details (or an error message)
func lookupSpell(mySpellbook *spellbook.Spellbook, name string) error {
	spell, found := mySpellbook.GetSpell(name)
	if !found {
		// TODO add some kind of logic to find similar names and prompt "did you mean...?"
		fmt.Fprintln(OutputWriter, "")
		fmt.Fprintf(OutputWriter, "Spell not found: %s\n", name)
		return fmt.Errorf("spell not found: %s", name)
	}

	// print spell details (and add newlines for some formatting)
	fmt.Fprintln(OutputWriter, "")
	fmt.Fprintln(OutputWriter, spellbook.FormatSpellDetails(spell))
	fmt.Fprintln(OutputWriter, "")

	return nil
}

// format input to be consistent and easily-matchable in processing code
func formatInput(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}
