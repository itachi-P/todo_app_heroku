# todo_app_heroku
Golang for Heroku deployment of the first CRUD system

Heroku無償プランが2022/11/28で終了した為、デプロイ後一瞬だけ課金しPostgresの挙動を一通り確認<br>
その後 `pg_restore` コマンドにより一旦HerokuのPostgresからバックアップを取り、Supabaseに移行<br>
Heroku上でのToDoアプリ公開（DynoとPostgresへのPlanMiniの課金）は数時間で停止<br>

---

以降のアプリ公開はどこにするか検討ちう
