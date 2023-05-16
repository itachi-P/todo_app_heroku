package controllers

import (
	"golang_udemy1/40_todo_app/app/models"
	"log"
	"net/http"
)

/* /signup/でアクセスするページのハンドラー関数 */
func signup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		//非ログイン状態の時だけアクセスできるようにする
		_, err := session(w, r)
		if err != nil {
			generateHTML(w, nil, "layout", "public_navbar", "signup")
		} else {
			//セッションが存在する場合はindex.htmlに飛ばす
			http.Redirect(w, r, "/todos", http.StatusFound)
		}
	} else if r.Method == "POST" {
		//入力フォームの値が送信された場合の処理を記述(ユーザー登録)
		err := r.ParseForm() //入力フォームの解析
		if err != nil {
			log.Fatalln(err)
		}
		//入力フォームの値からユーザーの構造体を新規作成
		user := models.User{
			Name:     r.PostFormValue("name"),
			Email:    r.PostFormValue("email"),
			PassWord: r.PostFormValue("password"),
		}
		//データベースに登録(INSERT)
		if err := user.CreateUser(); err != nil {
			log.Println(err)
		}
		//ユーザー登録が成功したらトップページにリダイレクトさせる
		//第3引数にリダイレクト先、第4引数にステータスコードを指定
		//http.Redirect(w, r, "/", 302)
		//ステータスコードを302のように数値リテラルで指定せず、定数の使用が推奨
		http.Redirect(w, r, "/", http.StatusFound)
	}
}

// ログイン用ページのハンドラー
func login(w http.ResponseWriter, r *http.Request) {
	//セッションが存在せず、非ログイン状態の場合のみ表示する
	_, err := session(w, r)
	if err != nil {
		generateHTML(w, nil, "layout", "public_navbar", "login")
	} else {
		//既にセッションが存在する場合はログイン後のindex.htmlに飛ばす
		http.Redirect(w, r, "/todos", http.StatusFound)
	}
}

func authenticate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	user, err := models.GetUserByEmail(r.PostFormValue("email"))
	if err != nil {
		log.Println(err)
		//エラーの場合はログイン失敗でログイン画面にリダイレクト
		http.Redirect(w, r, "/login", http.StatusFound)
	}
	//ユーザーが存在し、フォームで入力されたパスワードがDBと一致する場合はログイン可
	//入力されたパスワードを同じ方式で暗号化した上で比較して一致を確認する
	if user.PassWord == models.Encrypt(r.PostFormValue("password")) {
		session, err := user.CreateSession()
		if err != nil {
			log.Println(err)
		}
		//ここで作成されたセッションを元にしてCookieを作成し、ブラウザにログイン情報を保存
		cookie := http.Cookie{
			Name:     "_cookie",
			Value:    session.UUID, //Cookieの値としてUUIDをそのまま渡す
			HttpOnly: true,
		}
		//Cookieをブラウザに保存
		http.SetCookie(w, &cookie)
		//ログインに成功
		/* 本来ならログインに成功したユーザーのみ見れるページに遷移する段階だが
		今回まだ作成していないのでとりあえずトップページにリダイレクトさせる */
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		//パスワードが一致しない場合はログイン画面に戻す
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("_cookie")
	if err != nil {
		log.Println(err)
	}

	//Cookieが存在しない以外のエラーだった場合は取得したと同じCookieを削除する
	if err != http.ErrNoCookie {
		session := models.Session{UUID: cookie.Value}
		session.DeleteSessionByUUID()
	}
	//セッションのUUIDをデータベースから削除後、ログイン画面にリダイレクトさせる
	http.Redirect(w, r, "/login", http.StatusFound)
}
