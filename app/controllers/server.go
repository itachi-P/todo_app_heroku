package controllers

import (
	"fmt"
	"golang_udemy1/40_todo_app/app/models"
	"golang_udemy1/40_todo_app/config"
	"html/template"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

// 様々なハンドラー関数でテンプレートとしてHTML表示される部分を共通化する関数
// 第3引数は可変長引数を取れるように...<型>とする。実際は/views/templates/%s.htmlの%s部分
func generateHTML(w http.ResponseWriter, data interface{}, filenames ...string) {
	var files []string
	//第3引数で渡されたfilenamesの中身を取り出し、ファイルパスを加えfiles[]に格納
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("app/views/templates/%s.html", file))
	}
	//可変長引数で受け取ったファイルパス群からテンプレートのキャッシュを生成
	templates := template.Must(template.ParseFiles(files...))
	//上記で生成されたテンプレートのキャッシュに対し、適用するレイアウトを指定し実行
	//第2引数に実行するテンプレートを明示的に渡す(HTML中の{{define "templatename"}})
	templates.ExecuteTemplate(w, "layout", data)
}

// Cookie(特定のuser.UUID)を取得・DBと照合し、ログイン状態か否かを判定する
func session(w http.ResponseWriter, r *http.Request) (sess models.Session, err error) {
	cookie, err := r.Cookie("_cookie")
	if err == nil { //エラーが「無ければ」で、エラーハンドリングではない
		sess = models.Session{UUID: cookie.Value}
		if ok, _ := sess.CheckSession(); !ok {
			//もしデータベースのUUIDと一致しなければエラーを生成し強制ログアウト
			err = fmt.Errorf("invalid session")
		}
	}
	//ここでerrが返る場合に処理を分けることでアクセス制限を実現する
	return sess, err
}

// URLから正規表現により特定のパターンを検出し、格納する変数をコンパイルしておく
var validPath = regexp.MustCompile("^/todos/(edit|update|delete)/([0-9]+)$")

// URLからID部分を取得し、IDに紐付いた編集ページアクセス用のハンドラー関数を返す関数
func parseURL(fn func(http.ResponseWriter, *http.Request, int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// URLから /todos/edit/id(int値) というパターンをsliceで取得する
		q := validPath.FindStringSubmatch(r.URL.Path)
		if q == nil {
			//正規表現でマッチするURLが見つからずsliceが空になる場合
			http.NotFound(w, r)
			return
		}
		//URLのsliceのIDの部分は文字列型からINT型に変換しておく
		qi, err := strconv.Atoi(q[2])
		if err != nil {
			http.NotFound(w, r)
			return
		}
		//最終的な処理として、この関数じたいの引数である関数に引数を渡す
		fn(w, r, qi) //fnは下のhttp.HandleFunc第2引数内のtodoEdit()に等しい
	}
}

func StartMainServer() error {
	//CSSファイルとJavaScriptファイルを静的(static)ファイルとして読み込む
	files := http.FileServer(http.Dir(config.Config.Static))
	/* URLを/staticの下の階層とするが、実際のディレクトリ構成では
	/viewsと/css,/jsの間に/staticは無いのでプリフィックスを取り除く*/
	http.Handle("/static/", http.StripPrefix("/static/", files))

	//ハンドラー(topやsignup)と、それに対応させるURLの登録

	//以下の3つのページは非ログイン状態でのみアクセスできるようにする
	http.HandleFunc("/", top)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)

	http.HandleFunc("/authenticate", authenticate)

	//ログイン状態のユーザーのみアクセス可能なページ
	http.HandleFunc("/todos", index)
	http.HandleFunc("/logout", logout)
	http.HandleFunc("/todos/new", todoNew)
	http.HandleFunc("/todos/save", todoSave)
	//URLの末尾に"/"が無い場合は完全一致が要求され、付けた場合は後ろに何が付いても通す
	//parseURL(todoEdit)はハンドラー関数をチェインさせて実行している
	http.HandleFunc("/todos/edit/", parseURL(todoEdit))
	http.HandleFunc("/todos/update/", parseURL(todoUpdate))
	http.HandleFunc("/todos/delete/", parseURL(todoDelete))
	//サーバ起動: 第2引数にnilを渡すことで、デフォルトのマルチプレクサを使用
	//登録されていないURLへのアクセスはデフォルトで"404 page not found"を返す

	//Postgres用にポートを取得
	port := os.Getenv("PORT")
	return http.ListenAndServe(":"+port, nil)
}
