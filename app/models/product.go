package models

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Product struct {
	ID               string `gorm:"size:36;not null;uniqueIndex;primary_key"`
	ParentID         string `gorm:"size:36;index"`
	User             User
	UserID           string `gorm:"size:36;index"`
	ProductImages    []ProductImage
	Categories       []Category      `gorm:"many2many:product_categories;"`
	Sku              string          `gorm:"size:100;index"`
	Name             string          `gorm:"size:255"`
	Slug             string          `gorm:"size:255"`
	Price            decimal.Decimal `gorm:"type:decimal(16,2);"`
	Stock            int
	Weight           decimal.Decimal `gorm:"type:decimal(10,2);"`
	ShortDescription string          `gorm:"type:text"`
	Description      string          `gorm:"type:text"`
	Status           int             `gorm:"default:0"`
	CreatedAt        time.Time       `json:"-"`
	UpdatedAt        time.Time       `json:"-"`
	DeletedAt        gorm.DeletedAt  `json:"-"`
}

func (p *Product) GetProducts(db *gorm.DB, perPage int, page int) (*[]Product, int64, error) {
	var err error
	var products []Product
	var count int64

	err = db.Debug().Model(&Product{}).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage

	err = db.Debug().Preload("ProductImages").Model(&Product{}).Order("created_at desc").Limit(perPage).Offset(offset).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return &products, count, nil
}

func (p *Product) FindBySlug(db *gorm.DB, slug string) (*Product, error) {
	var err error
	var product Product

	err = db.Debug().Preload("ProductImages").Model(&Product{}).Where("slug = ?", slug).First(&product).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (p *Product) FindByID(db *gorm.DB, productID string) (*Product, error) {
	var err error
	var product Product

	err = db.Debug().Preload("ProductImages").Model(&Product{}).Where("id = ?", productID).First(&product).Error
	if err != nil {
		return nil, err
	}

	return &product, nil
}

//func (p *Product) LinkImage(productImage ProductImage) {
//	p.ProductImages = append(p.ProductImages, productImage)
//}
//
//func FillProductsWithImages(db *gorm.DB, products []Product) ([]Product, error) {
//	var productsImages []ProductImage
//	for _, product := range products {
//		productsImages = append(productsImages, ProductImage{
//			Product: product,
//		})
//	}
//	err := db.Debug().Model(&ProductImage{}).Find(&productsImages).Error
//	if err != nil {
//		//TODO log
//		return products, err
//	}
//
//	for _, productImage := range productsImages {
//		for i, product := range products {
//			if productImage.ProductID == product.ID {
//				products[i].LinkImage(productImage)
//			}
//		}
//	}
//
//	return products, nil
//}
