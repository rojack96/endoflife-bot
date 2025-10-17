package endoflife

import (
	"strings"

	"github.com/rojack96/endoflife-bot/endoflife/endpoints"
	"github.com/rojack96/endoflife-bot/endoflife/models"
	httpclient "github.com/rojack96/endoflife-bot/http"
)

type EndOfLifeService struct{}

// NewEndOfLifeService creates a new instance of EndOfLifeService
func NewEndOfLifeService() *EndOfLifeService {
	return &EndOfLifeService{}
}

func (e *EndOfLifeService) GetIndex() (*models.UriListResponse, error) {
	var (
		result models.UriListResponse
		err    error
	)

	url := endpoints.Index

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAllProducts retrieves all products from the End of Life API
func (e *EndOfLifeService) GetAllProducts() (*models.ProductListResponse, error) {
	var (
		result models.ProductListResponse
		err    error
	)

	url := endpoints.Products

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAllProductsFull retrieves all products with full details from the End of Life API
func (e *EndOfLifeService) GetAllProductsFull() (*models.FullProductListResponse, error) {
	var (
		result models.FullProductListResponse
		err    error
	)

	url := endpoints.ProductsFull

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetProduct retrieves a specific product by its name from the End of Life API
func (e *EndOfLifeService) GetProduct(product string) (*models.ProductResponse, error) {
	var (
		result models.ProductResponse
		err    error
	)

	url := strings.Replace(endpoints.ProductsOne, "{product}", product, 1)

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetProductReleases retrieves releases for a specific product from the End of Life API
func (e *EndOfLifeService) GetProductReleases(product, release string) (*models.ProductReleaseResponse, error) {
	var (
		result models.ProductReleaseResponse
		err    error
	)

	url := strings.Replace(endpoints.ProductsRelease, "{product}", product, 1)
	url = strings.Replace(url, "{release}", release, 1)

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetProductReleasesLatest retrieves the latest release for a specific product from the End of Life API
func (e *EndOfLifeService) GetProductReleasesLatest(product string) (*models.ProductReleaseResponse, error) {
	var (
		result models.ProductReleaseResponse
		err    error
	)

	url := strings.Replace(endpoints.ProductsReleasesLatest, "{product}", product, 1)

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCategories retrieves all categories from the End of Life API
func (e *EndOfLifeService) GetCategories() (*models.UriListResponse, error) {
	var (
		result models.UriListResponse
		err    error
	)

	url := endpoints.Categories

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCategoriesProducts retrieves products for a specific category from the End of Life API
func (e *EndOfLifeService) GetCategoriesProducts(category string) (*models.ProductListResponse, error) {
	var (
		result models.ProductListResponse
		err    error
	)

	url := strings.Replace(endpoints.CategoriesProducts, "{category}", category, 1)

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetTags retrieves all tags from the End of Life API
func (e *EndOfLifeService) GetTags() (*models.UriListResponse, error) {
	var (
		result models.UriListResponse
		err    error
	)

	url := endpoints.Tags

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetTagsProducts retrieves products for a specific tag from the End of Life API
func (e *EndOfLifeService) GetTagsProducts(tag string) (*models.UriListResponse, error) {
	var (
		result models.UriListResponse
		err    error
	)

	url := strings.Replace(endpoints.TagsProducts, "{tag}", tag, 1)

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetIdentifiers retrieves all identifier types from the End of Life API
func (e *EndOfLifeService) GetIdentifiers() (*models.IdentifierListResponse, error) {
	var (
		result models.IdentifierListResponse
		err    error
	)

	url := endpoints.Identifiers

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetIdentifiersType retrieves identifiers for a specific type from the End of Life API
func (e *EndOfLifeService) GetIdentifiersType(identifierType string) (*models.IdentifierListResponse, error) {
	var (
		result models.IdentifierListResponse
		err    error
	)

	url := strings.Replace(endpoints.IdentifiersType, "{identifier_type}", identifierType, 1)

	if err = httpclient.HttpRequest("GET", url, nil, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
