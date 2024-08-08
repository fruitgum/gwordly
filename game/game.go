package game

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Matches struct {
	Letter string
	Color  string
}

type Word struct {
	Item  string   `json:"word,omitempty"`
	Score int      `json:"score,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

func GetWords() string {

	word := ""
	for {
		var words []Word

		apiURL := BuildAPIQuery()

		r, _ := http.Get(apiURL)

		defer r.Body.Close()

		jsonData, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(jsonData, &words)
		word = GetWord(words)
		if len(word) > 0 {
			break
		} else {
			time.Sleep(1 * time.Second)
		}
	}

	return word
}

func GetWord(words []Word) string {
	var word string
	var gotWords []string
	for _, w := range words {
		var freqFloat float64
		for _, tag := range w.Tags {

			if len(tag) == 1 {
				if tag != "n" {
					continue
				}
			}

			if len(tag) == 10 {
				freqSplit := strings.Split(tag, ":")
				freqFloat, _ = strconv.ParseFloat(freqSplit[1], 8)
			}

		}

		if freqFloat >= 0.75 {
			check, _ := regexp.MatchString("^[a-zA-Z]{5}$", w.Item)
			if check {
				gotWords = append(gotWords, w.Item)
			}
		}

	}

	if len(gotWords) > 0 {
		minSizeFluct := 0
		maxSizeFluct := len(gotWords) - 1
		randWordIndex := rand.Intn(maxSizeFluct-minSizeFluct+1) + minSizeFluct
		word = gotWords[randWordIndex]
	}

	return word
}

func BuildAPIQuery() string {

	var alphabet []string
	for i := 'a'; i <= 'z'; i++ {
		char := fmt.Sprintf("%c", i)
		alphabet = append(alphabet, char)
	}

	startsFrom := alphabet[rand.Intn(len(alphabet))]

	endWith := alphabet[rand.Intn(len(alphabet))]

	apiURL := "https://api.datamuse.com/words?sp=" + startsFrom + "???" + endWith + "&md=f,p&max=1000"
	return apiURL

}

func Gwordly(word string) {
	f := 0
	fmt.Println("Type `giveup` for finish game and reveal hidden word")
	var gameField []string
	win := false

	fieldStr := "_____"
	for gf := 0; gf < 5; gf++ {
		gameField = append(gameField, fieldStr)
	}

	for g := range gameField {
		fmt.Println(gameField[g])
	}
	reader := bufio.NewReader(os.Stdin)

	for {
		if f == 5 {
			fmt.Println("Fail. Hidden word is " + word)
			Restart()
			break
		}
		fmt.Print("My suggestion is: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "giveup" {
			fmt.Println("Hidden word is " + word)
			Restart()
			break
		}
		if check := CheckInput(input); !check {
			fmt.Print("Invalid input\n")
		} else {
			wordExist, suggestedWord := CheckWordExist(input)
			if !wordExist {
				fmt.Print("Invalid word\n")
			} else {
				gameField[f], win = CheckMatches(word, suggestedWord)
				for g := range gameField {
					fmt.Println(gameField[g])
				}
				if win {
					fmt.Println("You won!")
					Restart()
					break
				}
				f++
			}
		}

	}
}

func CheckMatches(word, suggestedWord string) (string, bool) {

	swChars := strings.Split(suggestedWord, "")
	wChars := strings.Split(word, "")

	matches := make([]string, 5)
	mStruct := make([]Matches, 5)
	greens := 0
	win := false

	var swCharsCount map[string]int

	swCharsCount = make(map[string]int, len(swChars))

	for i := range swChars {
		swCharsCount[swChars[i]] += 1
	}

	for i := range swChars {
		if swChars[i] == wChars[i] {
			matches[i] = fmt.Sprintf(color.GreenString(swChars[i]))
			greens++
		}
	}

	yellowBreak := false

	for i := range swChars {
		for j := range wChars {
			if swChars[i] == wChars[j] {
				if mStruct[i].Color == "Green" {
					yellowBreak = true
				}
				if matches[i] == "" {
					if swCharsCount[swChars[i]] == 1 {
						matches[i] = color.YellowString(swChars[i])
					}
					if swCharsCount[swChars[i]] > 1 {
						matches[i] = color.YellowString(swChars[i])
						yellowBreak = true
					}
				}
			}
		}
		if yellowBreak {
			break
		}
	}

	for i := 0; i < len(swChars); i++ {
		if matches[i] == "" {
			matches[i] = swChars[i]
		}
	}

	if greens == 5 {
		win = true
	}

	return strings.Join(matches, ""), win

}

func CheckWordExist(word string) (bool, string) {

	var words []Word
	APIUrl := "https://api.datamuse.com/sug?s=" + word + "&md=f&max=10"
	r, _ := http.Get(APIUrl)
	defer r.Body.Close()

	jsonData, _ := io.ReadAll(r.Body)

	_ = json.Unmarshal(jsonData, &words)

	if len(words) == 0 {
		return false, ""
	}

	if words[0].Item != word {
		return false, ""
	}

	return true, words[0].Item
}

func CheckInput(word string) bool {

	check, _ := regexp.MatchString("^[a-zA-Z]{5}$", word)

	return check

}
