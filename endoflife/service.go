package endoflife

import "github.com/rojack96/endoflife-bot/endoflife/dto"

type EndOfLifeService interface {
	GetAllProducts() ([]string, error)
	GetProductLts(product string) (dto.Product, error)
	GetProductDetails(product string) ([]dto.Product, error)
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

func (e *endOfLifeServiceImpl) GetProductDetails(product string) ([]dto.Product, error) {
	var result []dto.Product
	return result, nil
}
