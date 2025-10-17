package endoflife

import "github.com/rojack96/endoflife-bot/endoflife/models"

type EndOfLifeRepository interface {
	GetIndex() (*models.UriListResponse, error)
	GetAllProducts() (*models.ProductListResponse, error)
	GetAllProductsFull() (*models.FullProductListResponse, error)
	GetProduct(product string) (*models.ProductResponse, error)
	GetProductReleases(product, release string) (*models.ProductReleaseResponse, error)
	GetProductReleasesLatest(product string) (*models.ProductReleaseResponse, error)
	GetCategories() (*models.UriListResponse, error)
	GetCategoriesProducts(category string) (*models.ProductListResponse, error)
	GetTags() (*models.UriListResponse, error)
	GetTagsProducts(tag string) (*models.UriListResponse, error)
	GetIdentifiers() (*models.IdentifierListResponse, error)
	GetIdentifiersType(identifierType string) (*models.IdentifierListResponse, error)
}
