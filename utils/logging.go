package utils

import (
	"io"
	"log"
	"os"
)

func LoggingSettings(logFile string) {
	//ログファイルの読み込み(読み書き・新規作成・追記可能、パーミッション0666)
	logfile, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln(err)
	}

	multiLogFile := io.MultiWriter(os.Stdout, logfile)
	//フォーマットを指定
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	//ログの出力先を設定
	log.SetOutput(multiLogFile)
}
