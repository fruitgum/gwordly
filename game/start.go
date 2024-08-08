package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Start() {
	fmt.Println("Choosing word...")
	word := GetWords()
	Gwordly(word)
}

func Restart() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Again? [y/n]: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	if input == "y" || input == "" {
		Start()
	} else {
		fmt.Println("Bye bye!")
		os.Exit(0)
	}
}
