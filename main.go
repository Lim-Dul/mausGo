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

// Maus blahblah
type Maus struct {
	STR         int
	DEX         int
	WIL         int
	HP          int
	PIPS        int
	Sign        string
	Disposition string
}

// Birthsigns blahblah
type Birthsigns struct {
	Birthsigns []Birthsign `json:"birthsigns"`
}

// Birthsign blahblah
type Birthsign struct {
	Sign        string `json:"sign"`
	Disposition string `json:"disposition"`
}

// Coat blahblah
type Coat struct {
	Colors   []Color   `json:"colors"`
	Patterns []Pattern `json:"patterns"`
}

// Color blahblah
type Color struct {
	Color string `json:"color"`
}

// Pattern blahblah
type Pattern struct {
	Pattern string `json:"pattern"`
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	minSTR := flag.Int("minSTR", 9, "Minimum STR")
	minDEX := flag.Int("minDEX", 9, "Minimum DEX")
	minWIL := flag.Int("minWIL", 9, "Minimum WIL")
	minHP := flag.Int("minHP", 3, "Minimum HP")
	minPIPS := flag.Int("minPIPS", 3, "Minimum Pips")

	flag.Parse()

	myMaus := new(Maus)

	tries := 0
	for {
		tries++
		*myMaus = myMaus.GenStats()
		if myMaus.STR >= *minSTR && myMaus.DEX >= *minDEX && myMaus.WIL >= *minWIL && myMaus.HP >= *minHP && myMaus.PIPS >= *minPIPS {
			fmt.Printf("STR: %d DEX: %d WIL: %d HP: %d Pips: %d Tries: %d\n", myMaus.STR, myMaus.DEX, myMaus.WIL, myMaus.HP, myMaus.PIPS, tries)
			birthsign := RollBirthsign()
			fmt.Printf("Sign: %s | Disposition: %s\n", birthsign.Sign, birthsign.Disposition)
			color, pattern := RollCoat()
			fmt.Printf("Color: %s | Pattern: %s\n", color.Color, pattern.Pattern)
			os.Exit(0)
		}
	}
}

// GenStats blahblah
func (myMaus Maus) GenStats() Maus {
	myMaus.STR = RollStat("3d6kh2")
	myMaus.DEX = RollStat("3d6kh2")
	myMaus.WIL = RollStat("3d6kh2")
	myMaus.HP = RollStat("1d6")
	myMaus.PIPS = RollStat("1d6")
	return myMaus
}

// RollStat blahblah
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

// RollBirthsign blahblah
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

// RollCoat blahblah
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

// ReadJSON blahblah
func ReadJSON(file string) []byte {
	// Open our jsonFile
	jsonFile, err := os.Open(file)
	// if  os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}
