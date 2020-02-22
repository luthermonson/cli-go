package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
)

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		loadConfig(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		loadConfig(fmt.Sprintf("%s/%s", home, ".civo.json"))
	}
}

func loadConfig(filename string) {
	configFile, err := os.Open(filename)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&CurrentConfig)
}

func saveConfig() {
	var filename string

	if cfgFile != "" {
		filename = cfgFile
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		filename = fmt.Sprintf("%s/%s", home, ".civo.json")
	}

	configJSON, _ := json.Marshal(CurrentConfig)
	err := ioutil.WriteFile(filename, configJSON, 0600)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func outputTable(headers []string, data [][]string) {
	if OutputFormat == "custom" {
		for _, items := range data {
			output := OutputFields
			for index, name := range headers {
				if strings.Contains(output, name) {
					output = strings.Replace(output, name, items[index], 1)
				}
			}
			output = strings.Replace(output, "\\t", "\t", -1)
			output = strings.Replace(output, "\\n", "\n", -1)
			fmt.Println(output)
		}
	} else if OutputFormat == "table" || OutputFormat == "human" || OutputFormat == "" {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader(headers)
		table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
		table.SetAutoFormatHeaders(false)
		table.AppendBulk(data)
		table.Render()
	}
}

func outputKeyValue(data map[string]string) {
	if OutputFormat == "custom" {
		output := OutputFields
		for key, value := range data {
			if strings.Contains(output, key) {
				output = strings.Replace(output, key, value, 1)
			}
			output = strings.Replace(output, "\\t", "\t", -1)
			output = strings.Replace(output, "\\n", "\n", -1)
			fmt.Println(output)
		}
	} else if OutputFormat == "table" || OutputFormat == "human" || OutputFormat == "" {
		table := tablewriter.NewWriter(os.Stdout)
		for key, value := range data {
			table.Append([]string{key, value})
		}
		table.Render()
	}
}