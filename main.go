package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/justinian/dice"
)

// Maus defines the base stats of a Mausritter character without cosmetics
type Maus struct {
	STR         int
	DEX         int
	WIL         int
	HP          int
	PIPS        int
	Sign        string
	Disposition string
}

// Birthsigns is a collection of birthsigns from the Mausritter rulebook
type Birthsigns struct {
	Birthsigns []Birthsign `json:"birthsigns"`
}

// Birthsign is a concrete combination of sign and disposition from the Mausritter rulebook
type Birthsign struct {
	Sign        string `json:"sign"`
	Disposition string `json:"disposition"`
}

// Coat is a combination of colors and patterns from the Mausritter rulebook
type Coat struct {
	Colors   []Color   `json:"colors"`
	Patterns []Pattern `json:"patterns"`
}

// Color is a concrete coat color from the Mausritter rulebook
type Color struct {
	Color string `json:"color"`
}

// Pattern is a concrete coat pattern from the Mausritter rulebook
type Pattern struct {
	Pattern string `json:"pattern"`
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	minSTR := flag.Int("minSTR", 9, "Minimum STR")
	minDEX := flag.Int("minDEX", 9, "Minimum DEX")
	minWIL := flag.Int("minWIL", 9, "Minimum WIL")
	minHP := flag.Int("minHP", 4, "Minimum HP")
	minPIPS := flag.Int("minPIPS", 4, "Minimum Pips")

	flag.Parse()

	myMaus := new(Maus)

	tries := 0
	for {
		tries++
		*myMaus = myMaus.GenStats()
		if myMaus.STR >= *minSTR && myMaus.DEX >= *minDEX && myMaus.WIL >= *minWIL && myMaus.HP >= *minHP && myMaus.PIPS >= *minPIPS {
			fmt.Printf("STR: %d DEX: %d WIL: %d HP: %d Pips: %d Tries: %d\n", myMaus.STR, myMaus.DEX, myMaus.WIL, myMaus.HP, myMaus.PIPS, tries)
			fmt.Printf("Sign: %s | Disposition: %s\n", myMaus.Sign, myMaus.Disposition)
			color, pattern := RollCoat()
			fmt.Printf("Color: %s | Pattern: %s\n", color.Color, pattern.Pattern)
			os.Exit(0)
		}
	}
}

// GenStats generates the base stats of a Mausritter character
func (myMaus Maus) GenStats() Maus {
	myMaus.STR = RollStat("3d6kh2")
	myMaus.DEX = RollStat("3d6kh2")
	myMaus.WIL = RollStat("3d6kh2")
	myMaus.HP = RollStat("1d6")
	myMaus.PIPS = RollStat("1d6")
	birthsign := RollBirthsign()
	myMaus.Sign = birthsign.Sign
	myMaus.Disposition = birthsign.Disposition
	return myMaus
}

// RollStat rolls dice according to a pattern (dice lib) and returns just the integer
func RollStat(input string) int {
	res, _, err := dice.Roll(input)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	} else {
		return res.Int()
	}
	return 0
}

// RollBirthsign rolls a random birthsign combination from the Mausritter rulebook
func RollBirthsign() Birthsign {
	rawData := ReadJSON("config/birthsigns.json")
	var birthsigns Birthsigns
	json.Unmarshal(rawData, &birthsigns)
	nString := "1d" + strconv.Itoa(len(birthsigns.Birthsigns))
	n := RollStat(nString) - 1
	var birthsign Birthsign
	birthsign = birthsigns.Birthsigns[n]
	return birthsign
}

// RollCoat rolls a random coat combination from the Mausritter rulebook
func RollCoat() (Color, Pattern) {
	rawData := ReadJSON("config/coat.json")
	var coat Coat
	json.Unmarshal(rawData, &coat)
	nString := "1d" + strconv.Itoa(len(coat.Colors))
	mString := "1d" + strconv.Itoa(len(coat.Patterns))
	n := RollStat(nString) - 1
	m := RollStat(mString) - 1
	var color Color
	var pattern Pattern
	color = coat.Colors[n]
	pattern = coat.Patterns[m]
	return color, pattern
}

// ReadJSON reads raw data from a JSON file
func ReadJSON(file string) []byte {
	// Open our jsonFile
	jsonFile, err := os.Open(file)
	// if  os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened JSON file as a byte array
	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}
