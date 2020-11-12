package lib

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func ParseFromJsonFile(path string, ptr interface{}) error {
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("util: error reading file. err: [%v]\n", err)
		return err
	}

	err = json.Unmarshal(fileData, ptr)

	if err != nil {
		log.Printf("util: fail to unmarshal. Error: [%v]\n", err)
		log.Printf("util: data = [%v]\n", fileData)
	}

	return err
}
