package main

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

var AppConfig Conf

type Conf struct {
	DataSourceName  string        // Mysql
	Ip              string        // Ip, на котором прослушиваются порты
	StartPort       uint16        // Начальный порт
	RangePort       uint16        // Количество занимаемых портов + максимальное ограничение на количество потоков
	DelayUpdatePack time.Duration // Задержка обновления Списков (Минуты)
	PathToConfDir   string
	CheckerUrl      string
}

func GetConfig() {
	file, err := os.Open("conf.json")
	defer file.Close()
	if err != nil {
		log.Panicln("Error occurred while reading config")
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&AppConfig)
	if err != nil {
		log.Panicln("Invalid json")
	}
}
