package main

import (
	"encoding/json"
	"os"
	_ "log"
	"flag"
	"fmt"
	"path"
)
var lg language

func main() {
	var lang string
	var outputFolder string
	var itemsFolder string
	var resourceFolder string
	var staticFile string
	var medals bool

	flag.StringVar(&lang, "l", "english", "Language")
	flag.BoolVar(&medals, "m", false, "Tournament medals")
	flag.StringVar(&outputFolder, "o", "", "Output folder")
	flag.StringVar(&itemsFolder, "i", "", "Items folder")
	flag.StringVar(&resourceFolder, "r", "", "Resource folder")
	flag.StringVar(&staticFile, "s", "", "Static file")
	flag.Parse()

	if itemsFolder == "" {
		fmt.Println("No items folder provided. Use the flag -i")
		os.Exit(1)
	}
	if resourceFolder == "" {
		fmt.Println("No resource folder provided. Use the flag -r")
		os.Exit(1)
	}
	if outputFolder == "" {
		fmt.Println("No output folder provided. Use the flag -o")
		os.Exit(1)
	}

	lg = language{}
	lg.init(path.Join(resourceFolder, "csgo_" + lang + ".txt"))

	ig := itemsGame{}
	itemsGameDatas, _ := os.ReadFile(path.Join(itemsFolder, "items_game.txt"))
	ig.init(itemsGameDatas)

	j, _ := json.MarshalIndent(&ig, "", "\t")
	os.WriteFile(path.Join(outputFolder, "items_" + lang + ".json"), j, 0666)


	/*heroes := ig.getItemsPerHero()

	for npc, hero := range heroes {
		j, _ := json.MarshalIndent(hero, "", "\t")
		os.WriteFile(path.Join(outputFolder, npc + ".json"), j, 0666)
	}*/
}

func getStringToken(token string) string {
	s, exist := lg.getToken(token)

	if (exist) {
		return s
	} else {
		return token
	}
}
