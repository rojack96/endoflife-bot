package endoflife

type EndOfLifeService interface {
	GetAllProducts() ([]string, error)
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
