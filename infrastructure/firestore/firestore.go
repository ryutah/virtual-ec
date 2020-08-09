//go:generate interfacer -for cloud.google.com/go/datastore.Client -as firestore.Client -o client.go

package firestore

var kinds = struct {
	product string
	review  string
}{
	product: "product",
	review:  "review",
}
