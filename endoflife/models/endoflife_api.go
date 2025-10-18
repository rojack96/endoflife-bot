package models

type Uri struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
}

type Identifier struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type UnknownProperties struct {
	AdditionalProperties *string `json:"additionalProperties"`
}

type ProductVersion struct {
	Name string `json:"name"`
	Date string `json:"date,omitempty"`
	Link string `json:"link,omitempty"`
}

type ProductRelease struct {
	Name             string          `json:"name"`
	Codename         *string         `json:"codename,omitempty"`
	Label            string          `json:"label"`
	ReleaseDate      string          `json:"releaseDate"`
	IsLts            bool            `json:"isLts"`
	LtsFrom          *string         `json:"ltsFrom,omitempty"`
	IsEoas           bool            `json:"isEoas"`
	EoasFrom         *string         `json:"eoaFrom,omitempty"`
	IsEol            bool            `json:"isEol"`
	EolFrom          *string         `json:"eolFrom,omitempty"`
	IsDiscontinued   bool            `json:"isDiscontinued"`
	DiscontinuedFrom *string         `json:"discontinuedFrom,omitempty"`
	IsEoes           bool            `json:"isEoes"`
	EoesFrom         *string         `json:"eoesFrom,omitempty"`
	IsMaintained     bool            `json:"isMaintained"`
	Latest           *ProductVersion `json:"latest,omitempty"`
	Custom           *any            `json:"custom,omitempty"`
}

type ProductSummary struct {
	Name     string   `json:"name"`
	Label    string   `json:"label"`
	Aliases  []string `json:"aliases,omitempty"`
	Category string   `json:"category"`
	Tags     []string `json:"tags"`
	URI      string   `json:"uri"`
}

type ProductDetails struct {
	Name           string       `json:"name"`
	Label          string       `json:"label"`
	Aliases        []string     `json:"aliases,omitempty"`
	Category       string       `json:"category"`
	Tags           []string     `json:"tags"`
	VersionCommand *string      `json:"versionCommand,omitempty"`
	Identifiers    []Identifier `json:"identifiers,omitempty"`
	Labels         *struct {
		Eoas         *string `json:"eos,omitempty"`
		Discontinued *string `json:"discontinued,omitempty"`
		Eol          string  `json:"eol"`
		Eoes         *string `json:"eoes,omitempty"`
	} `json:"labels,omitempty"`
	Links *struct {
		Icon          *string `json:"icon,omitempty"`
		Html          string  `json:"html"`
		ReleasePolicy *string `json:"releasePolicy,omitempty"`
	} `json:"links,omitempty"`
	Releases []ProductRelease `json:"releases,omitempty"`
}

// Responses

type UriListResponse struct {
	SchemaVersion string `json:"schema_version"`
	Total         int32  `json:"total"`
	Result        []Uri  `json:"result"`
}

type ProductListResponse struct {
	SchemaVersion string           `json:"schema_version"`
	Total         int32            `json:"total"`
	Result        []ProductSummary `json:"result"`
}

type FullProductListResponse struct {
	SchemaVersion string           `json:"schema_version"`
	Total         int32            `json:"total"`
	Result        []ProductDetails `json:"result"`
}

type ProductReleaseResponse struct {
	SchemaVersion string           `json:"schema_version"`
	Result        []ProductRelease `json:"result"`
}

type ProductResponse struct {
	SchemaVersion string         `json:"schema_version"`
	LastModified  string         `json:"last_modified"`
	Result        ProductDetails `json:"result"`
}

type IdentifierListResponse struct {
	SchemaVersion string `json:"schema_version"`
	Total         int32  `json:"total"`
	Result        []struct {
		Identifier string `json:"identifier"`
		Product    struct {
			Name string `json:"name"`
			Uri  string `json:"uri"`
		} `json:"product"`
	} `json:"result"`
}
