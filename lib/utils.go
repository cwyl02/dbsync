package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ParseFromJsonFile(path string, obj interface{}) error {
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("util: error reading file. err: [%v]\n", err)
		return err
	}

	err = json.Unmarshal(fileData, &obj)

	if err != nil {
		fmt.Printf("util: fail to unmarshal. Error: [%v]\n", err)
		fmt.Printf("util: data = [%v]\n", fileData)
	}

	return err
}
