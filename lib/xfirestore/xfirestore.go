//go:generate interfacer -for cloud.google.com/go/datastore.Client -as xfirestore.DatastoreClient -o client_gen.go
//go:generate interfacer -for cloud.google.com/go/datastore.Iterator -as xfirestore.iterator -o iterator_gen.go
//go:generate interfacer -for cloud.google.com/go/datastore.Cursor -as xfirestore.cursor -o cursor_gen.go
//go:generate interfacer -for cloud.google.com/go/datastore.Transaction -as xfirestore.transaction -o transaction_gen.go

package xfirestore
