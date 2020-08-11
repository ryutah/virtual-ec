package admin

import (
	"context"
	"testing"

	"cloud.google.com/go/datastore"
	"github.com/ryutah/virtual-ec/infrastructure/firestore"
	"github.com/ryutah/virtual-ec/usecase/admin"

	. "github.com/smartystreets/goconvey/convey"
)

type productCreateInput struct {
	name  string
	price int
}

func (p productCreateInput) Name() string {
	return p.name
}

func (p productCreateInput) Price() int {
	return p.price
}

type productCreateOutputPort struct {
	success admin.ProductCreateSuccess
	failed  admin.ProductCreateFailed
}

func (p *productCreateOutputPort) Success(s admin.ProductCreateSuccess) {
	p.success = s
}

func (p *productCreateOutputPort) Failed(f admin.ProductCreateFailed) {
	p.failed = f
}

type productFindOutputPort struct {
	success admin.ProductFindSuccess
	failed  admin.ProductFindFailed
}

func (p *productFindOutputPort) Success(s admin.ProductFindSuccess) {
	p.success = s
}

func (p *productFindOutputPort) NotFound(f admin.ProductFindFailed) {
	p.failed = f
}

func (p *productFindOutputPort) Failed(f admin.ProductFindFailed) {
	p.failed = f
}

func TestProduct_CreateAndConfirm(t *testing.T) {
	ctx := context.Background()
	client, err := datastore.NewClient(ctx, "test-project")
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	productRepo := firestore.NewProduct(client)

	Convey("商品の追加と登録データの確認をする", t, func() {
		Convey("新商品を作成する", func() {
			createOutputPort := new(productCreateOutputPort)
			success := admin.NewProductCreate(createOutputPort, productRepo).Create(ctx, productCreateInput{
				name:  "新商品",
				price: 1000,
			})
			Convey("作成が正常に終了する", func() {
				So(success, ShouldBeTrue)
				Convey("生成されたIDを指定してProductを取得する", func() {
					findOutputPort := new(productFindOutputPort)
					success := admin.NewProductFind(findOutputPort, productRepo).Find(ctx, createOutputPort.success.ID)
					Convey("取得が正常に終了する", func() {
						So(success, ShouldBeTrue)
					})
					Convey("新規作成したProductが取得できている", func() {
						So(findOutputPort.success, ShouldResemble, admin.ProductFindSuccess{
							ID:    createOutputPort.success.ID,
							Name:  createOutputPort.success.Name,
							Price: createOutputPort.success.Price,
						})
					})
				})
			})
		})
	})
}
