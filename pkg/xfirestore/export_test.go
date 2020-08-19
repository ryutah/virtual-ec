package xfirestore

type (
	GenIterator    = iterator
	GenTransaction = transaction
	GenCursor      = cursor
)

var NewIterator = newIterator

var WrapRunInTransactionFunc = wrapRunInTransactionFunc
