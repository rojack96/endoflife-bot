package endoflife

import (
	"strings"

	"github.com/rojack96/endoflife-bot/endoflife/endpoints"
	"github.com/rojack96/endoflife-bot/endoflife/models"
	httpclient "github.com/rojack96/endoflife-bot/http"
	"go.uber.org/zap"
)

const (
	baseURL           = "https://endoflife.date/api/v1"
	productsPathParam = "{product}"
)

type endOfLifeRepositoryImpl struct {
	log *zap.Logger
}

// NewEndOfLifeRepository creates a new instance of EndOfLifeRepository
func NewEndOfLifeRepository(log *zap.Logger) EndOfLifeRepository {
	return &endOfLifeRepositoryImpl{log: log}
}

func (e *endOfLifeRepositoryImpl) GetIndex() (*models.UriListResponse, error) {
	var (
		result models.UriListResponse
		err    error
	)

	url := baseURL + endpoints.Index

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		e.log.Error("failed request GET index", zap.Error(err))
		return nil, err
	}

	return &result, nil
}

// GetAllProducts retrieves all products from the End of Life API
func (e *endOfLifeRepositoryImpl) GetAllProducts() (*models.ProductListResponse, error) {
	var (
		result models.ProductListResponse
		err    error
	)

	url := baseURL + endpoints.Products

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		e.log.Error("failed request GET all products",
			zap.String("endpoint", endpoints.Products),
			zap.Error(err))
		return nil, err
	}

	return &result, nil
}

// GetAllProductsFull retrieves all products with full details from the End of Life API
func (e *endOfLifeRepositoryImpl) GetAllProductsFull() (*models.FullProductListResponse, error) {
	var (
		result models.FullProductListResponse
		err    error
	)

	url := baseURL + endpoints.ProductsFull

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		e.log.Error("failed request GET all products full",
			zap.String("endpoint", endpoints.ProductsFull),
			zap.Error(err))
		return nil, err
	}

	return &result, nil
}

// GetProduct retrieves a specific product by its name from the End of Life API
func (e *endOfLifeRepositoryImpl) GetProduct(product string) (*models.ProductResponse, error) {
	var (
		result models.ProductResponse
		err    error
	)

	url := baseURL + strings.Replace(endpoints.ProductsOne, productsPathParam, product, 1)

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		e.log.Error("failed request GET product",
			zap.String("endpoint", strings.Replace(endpoints.ProductsOne, productsPathParam, product, 1)),
			zap.Error(err))
		return nil, err
	}

	return &result, nil
}

// GetProductReleases retrieves releases for a specific product from the End of Life API
func (e *endOfLifeRepositoryImpl) GetProductReleases(product, release string) (*models.ProductReleaseResponse, error) {
	var (
		result models.ProductReleaseResponse
		err    error
	)

	url := baseURL + strings.Replace(endpoints.ProductsRelease, productsPathParam, product, 1)
	url = strings.Replace(url, "{release}", release, 1)

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		e.log.Error("failed request GET product releases",
			zap.String("endpoint", strings.Replace(endpoints.ProductsRelease, productsPathParam, product, 1)),
			zap.Error(err))
		return nil, err
	}

	return &result, nil
}

// GetProductReleasesLatest retrieves the latest release for a specific product from the End of Life API
func (e *endOfLifeRepositoryImpl) GetProductReleasesLatest(product string) (*models.ProductReleaseResponse, error) {
	var (
		result models.ProductReleaseResponse
		err    error
	)

	url := baseURL + strings.Replace(endpoints.ProductsReleasesLatest, productsPathParam, product, 1)

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		e.log.Error("failed request GET product releases latest",
			zap.String("endpoint", strings.Replace(endpoints.ProductsReleasesLatest, productsPathParam, product, 1)),
			zap.Error(err))
		return nil, err
	}

	return &result, nil
}

// GetCategories retrieves all categories from the End of Life API
func (e *endOfLifeRepositoryImpl) GetCategories() (*models.UriListResponse, error) {
	var (
		result models.UriListResponse
		err    error
	)

	url := baseURL + endpoints.Categories

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		e.log.Error("failed request GET categories",
			zap.String("endpoint", endpoints.Categories),
			zap.Error(err))
		return nil, err
	}

	return &result, nil
}

// GetCategoriesProducts retrieves products for a specific category from the End of Life API
func (e *endOfLifeRepositoryImpl) GetCategoriesProducts(category string) (*models.ProductListResponse, error) {
	var (
		result models.ProductListResponse
		err    error
	)

	url := baseURL + strings.Replace(endpoints.CategoriesProducts, "{category}", category, 1)

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		e.log.Error("failed request GET category products",
			zap.String("endpoint", strings.Replace(endpoints.CategoriesProducts, "{category}", category, 1)),
			zap.Error(err))
		return nil, err
	}

	return &result, nil
}

// GetTags retrieves all tags from the End of Life API
func (e *endOfLifeRepositoryImpl) GetTags() (*models.UriListResponse, error) {
	var (
		result models.UriListResponse
		err    error
	)

	url := baseURL + endpoints.Tags

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		e.log.Error("failed request GET tags",
			zap.String("endpoint", endpoints.Tags),
			zap.Error(err))
		return nil, err
	}

	return &result, nil
}

// GetTagsProducts retrieves products for a specific tag from the End of Life API
func (e *endOfLifeRepositoryImpl) GetTagsProducts(tag string) (*models.UriListResponse, error) {
	var (
		result models.UriListResponse
		err    error
	)

	url := baseURL + strings.Replace(endpoints.TagsProducts, "{tag}", tag, 1)

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		e.log.Error("failed request GET tag products",
			zap.String("endpoint", strings.Replace(endpoints.TagsProducts, "{tag}", tag, 1)),
			zap.Error(err))
		return nil, err
	}

	return &result, nil
}

// GetIdentifiers retrieves all identifier types from the End of Life API
func (e *endOfLifeRepositoryImpl) GetIdentifiers() (*models.IdentifierListResponse, error) {
	var (
		result models.IdentifierListResponse
		err    error
	)

	url := baseURL + endpoints.Identifiers

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		e.log.Error("failed request GET identifiers",
			zap.String("endpoint", endpoints.Identifiers),
			zap.Error(err))
		return nil, err
	}

	return &result, nil
}

// GetIdentifiersType retrieves identifiers for a specific type from the End of Life API
func (e *endOfLifeRepositoryImpl) GetIdentifiersType(identifierType string) (*models.IdentifierListResponse, error) {
	var (
		result models.IdentifierListResponse
		err    error
	)

	url := baseURL + strings.Replace(endpoints.IdentifiersType, "{identifier_type}", identifierType, 1)

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		e.log.Error("failed request GET identifiers by type",
			zap.String("endpoint", strings.Replace(endpoints.IdentifiersType, "{identifier_type}", identifierType, 1)),
			zap.Error(err))
		return nil, err
	}

	return &result, nil
}
