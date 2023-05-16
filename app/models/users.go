package models

import (
	"log"
	"time"
)

// usersテーブル用構造体
type User struct {
	ID        int
	UUID      string
	Name      string
	Email     string
	PassWord  string
	CreatedAt time.Time
	Todos     []Todo //追加でユーザー情報にToDo一覧を持たせる

}

// セッションテーブル用構造体
type Session struct {
	ID        int
	UUID      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

// User構造体のポインタをレシーバーとして受け取り、errorを返すメソッド
func (u *User) CreateUser() (err error) {
	//idは自動で一意な値としてインクリメントされるので指定しない
	cmd := `INSERT INTO users (
		uuid,
		name,
		email,
		password,
		created_at) VALUES ($1, $2, $3, $4, $5)`

	_, err = Db.Exec(cmd,
		createUUID(),
		u.Name,
		u.Email,
		Encrypt(u.PassWord),
		time.Now())

	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// ユーザー情報を1件取得する（メソッドではなく関数として定義）
func GetUser(id int) (user User, err error) {
	user = User{}
	cmd := `SELECT id, uuid, name, email, password, created_at
	FROM users WHERE id = $1`
	err = Db.QueryRow(cmd, id).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PassWord,
		&user.CreatedAt,
	)
	return user, err
}

// ユーザー情報の更新:Userのポインタをレシーバに持つメソッドとして定義
func (u *User) UpdateUser() (err error) {
	cmd := `UPDATE users SET name = $1, email = $2 WHERE id = $3`
	_, err = Db.Exec(cmd, u.Name, u.Email, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// ユーザーの削除メソッド
func (u *User) DeleteUser() (err error) {
	cmd := `DELETE FROM users WHERE id = $1`
	_, err = Db.Exec(cmd, u.ID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// ログイン処理用メソッド
// 入力されたEmailアドレスからユーザー情報を取得する
func GetUserByEmail(email string) (user User, err error) {
	user = User{}
	cmd := `SELECT id, uuid, name, email, password, created_at
	FROM users WHERE email = $1`
	err = Db.QueryRow(cmd, email).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.PassWord,
		&user.CreatedAt)

	return user, err
}

// セッションを作成/取得するメソッド
func (u *User) CreateSession() (session Session, err error) {
	session = Session{}
	cmd1 := `INSERT INTO sessions (
		uuid, 
		email, 
		user_id, 
		created_at)
	VALUES ($1, $2, $3, $4)`

	_, err = Db.Exec(cmd1, createUUID(), u.Email, u.ID, time.Now())
	if err != nil {
		log.Println(err)
	}

	cmd2 := `SELECT id, uuid, email, user_id, created_at
	FROM sessions WHERE user_id = $1 AND email = $2`

	err = Db.QueryRow(cmd2, u.ID, u.Email).Scan(
		&session.ID,
		&session.UUID,
		&session.Email,
		&session.UserId,
		&session.CreatedAt)

	return session, err
}

// セッションがデータベースに存在するか(ログイン状態か否か)の確認メソッド
func (sess *Session) CheckSession() (valid bool, err error) {
	cmd := `SELECT id, uuid, email, user_id, created_at
	FROM sessions WHERE uuid = $1`

	err = Db.QueryRow(cmd, sess.UUID).Scan(
		&sess.ID,
		&sess.UUID,
		&sess.Email,
		&sess.UserId,
		&sess.CreatedAt)
	//errが存在する=セッションが存在しない
	if err != nil {
		valid = false
		return
	}
	//errでなく、セッションIDが初期値の0ではない=セッションが存在する
	if sess.ID != 0 {
		valid = true
	}
	return valid, err
}

// ログアウト処理
// セッションから取得したCookieと一致したUUIDをテーブルから削除する
func (sess *Session) DeleteSessionByUUID() (err error) {
	cmd := `DELETE FROM sessions WHERE uuid = $1`
	_, err = Db.Exec(cmd, sess.UUID)
	if err != nil {
		log.Fatalln(err)
	}
	return err
}

// セッションからユーザー情報を取得するメソッド
func (sess *Session) GetUserBySession() (user User, err error) {
	user = User{}
	cmd := `SELECT id, uuid, name, email, created_at FROM users
	WHERE id = $1`
	err = Db.QueryRow(cmd, sess.UserId).Scan(
		&user.ID,
		&user.UUID,
		&user.Name,
		&user.Email,
		&user.CreatedAt)

	return user, err
}
