package utils

import (
	"errors"
	"log"

	"github.com/spf13/viper"
)

func CheckConfig() {
	var err error
	v := viper.GetString("TipPrivateKey")
	if v == "" {
		err = errors.New("please set TipPrivateKey in config file")
	}
	v = viper.GetString("PrivateKey")
	if v == "" {
		err = errors.New("please set PrivateKey in config file")
	}
	v = viper.GetString("Https")
	if v == "" {
		err = errors.New("please set Https in config file")
	}
	v = viper.GetString("ComputeUnitPrice")
	if v == "" {
		err = errors.New("please set ComputeUnitPrice in config file")
	}
	v = viper.GetString("ComputeUnitLimit")
	if v == "" {
		err = errors.New("please set ComputeUnitLimit in config file")
	}
	if err != nil {
		log.Fatal(err.Error())
	}

}
