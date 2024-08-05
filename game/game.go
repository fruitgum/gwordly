package game

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"gwordly/config"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Word struct {
	Item  string   `json:"word,omitempty"`
	Score int      `json:"score,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

func GetWords() string {
	logger := config.BotLogger
	logger.Debug("Getting words")
	var words []Word

	apiURL := BuildAPIQuery()

	r, err := http.Get(apiURL)
	if err != nil {
		logger.Fatal("Error getting words: ", err)
	}

	defer r.Body.Close()

	jsonData, _ := io.ReadAll(r.Body)
	err = json.Unmarshal(jsonData, &words)
	if err != nil {
		logger.Fatal("Can't unmarshal JSON: ", err)
	}
	word := GetWord(words)

	if word == "" {
		logger.Debug("Can't find word")
		time.Sleep(1 * time.Second)
		words = []Word{}
		GetWords()
	}

	return word
}

func GetWord(words []Word) string {
	logger := config.BotLogger
	logger.Debug("Selecting word")
	var word string

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

			if freqFloat > 0.5 {
				word = w.Item
			}

			logger.Debug("Word %v, freq: %v", w.Item, freqFloat)

		}
		if wordSplit := strings.Split(w.Item, " "); len(wordSplit) > 1 {
			word = ""
		}
	}

	return word
}

func BuildAPIQuery() string {

	logger := config.BotLogger

	var alphabet []string
	for i := 'a'; i <= 'z'; i++ {
		char := fmt.Sprintf("%c", i)
		alphabet = append(alphabet, char)
	}

	startsFrom := alphabet[rand.Intn(len(alphabet))]

	endWith := alphabet[rand.Intn(len(alphabet))]

	apiURL := "https://api.datamuse.com/words?sp=" + startsFrom + "???" + endWith + "&md=f,p&max=1000"
	logger.Debug("API URL: %v", apiURL)
	return apiURL

}

func Gwordly(word string) {
	f := 0

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

	logger := config.BotLogger
	logger.Info("Word %v", word)

	for {
		if f == 5 {
			fmt.Println("Fail. Hidden word is " + word)
			os.Exit(0)
		}
		fmt.Print("My suggestion is:")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

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
					if win {
						fmt.Println("You won!")
						os.Exit(0)
					}
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
	greens := 1
	win := false

	for i := 0; i < len(swChars); i++ {
		for j := 0; j < len(wChars); j++ {
			if swChars[i] == wChars[j] {
				if j != i {
					matches[i] = fmt.Sprintf(color.YellowString(swChars[i]))
					greens++
				} else {
					matches[i] = fmt.Sprintf(color.GreenString(swChars[i]))
				}
			}
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
	logger := config.BotLogger
	var words []Word
	APIUrl := "https://api.datamuse.com/sug?s=" + word + "&md=f&max=1"
	r, err := http.Get(APIUrl)
	if err != nil {
		logger.Error("Error cheking if word exists: ", err)
	}
	defer r.Body.Close()

	jsonData, _ := io.ReadAll(r.Body)

	err = json.Unmarshal(jsonData, &words)
	if err != nil {
		logger.Error("Check Word: Can't unmarshal JSON: ", err)
	}

	if len(words) == 0 {
		return false, ""
	}

	return true, words[0].Item
}

func CheckInput(word string) bool {
	logger := config.BotLogger
	check, err := regexp.MatchString("^[a-zA-Z]{5}$", word)
	if err != nil {
		logger.Error("Error checking word: ", err)
	}

	return check

}
