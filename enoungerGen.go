package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
)

//In Dungeons And Dragons, To build an encounter (Monsters for a party to fight against), there are mathmatical rules to keep it balanced.
//This Program takes the amount of monsters, and the amount of XP that your players can handle based on a chart from the book as determined by the GM,  as arguements
//So if you were building a one monster encounter, this enough information
//However, fighting more than a single monster, increases the probability of player loss
//If you were to try and build an encounter with multiple monsters, you would have to do the calculation after consturcting the encounter.
//This results in a sort of backwards construction of the encounter.
//Instead this Program takes, the amount of monsters, and your parties limitations Players(XP)./monjs 5 500
//This will precalculate A number of encounters using the statistics in the monsters manual which has been converted to JSON via a bash Script
type Monsters struct {
	//Slice [Array] of monsters
	Monsters []Monster `json:"monsters"`
}

type Monster struct {
	//We only need the name and XP of the monsters, however other critical stats are included for future expansion
	Name        string `json:"name"`
	Id          int    `json:"id"`
	Age         string `json:"HITPOINTS"`
	AC          int    `json:"AC"`
	CR          int    `json:"CR"`
	XP          string `json:"XP"`
	Description string `json:"Descriptions"`
}

func randomencounters(i int) []int {
	//Converts String to Number, arguement two is Player Number
	argTwo, _ := strconv.Atoi(os.Args[2])
	//Creates Slice One larger than the number of Monsters
	s := make([]int, argTwo+1)
	//Generates Random numbers, for each monster, based on the amount of monsters in dummy.json
	for n := range s {
		s[n] = rand.Intn(i)
	}
	return s
}

func getXPandName(is []int, m Monsters) []string {
	// This gets the names and XP of the random monsters
	total := 0
	argTwo, _ := strconv.Atoi(os.Args[2])

	maxTwo := argTwo*2 + 1
	var intXp int
	namesAndXP := make([]string, maxTwo+1)
	for i := range is {
		namesAndXP[i] = m.Monsters[is[i]].XP
		namesAndXP[i+argTwo] = m.Monsters[is[i]].Name
		intXp, _ = strconv.Atoi(m.Monsters[is[i]].XP)

		total = total + intXp
	}
	namesAndXP[maxTwo] = strconv.Itoa(total)
	return namesAndXP
}
func findWinner(s []string) ([]string, bool) {
	var modify [16]float64
	// Arbitrary Table from the DND DMG
	//Increases as the amount of monsters increase
	modify[0] = 1
	modify[1] = 1
	modify[2] = 1.5
	modify[3] = 2
	modify[4] = 2
	modify[5] = 2
	modify[6] = 2
	modify[7] = 2.5
	modify[8] = 2.5
	modify[9] = 2.5
	modify[10] = 2.5
	modify[11] = 3
	modify[12] = 3
	modify[13] = 3
	modify[14] = 3
	modify[15] = 4

	argTwo, _ := strconv.Atoi(os.Args[2])
	argOne, _ := strconv.ParseFloat(os.Args[1], 64)

	total, _ := strconv.ParseFloat(s[argTwo*2+1], 64)
	//total multiplied by modifier
	adjTotal := total * modify[argTwo]
	//If the adjusted XP is lower than the XP expectations, return true
	if adjTotal < argOne {
		return s, true
	}
	return s, false
}
func (m Monster) N() string {
	return m.Name
}

func main() {

	jsonFile, err := os.Open("dummys.json")

	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var monsters Monsters

	json.Unmarshal(byteValue, &monsters)
	//Generate Encounters forever, that meet the XP criteria
	for {
		s := randomencounters(len(monsters.Monsters))
		final := getXPandName(s, monsters)
		msg, tf := findWinner(final)
		if tf {
			fmt.Println("success: ", msg)
		}
	}

	fmt.Println(monsters.Monsters)

}
