package endpoints

// Index
const (
	// API index
	Index = "/"
)

// Products
const (
	// List all products
	Products = "/products"
	// List all products (full)
	ProductsFull = "/products/full"
	// Get a product
	ProductsOne = "/products/{product}"
	// Get a product release cycle
	ProductsRelease = "/products/{product}/releases/{release}"
	// Get a product latest release cycle
	ProductsReleasesLatest = "/products/{product}/releases/latest"
)

// Categories
const (
	// List all categories
	Categories = "/categories"
	// List products in a category
	CategoriesProducts = "/categories/{category}"
)

// Tags
const (
	// List all tags
	Tags = "/tags"
	// List all products with a tag
	TagsProducts = "/tags/{tag}"
)

// Identifiers
const (
	// List all identifier types
	Identifiers = "/identifiers"
	// List all identifiers for a given type
	IdentifiersType = "/identifiers/{identifier_type}"
)
