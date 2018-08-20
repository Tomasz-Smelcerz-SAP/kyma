package config

import (
	"log"
	"bytes"
	"io/ioutil"
	"os"
)

const fileName = "overrides.yaml"
var filePath = "installer/" + fileName

//LoadComponentsVersionsYaml reads overrides.yaml file if it exists in installer/ dir
func LoadComponentsVersionsYaml() (*bytes.Buffer, error) {

	_, err := os.Stat(filePath)
	if err != nil {

		if os.IsNotExist(err) {
			log.Printf("File %v does not exist", fileName)
			return nil, nil
		}

		log.Println("Error reading file")
		return nil, err
	}

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(contents), nil
}

