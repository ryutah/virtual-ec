// Package internal provides primitives to interact the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen DO NOT EDIT.
package internal

// Error defines model for Error.
type Error struct {
	Details *[]string `json:"details,omitempty"`
	Message string    `json:"message"`
}

// Product defines model for Product.
type Product struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Price int64  `json:"price"`
}

// NotFound defines model for NotFound.
type NotFound Error

// ProductCreateSuccess defines model for ProductCreateSuccess.
type ProductCreateSuccess Product

// ProductGetSuccess defines model for ProductGetSuccess.
type ProductGetSuccess Product

// ProductSearchSuccess defines model for ProductSearchSuccess.
type ProductSearchSuccess struct {
	Products []Product `json:"products"`
}

// ServerError defines model for ServerError.
type ServerError Error

// ProductCreate defines model for ProductCreate.
type ProductCreate Product

// ProductSearchParams defines parameters for ProductSearch.
type ProductSearchParams struct {

	// Product名の前方一致検索をする
	Name *string `json:"name,omitempty"`
}

// ProductCreateRequestBody defines body for ProductCreate for application/json ContentType.
type ProductCreateJSONRequestBody ProductCreate
