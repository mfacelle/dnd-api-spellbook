package main

import (
	"bufio"
	"dnd-api-spellbook/dndapi"
	"dnd-api-spellbook/spellbook"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {

	apiClient := dndapi.NewDndApiClient(http.Client{
		Timeout: 10 * time.Second,
	})
	mySpellbook := spellbook.NewSpellbook(apiClient)

	fmt.Println("")
	fmt.Println("Welcome to Mike's D&D Spellbook!")

	// loop until user types "exit"
	for {
		// TODO kind of a clunky interface, but I guess it works for now
		fmt.Println("")
		fmt.Println("Special Commands:")
		fmt.Println("   mine: Display only the spells currently in my spellbook")
		fmt.Println("   all:  Display all possible spells")
		fmt.Println("   exit: Close the book (exit the program)")
		fmt.Print("Enter a spell name or special command: ")
		name, err := processSpellName(os.Stdin)

		// handle exit command
		if strings.EqualFold(formatInput(name), "exit") {
			fmt.Println("")
			fmt.Println("Exiting...")
			return
		}

		// handle non-exit special commands
		if strings.EqualFold(formatInput(name), "mine") {
			fmt.Println("")
			fmt.Println("My Spells:")
			fmt.Println(mySpellbook.FormatMySpells())
			continue
		}
		if strings.EqualFold(formatInput(name), "all") {
			fmt.Println("")
			fmt.Println("All Spells:")
			fmt.Println(mySpellbook.FormatAllSpells())
			continue
		}

		// normal processing for spell lookup
		if err == nil {
			err = lookupSpell(mySpellbook, name)
		} else {
			fmt.Println("Error: ", err)
		}
	}
}

// process input from user to get a spell name to use for searching
func processSpellName(input io.Reader) (string, error) {
	line, err := bufio.NewReader(input).ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("Error reading spell name: %w", err)
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
		fmt.Println("")
		fmt.Printf("Spell not found: %s\n", name)
		return fmt.Errorf("spell not found: %s", name)
	}

	// print spell details (and add newlines for some formatting)
	fmt.Println("")
	fmt.Println(spellbook.FormatSpellDetails(spell))
	fmt.Println("")

	return nil
}

// format input to be consistent and easily-matchable in processing code
func formatInput(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}
