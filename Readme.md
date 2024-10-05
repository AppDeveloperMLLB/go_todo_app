## フォルダ構成
### handler
リクエストハンドラ
service.goにインターフェースを定義

//go:generateから始まるコメントは、ソースコードを自動生成できる
//go:generate go run github.com/matryer/moq -out moq_test.go . ListTasksService

### JWTを発行する際のKeyについて
opensslコマンドで、秘密鍵と公開鍵を生成し、それを利用してRS256形式で署名を行なっている

まず秘密鍵を生成する
```
openssl genrsa 4096 > secret.pem
```

その後、公開鍵を生成する
```
openssl rsa -pubout < secret.pem > public.pem
```

作成した鍵は、auth/sertに配置する
鍵は下記のように埋め込んでいます
```
//go:embed cert/secret.pem
var rawPrivKey []byte

//go:embed cert/public.pem
var rawPubKey []byte
```