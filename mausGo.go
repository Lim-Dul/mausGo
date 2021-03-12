package main

import (
	"encoding/json" // for reading JSON
	"flag"          // for parsing CLI flags
	"fmt"           // for basic I/O
	"io/ioutil"     // for disk i/o
	"math/rand"     // for random number generation
	"os"            // for exit codes
	"strconv"       // for converting strings
	"time"          // for timing script execution and creating random seed

	"github.com/justinian/dice" // for convenient dice rolling
)

// Maus defines a complete Mausritter character
type Maus struct {
	STR         int
	DEX         int
	WIL         int
	HP          int
	PIPS        int
	Sign        string
	Disposition string
	Color       string
	Pattern     string
	Detail      string
	Background  string
	Item1       string
	Item2       string
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

// Details is a collection of details from the Mausritter rulebook
type Details struct {
	Details []Detail `json:"details"`
}

// Detail is a concrete detail from the Mausritter rulebook
type Detail struct {
	Detail string `json:"detail"`
}

// Backgrounds is a collection of backgrounds from the Mausritter rulebook (based on HP / Pips)
type Backgrounds struct {
	HP []HP `json:"hp"`
}

// HP is the first layer defining the background from the Mausritter rulebook
type HP struct {
	Pips []Pip `json:"pips"`
}

// Pip is the second layer defining the background (and starting items) from the Mausritter rulebook
type Pip struct {
	Background string `json:"background"`
	Item1      string `json:"item1"`
	Item2      string `json:"item2"`
}

func main() {
	rand.Seed(time.Now().UTC().UnixNano()) // seed the randomizer

	// CLI flags with name, default value and help string
	minSTR := flag.Int("minSTR", 2, "Minimum STR")
	minDEX := flag.Int("minDEX", 2, "Minimum DEX")
	minWIL := flag.Int("minWIL", 2, "Minimum WIL")
	minHP := flag.Int("minHP", 1, "Minimum HP")
	minPIPS := flag.Int("minPIPS", 1, "Minimum Pips")
	// parse the flags
	flag.Parse()

	// check if user input impossible to reach minimum values
	if *minSTR > 12 || *minDEX > 12 || *minWIL > 12 {
		fmt.Println("Min values too high! Base attributes cannot be higher than 12!")
		os.Exit(1)
	}
	if *minHP > 6 || *minPIPS > 6 {
		fmt.Println("Min values too high! HP/Pips cannot be higher than 6!")
		os.Exit(2)
	}

	//initialize myMaus
	myMaus := new(Maus)
	// count how many tries we needed
	tries := 0
	for {
		tries++
		// generate base stats of myMaus without details to increase performance
		*myMaus = myMaus.GenStats()
		// check if myMaus passed the required attribute values
		if myMaus.STR >= *minSTR && myMaus.DEX >= *minDEX && myMaus.WIL >= *minWIL && myMaus.HP >= *minHP && myMaus.PIPS >= *minPIPS {
			// now that we passed the check, we can generate the rest of the details
			*myMaus = myMaus.GenDetails()
			// print out everything and exit
			fmt.Printf("STR: %d DEX: %d WIL: %d HP: %d Pips: %d Tries: %d\n", myMaus.STR, myMaus.DEX, myMaus.WIL, myMaus.HP, myMaus.PIPS, tries)
			fmt.Printf("Sign: %s | Disposition: %s\n", myMaus.Sign, myMaus.Disposition)
			fmt.Printf("Color: %s | Pattern: %s | Detail: %s\n", myMaus.Color, myMaus.Pattern, myMaus.Detail)
			fmt.Printf("Background: %s\n", myMaus.Background)
			fmt.Printf("Items: %s, %s, Torches, Rations, <+Weapon of Choice>\n", myMaus.Item1, myMaus.Item2)
			os.Exit(0)
		}
	}
}

// GenStats generates the base stats of a Mausritter character and returns the struct
func (myMaus Maus) GenStats() Maus {
	myMaus.STR = RollStat("3d6kh2")
	myMaus.DEX = RollStat("3d6kh2")
	myMaus.WIL = RollStat("3d6kh2")
	myMaus.HP = RollStat("1d6")
	myMaus.PIPS = RollStat("1d6")

	return myMaus
}

// RollStat rolls dice according to a pattern (dice lib) and returns just the integer
func RollStat(input string) int {
	// only need result and error
	res, _, err := dice.Roll(input)

	// handle error, otherwise just return the result int, not the whole result struct
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	} else {
		return res.Int()
	}
	return 0
}

// GenDetails generates the details of a Mausritter character and returns the struct
func (myMaus Maus) GenDetails() Maus {
	// roll for Sign, Disposition, Color, Pattern and Detail
	myMaus.Sign, myMaus.Disposition = RollBirthsign()
	myMaus.Color, myMaus.Pattern = RollCoat()
	myMaus.Detail = RollDetail()

	// Read in the appropriate background according to HP / Pips
	myMaus.Background, myMaus.Item1, myMaus.Item2 = GetBackground(myMaus.HP, myMaus.PIPS)

	return myMaus
}

// RollBirthsign rolls a random birthsign combination from the Mausritter rulebook
func RollBirthsign() (string, string) {
	// birthsigns stored in JSON
	rawData := ReadJSON("config/birthsigns.json")
	// initialize Birthsigns struct and unmarshal JSON into it
	var birthsigns Birthsigns
	json.Unmarshal(rawData, &birthsigns)
	// roll for random Birthsign up to len of Birthsigns
	nString := "1d" + strconv.Itoa(len(birthsigns.Birthsigns))
	n := RollStat(nString) - 1
	// sign and disposition are dependent, so return in the apropriate entries
	return birthsigns.Birthsigns[n].Sign, birthsigns.Birthsigns[n].Disposition
}

// RollCoat rolls a random coat combination from the Mausritter rulebook
func RollCoat() (string, string) {
	// coat stored in JSON
	rawData := ReadJSON("config/coat.json")
	// initialize Coat struct and unmarshal JSON into it
	var coat Coat
	json.Unmarshal(rawData, &coat)
	// roll for random Color and Pattern up to len of each
	nString := "1d" + strconv.Itoa(len(coat.Colors))
	mString := "1d" + strconv.Itoa(len(coat.Patterns))
	n := RollStat(nString) - 1
	m := RollStat(mString) - 1
	// color and pattern are independent so read out at separate indexes
	return coat.Colors[n].Color, coat.Patterns[m].Pattern
}

// RollDetail rolls a random detail from the Mausritter rulebook
func RollDetail() string {
	// details stored in JSON
	rawData := ReadJSON("config/details.json")
	// initialize Details struct and unmarshal JSON into it
	var details Details
	json.Unmarshal(rawData, &details)
	// roll for random Detail up to len of Details
	nString := "1d" + strconv.Itoa(len(details.Details))
	n := RollStat(nString) - 1
	// return detail
	return details.Details[n].Detail
}

// GetBackground reads out the background and starting items from the Mausritter rulebook based on HP/Pips rolled
func GetBackground(hp, pips int) (string, string, string) {
	// backgrounds stored in JSON
	rawData := ReadJSON("config/backgrounds.json")
	// initialize Details struct and unmarshal JSON into it
	var backgrounds Backgrounds
	json.Unmarshal(rawData, &backgrounds)
	// read out the entry at matching indexes watching out for off-by-one error
	background := backgrounds.HP[hp-1].Pips[pips-1].Background
	item1 := backgrounds.HP[hp-1].Pips[pips-1].Item1
	item2 := backgrounds.HP[hp-1].Pips[pips-1].Item2
	// return vars
	return background, item1, item2
}

// ReadJSON reads raw data from a JSON file
func ReadJSON(file string) []byte {
	// Open JSON file
	jsonFile, err := os.Open(file)
	// if  os.Open returns an error, handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so it can be processed until the function finishes
	defer jsonFile.Close()

	// read our opened JSON file as a byte array and return it
	byteValue, _ := ioutil.ReadAll(jsonFile)
	return byteValue
}
