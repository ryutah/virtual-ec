//go:generate interfacer -for cloud.google.com/go/datastore.Client -as xfirestore.DatastoreClient -o client_gen.go
//go:generate interfacer -for cloud.google.com/go/datastore.Iterator -as xfirestore.iterator -o iterator_gen.go
//go:generate interfacer -for cloud.google.com/go/datastore.Cursor -as xfirestore.cursor -o cursor_gen.go
//go:generate interfacer -for cloud.google.com/go/datastore.Transaction -as xfirestore.transaction -o transaction_gen.go

// Package xfirestore cloud.google.com/go/datastore のWrapperコード
//
// デフォルトのSDKでは、インターフェース化されていないオブジェクトを利用するAPIが多く、モック等を使ったテストコードが
// 非常に書きづらい状態なので、各APIをインターフェースでラッピングすることで、datastore SDKの利用元のコードでモックを
// 定義しやすくできるようにした
package xfirestore
