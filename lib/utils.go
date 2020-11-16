package lib

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func ParseFromYamlFile(path string, ptr interface{}) error {
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Panicf("error reading file. err: [%v]\n", err)
		return err
	}

	err = yaml.Unmarshal(fileData, ptr)

	if err != nil {
		logger.Panicf("fail to unmarshal. Error: [%v]\n", err)
		logger.Panicf("data = [%v]\n", fileData)
	}

	return err
}

// https://gobyexample.com/collection-functions
func Map(vs []string, f func(string) string) []string {
	vsm := make([]string, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}
