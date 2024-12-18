package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	MaxPower = 256
	OneDay   = 24 * time.Hour
)

type Card struct {
	id           int
	word         string
	example      string
	description  string
	power        int
	revisionDate time.Time
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		EnterNewCard()
		return
	}

	if len(args) == 1 {
		if args[0] == "revise" {
			ReviseNextCard()
			return
		}
		fmt.Println("ERROR: invalid argument.")
		Usage()
		os.Exit(1)
	}

	fmt.Println("ERROR: invalid number of arguments. only one argument is allowed.")
	Usage()
	os.Exit(1)
}

func EnterNewCard() {
	var (
		word        string
		example     string
		description string
	)

	for word == "" {
		word = strings.TrimSpace(Input("word: "))
	}
	for example == "" {
		example = strings.TrimSpace(Input("example: "))
	}
	for description == "" {
		description = strings.TrimSpace(Input("description: "))
	}

	if exists, err := CheckIfWordExists(word); err != nil {
		fmt.Printf("ERROR: something went wrong!! err: %+v\n", err)
		os.Exit(1)
	} else if exists {
		fmt.Printf("ERROR: word '%s' already exists.\n", word)
		os.Exit(1)
	}

	card := &Card{
		word:         word,
		example:      example,
		description:  description,
		power:        1,
		revisionDate: time.Now().Add(24 * time.Hour),
	}

	if err := AddCard(card); err != nil {
		fmt.Printf("ERROR: something went wrong!! err: %+v\n", err)
		os.Exit(1)
	}

	fmt.Println("INFO: card added successfully.")
}

func ReviseNextCard() {
	card, err := GetNextCard()
	if err != nil {
		fmt.Printf("ERROR: something went wrong!! err: %+v\n", err)
		os.Exit(1)
	} else if card == nil {
		return // no cards to revise today
	}

	fmt.Printf("%s - %s\n", card.word, card.example)

	gotit := strings.TrimSpace(Input("got it [y/n] "))
	for gotit != "y" && gotit != "n" {
		gotit = strings.TrimSpace(Input("got it [y/n] "))
	}

	fmt.Println("==>", card.description)

	switch gotit {
	case "y":
		card.power *= 2
		if card.power >= MaxPower {
			if err := DeleteCard(card.id); err != nil {
				fmt.Printf("ERROR: something went wrong!! err: %+v\n", err)
				os.Exit(1)
			}
			return
		}
		card.revisionDate = card.revisionDate.Add(time.Duration(card.power) * OneDay)
		UpdateCard(card)
	case "n":
		card.revisionDate = card.revisionDate.Add(OneDay)
		UpdateCard(card)
	}
}

func Input(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewScanner(os.Stdin)
	reader.Scan()
	return reader.Text()
}

func Usage() {

}
