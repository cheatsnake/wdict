package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cheatsnake/wdict/types"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var dictionaryApi string = "https://api.dictionaryapi.dev/api/v2/entries/en/"
var rootCmd = &cobra.Command{
	Use:   "wdict [word]",
	Short: "A simple dictionary CLI app",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			for _, word := range args {
				wdchan := make(chan []types.WordDictionary)
				go getDictionary(word, wdchan)
				showDictionary(<-wdchan)
			}
		} else {
			fmt.Println("Please input one or more words to search.")
		}
	},
}

func getDictionary(word string, ch chan []types.WordDictionary) {
	response, err := http.Get(dictionaryApi + word)
	if err != nil {
		log.Println("Error receiving data...")
	}

	responseByte, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading data...")
	}

	var data []types.WordDictionary

	err = json.Unmarshal(responseByte, &data)
	if err != nil {
		log.Println("Error to unmarshal the response...")
	}

	ch <- data
}

func showDictionary(data []types.WordDictionary) {
	boldWhite := color.New(color.FgHiWhite, color.Bold)
	boldWhite.Println(data[0].Word)

	for _, mean := range data[0].Meanings {
		color.New(color.FgMagenta, color.Italic).Printf("*%s\n", mean.PartOfSpeech)
		for idx, defs := range mean.Definitions {
			color.Cyan("%4d. %4s", idx+1, defs.Definition)
		}
		if len(mean.Synonyms) > 0 {
			color.New(color.FgGreen, color.Bold).Printf("%11s: ", "Synonyms")
			color.Green("%s", strings.Join(mean.Synonyms, ", "))
		}
		if len(mean.Antonyms) > 0 {
			color.New(color.FgRed, color.Bold).Printf("%11s: ", "Antonyms")
			color.Red("%s", strings.Join(mean.Antonyms, ", "))
		}
		fmt.Println("")
	}
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags()
}
