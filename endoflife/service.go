package endoflife

import (
	"github.com/rojack96/endoflife-bot/endoflife/dto"
)

type EndOfLifeService interface {
	GetAllProducts() ([]string, error)
	GetProductLts(product string) (dto.Product, error)
	GetProducts(product string) ([]dto.Product, error)
	GetProductReleases(product, release string) (dto.Product, error)
}
type endOfLifeServiceImpl struct {
	repo EndOfLifeRepository
}

func NewEndOfLifeService(repo EndOfLifeRepository) EndOfLifeService {
	return &endOfLifeServiceImpl{repo}
}

func (e *endOfLifeServiceImpl) GetAllProducts() ([]string, error) {
	var result []string

	allProducts, err := e.repo.GetAllProducts()
	if err != nil {
		return nil, err
	}

	for _, product := range allProducts.Result {
		result = append(result, product.Name)
	}

	return result, err
}

func (e *endOfLifeServiceImpl) GetProductLts(product string) (dto.Product, error) {
	var result dto.Product
	p, err := e.repo.GetProduct(product)
	if err != nil {
		return dto.Product{}, err
	}

	for _, release := range p.Result.Releases {
		if release.IsLts {
			result.Name = product
			result.Release = release.Name
			result.Released = release.ReleaseDate
			result.EndOfActiveSupport = release.EoasFrom
			result.EndOfSecuritySupport = release.EolFrom
			result.Latest.Version = release.Latest.Name
			result.Latest.Date = release.Latest.Date
			result.Latest.Link = release.Latest.Link
			break
		}
	}

	return result, nil
}

func (e *endOfLifeServiceImpl) GetProducts(product string) ([]dto.Product, error) {
	var result []dto.Product
	p, err := e.repo.GetProduct(product)
	if err != nil {
		return result, err
	}

	for _, release := range p.Result.Releases {
		result = append(result, dto.Product{
			Name:                 product,
			Release:              release.Name,
			Released:             release.ReleaseDate,
			EndOfActiveSupport:   release.EoasFrom,
			EndOfSecuritySupport: release.EolFrom,
			Latest: struct {
				Version string `json:"version"`
				Date    string `json:"date"`
				Link    string `json:"link"`
			}{
				Version: release.Latest.Name,
				Date:    release.Latest.Date,
				Link:    release.Latest.Link,
			},
		})
	}

	return result, nil
}

func (e *endOfLifeServiceImpl) GetProductReleases(product, release string) (dto.Product, error) {
	{
		var result dto.Product
		p, err := e.repo.GetProductReleases(product, release)
		if err != nil {
			return result, err
		}

		result = dto.Product{
			Name:     product,
			Release:  p.Result.Name,
			Released: p.Result.ReleaseDate,
			Latest: struct {
				Version string `json:"version"`
				Date    string `json:"date"`
				Link    string `json:"link"`
			}{
				Version: p.Result.Latest.Name,
				Date:    p.Result.Latest.Date,
				Link:    p.Result.Latest.Link,
			},
		}

		if p.Result.EoasFrom != nil {
			result.EndOfActiveSupport = p.Result.EoasFrom
		}

		if p.Result.EolFrom != nil {
			result.EndOfSecuritySupport = p.Result.EolFrom
		}

		return result, nil
	}
}
