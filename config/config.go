package config

import (
	"golang_udemy1/40_todo_app/utils"
	"log"

	"gopkg.in/go-ini/ini.v1"
)

type ConfigList struct {
	Port      string
	SQLDriver string
	DbName    string
	LogFile   string
	//Added Bootstrap and JQuery to Views
	Static string
}

// 外部パッケージからも参照できるようにグローバルで変数宣言
var Config ConfigList

// main関数より先にiniファイルを読み込みたいのでinit()で呼び出す
func init() {
	LoadConfig()
	utils.LoggingSettings(Config.LogFile)
}

// iniファイルを読み込む
func LoadConfig() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalln(err)
	}
	Config = ConfigList{
		Port:      cfg.Section("web").Key("port").MustString("8080"),
		SQLDriver: cfg.Section("db").Key("driver").String(),
		DbName:    cfg.Section("db").Key("name").String(),
		LogFile:   cfg.Section("web").Key("logfile").String(),
		Static:    cfg.Section("web").Key("static").String(),
	}
}
