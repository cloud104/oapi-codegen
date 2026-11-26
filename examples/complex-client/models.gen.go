package examples

type Beer struct {
	Abv     float32 `json:"abv"`
	Brewery *string `json:"brewery,omitempty"`
	IBU     *int    `json:"ibu,omitempty"`
	ID      int     `json:"id"`
	Name    string  `json:"name"`
	Style   string  `json:"style"`
}

type BeerCreate struct {
	Abv     float32 `json:"abv"`
	Brewery *string `json:"brewery,omitempty"`
	IBU     *int    `json:"ibu,omitempty"`
	Name    string  `json:"name"`
	Style   string  `json:"style"`
}

type BeerPage struct {
	HasNext     *bool  `json:"hasNext,omitempty"`
	HasPrevious *bool  `json:"hasPrevious,omitempty"`
	Items       []Beer `json:"items"`
	Page        int    `json:"page"`
	PageSize    int    `json:"pageSize"`
	TotalItems  int    `json:"totalItems"`
	TotalPages  int    `json:"totalPages"`
}

type BeerReplace struct {
	Abv     float32 `json:"abv"`
	Brewery *string `json:"brewery,omitempty"`
	IBU     *int    `json:"ibu,omitempty"`
	Name    string  `json:"name"`
	Style   string  `json:"style"`
}

type BeerUpdate struct {
	Abv     *float32 `json:"abv,omitempty"`
	Brewery *string  `json:"brewery,omitempty"`
	IBU     *int     `json:"ibu,omitempty"`
	Name    *string  `json:"name,omitempty"`
	Style   *string  `json:"style,omitempty"`
}

type Error struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

type BadRequest = Error

type Conflict = Error

type NotFound = Error

type UnsupportedMediaType = Error

type ListBeersParams struct {
	Page     *int `form:"page,omitempty" json:"page,omitempty"`
	PageSize *int `form:"pageSize,omitempty" json:"pageSize,omitempty"`
}

type CreateBeerTextBody = string

type CreateBeerJSONRequestBody = BeerCreate

type CreateBeerFormdataRequestBody = BeerCreate

type CreateBeerMultipartRequestBody = BeerCreate

type CreateBeerTextRequestBody = CreateBeerTextBody

type UpdateBeerJSONRequestBody = BeerUpdate

type ReplaceBeerJSONRequestBody = BeerReplace
