## フォルダ構成
### handler
リクエストハンドラ
service.goにインターフェースを定義

//go:generateから始まるコメントは、ソースコードを自動生成できる
//go:generate go run github.com/matryer/moq -out moq_test.go . ListTasksService