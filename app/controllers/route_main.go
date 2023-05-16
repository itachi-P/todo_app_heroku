package controllers

import (
	"golang_udemy1/40_todo_app/app/models"
	"log"
	"net/http"
)

// 固定の2つの引数を取るハンドラー関数
// http.HandleFuncにて指定するURLで開かれるページの設定
func top(w http.ResponseWriter, r *http.Request) {
	//Cookieを取得し、ログイン状態では表示されないようにする
	_, err := session(w, r)
	if err != nil {
		//様々なハンドラー関数で使うHTMLページ生成の共通関数でトップページを表示
		generateHTML(w, "Hello, Gopher!", "layout", "public_navbar", "top")
	} else {
		//既にログインしている場合はindex.htmlに飛ばす
		http.Redirect(w, r, "/todos", http.StatusFound)
	}
}

// ログイン成功後のホーム画面用ハンドラー
func index(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		//エラーがあればログイン状態ではないのでトップページに飛ばす
		http.Redirect(w, r, "/", http.StatusFound)
	} else { //セッションが存在すれば、indexページを表示する

		//セッションを使ってユーザー情報を取得
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		//このユーザーが作成したToDo一覧を取得
		todos, _ := user.GetTodosByUser()
		user.Todos = todos

		//ログイン状態であればtop,signup,login画面とそのリンクは表示しない
		//第2引数をnilではなくuserを渡し、画面にユーザーの持つ情報を表示する
		//テンプレートHTML側では {{.Name}}でmodels.User.Nameにアクセスできる
		generateHTML(w, user, "layout", "private_navbar", "index")
	}
}

// 新規ToDo作成画面
func todoNew(w http.ResponseWriter, r *http.Request) {
	_, err := session(w, r) //sessinの値は使わないのでエラーハンドリングのみ
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		generateHTML(w, nil, "layout", "private_navbar", "todo_new")
	}
}

// 新規ToDo保存（画面表示は必ずしも不要）
func todoSave(w http.ResponseWriter, r *http.Request) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		//POSTメソッドで送信されたフォームの値を取得
		err = r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		content := r.PostFormValue("content")
		if err := user.CreateTodo(content); err != nil {
			log.Println(err)
		}
		//ToDoの新規登録が完了したらToDo一覧画面に戻る
		http.Redirect(w, r, "/todos", http.StatusFound)
	}
}

// ToDoの編集(対象IDを引数として持たせる)
// server.go内のparseURL関数の、同じ引数を取る内部関数として実行される
func todoEdit(w http.ResponseWriter, r *http.Request, id int) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		//ユーザーの確認（ただし今回ユーザーは使わないので第1引数は破棄）
		_, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		//データベースから引数のIDによりToDoを取得
		td, err := models.GetTodo(id)
		if err != nil {
			//数値型でなければエラーになる
			log.Println(err)
		}
		//受け取った単一のToDoデータを編集画面のテンプレートに渡す
		generateHTML(w, td, "layout", "private_navbar", "todo_edit")
	}
}

// ToDo編集画面からPOSTメソッドで送られたToDoの変更を反映
func todoUpdate(w http.ResponseWriter, r *http.Request, id int) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		//セッションが存在する場合は、フォームの解析と値の取得を行う
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		user, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		content := r.PostFormValue("content")
		td := &models.Todo{ID: id, Content: content, UserID: user.ID}
		if err := td.UpdateTodo(); err != nil {
			log.Println(err)
		}
		//ToDoの更新が終わったらToDo一覧画面(index.html)に戻す
		http.Redirect(w, r, "/todos", http.StatusFound)
	}
}

// ToDoの削除（特定の1件）
func todoDelete(w http.ResponseWriter, r *http.Request, id int) {
	sess, err := session(w, r)
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		//セッションが存在すれば、念の為ユーザーの存在の有無だけ確認
		_, err := sess.GetUserBySession()
		if err != nil {
			log.Println(err)
		}
		td, err := models.GetTodo(id)
		if err != nil {
			log.Println(err)
		}
		if err := td.DeleteTodo(); err != nil {
			log.Println(err)
		}
		//対象ToDoの削除が完了したらリダイレクト
		http.Redirect(w, r, "/todos", http.StatusFound)
	}
}
