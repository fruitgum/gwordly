package main

import (
	"fmt"
	"gwordly/game"
)

func main() {

	appVersion := "0.2b"

	fmt.Println("GWordly by Fruitgum - a console version of Wordly written on Go")
	fmt.Println("https://github.com/fruitgum/GWordly")
	fmt.Println("Version: " + appVersion)
	fmt.Println("2024")
	game.Start()
}
