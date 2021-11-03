package config

import (
	"crud-product/constant"
	"crud-product/model"
	"encoding/json"
	"io/ioutil"
)

func GetConfig() (*model.Config, error) {
	cfg := &model.Config{}

	jsonFile, err := ioutil.ReadFile(constant.ConfigProjectFilepath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonFile, &cfg)

	return cfg, nil
}