package test

import (
	"context"
	"testing"

	"cloud.google.com/go/datastore"
	"github.com/ryutah/virtual-ec/infrastructure/firestore"
	"github.com/ryutah/virtual-ec/usecase"

	. "github.com/smartystreets/goconvey/convey"
)

func TestProduct(t *testing.T) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "test-project")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	productRepo := firestore.NewProduct(client)

	Convey("商品の追加と登録データの確認をする", t, func() {
		Convey("新商品を作成する", func() {
			createResult, err := usecase.NewProductCreator(productRepo).Append(ctx, usecase.ProductAddRequest{
				Name:  "新商品",
				Price: 1000,
			})
			Convey("作成が正常に終了する", func() {
				So(err, ShouldBeNil)
				Convey("生成されたIDを指定してProductを取得する", func() {
					findResult, err := usecase.NewProductFind(productRepo).Find(ctx, createResult.ID)
					Convey("取得が正常に終了する", func() {
						So(err, ShouldBeNil)
					})
					Convey("新規作成したProductが取得できている", func() {
						So(*findResult, ShouldResemble, usecase.ProductFindResponse{
							ID:    createResult.ID,
							Name:  createResult.Name,
							Price: createResult.Price,
						})
					})
				})
			})
		})
	})
}
