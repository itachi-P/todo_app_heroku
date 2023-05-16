package models

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"golang_udemy1/40_todo_app/config"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

//データ取得元・格納先モデルの生成

var Db *sql.DB

var err error

/*
const (
	//Userモデル
	//user情報を格納するテーブルの新規作成
	tableNameUser = "users"
	//ToDoモデル
	tableNameTodo = "todos"
	//ログイン状態保持用のセッションを保存するテーブル
	tableNameSession = "sessions"
)
*/

func init() {
	//Heroku対応：Postges用に修正
	url := os.Getenv("DATABASE_URL")
	connection, _ := pq.ParseURL(url)
	connection += "sslmode=require"
	Db, err = sql.Open(config.Config.SQLDriver, connection)
	if err != nil {
		log.Fatalln(err)
	}
	/*
	   Db, err = sql.Open(config.Config.SQLDriver, config.Config.DbName)

	   	if err != nil {
	   		log.Fatalln(err)
	   	}

	   //usersテーブル作成SQLコマンド
	   cmdU := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(

	   	id INTEGER PRIMARY KEY AUTOINCREMENT,
	   	uuid STRING NOT NULL UNIQUE,
	   	name STRING,
	   	email STRING,
	   	password STRING,
	   	created_at DATETIME)`, tableNameUser)

	   Db.Exec(cmdU)

	   //todosテーブル作成SQLコマンド
	   cmdT := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(

	   	id INTEGER PRIMARY KEY AUTOINCREMENT,
	   	content TEXT,
	   	user_id INTEGER,
	   	created_at DATETIME)`, tableNameTodo)

	   Db.Exec(cmdT)

	   //sessionsテーブル作成コマンド
	   cmdS := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(

	   	id INTEGER PRIMARY KEY AUTOINCREMENT,
	   	uuid STRING NOT NULL UNIQUE,
	   	email STRING,
	   	user_id INTEGER,
	   	created_at DATETIME)`, tableNameSession)

	   Db.Exec(cmdS)
	*/
}

// UUIDを生成する関数
func createUUID() (uuidobj uuid.UUID) {
	uuidobj, _ = uuid.NewUUID()
	return uuidobj
}

// パスワードの保存はハッシュ値に暗号化する必要がある
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return cryptext
}
