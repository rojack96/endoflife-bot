package endoflife

import "github.com/rojack96/endoflife-bot/endoflife/dto"

type EndOfLifeService interface {
	GetAllProducts() ([]string, error)
	GetProductLts() (dto.Product, error)
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

func (e *endOfLifeServiceImpl) GetProductLts() (dto.Product, error) {
	var result dto.Product
	return result, nil
}

func (e *endOfLifeServiceImpl) GetProductDetails(product string) ([]dto.Product, error) {
	var result []dto.Product
	return result, nil
}
