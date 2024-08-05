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
	Word  string   `json:"word,omitempty"`
	Score int      `json:"score,omitempty"`
	Tags  []string `json:"tags,omitempty"`
}

func GetWords() string {
	logger := config.BotLogger
	var words []Word

	apiURL := BuildAPIQuery()

	r, err := http.Get(apiURL)
	if err != nil {
		logger.Error("Error getting words: ", err)
	}

	defer r.Body.Close()

	jsonData, _ := io.ReadAll(r.Body)
	err = json.Unmarshal(jsonData, &words)
	if err != nil {
		logger.Error("Can't unmarshal JSON: ", err)
	}
	word := GetWord(words)

	if word == "" {
		logger.Debug("No word found")
		time.Sleep(1 * time.Second)
		GetWords()
	}

	return word
}

func GetWord(words []Word) string {
	logger := config.BotLogger
	logger.Info("Getting word")
	var word string

	for _, w := range words {

		var freqFloat float64
		var partOfSpeech string

		for _, tag := range w.Tags {

			if len(tag) == 1 {
				if tag != "n" {
					continue
				} else {
					partOfSpeech = tag
				}
			}

			if len(tag) == 10 {
				freqSplit := strings.Split(tag, ":")
				freqFloat, _ = strconv.ParseFloat(freqSplit[1], 64)
			}
			logger.Debug("Word: %v, freq: %v, part of speech: %v", w.Word, freqFloat, partOfSpeech)
			if freqFloat > 0.5 {
				word = w.Word
			}
		}
	}

	logger.Debug("Word: %v", word)
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

	logger.Debug("Word starts with %v, ends with %v", startsFrom, endWith)

	apiURL := "https://api.datamuse.com/words?sp=" + startsFrom + "???" + endWith + "&md=f,p&max=1000"
	logger.Debug("API URL: %v", apiURL)
	return apiURL

}

func Gwordly(word string) {
	f := 0

	var gameField []string
	//inputError := false

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
			fmt.Print("Fail. Hidden word is " + word)
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
				gameField[f] = CheckMatches(word, suggestedWord)
				for g := range gameField {
					fmt.Println(gameField[g])
				}
				f++
			}
		}

	}
}

func CheckMatches(word, suggestedWord string) string {

	logger := config.BotLogger

	swChars := strings.Split(suggestedWord, "")
	wChars := strings.Split(word, "")

	matches := make([]string, 5)

	for i := 0; i < len(swChars); i++ {
		for j := 0; j < len(wChars); j++ {
			logger.Debug("wc: %v, index %v; swc: %v, index: %v", wChars[j], j, swChars[i], i)
			if wChars[j] == swChars[i] {
				if j != i {
					matches[i] = fmt.Sprintf(color.YellowString(swChars[i]))
				} else {
					matches[i] = fmt.Sprintf(color.GreenString(swChars[i]))
				}
			} else {
				matches[i] = swChars[i]
			}
		}
	}

	return strings.Join(matches, "")

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

	return true, words[0].Word
}

func CheckInput(word string) bool {
	logger := config.BotLogger
	//logger.Debug("Suggested word %v", word)
	check, err := regexp.MatchString("^[a-zA-Z]{5}$", word)
	if err != nil {
		logger.Error("Error checking word: ", err)
	}

	return check

}
