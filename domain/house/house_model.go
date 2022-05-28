package house

type HouseCreateRequest struct {
	Title       string `form:"title" binding:"required"`
	Bedrooms    int    `form:"bedrooms" binding:"required"`
	Bathrooms   int    `form:"bathrooms" binding:"required"`
	TypeRent    string `form:"type_rent" binding:"required"`
	Price       int64  `form:"price" binding:"required"`
	ProvinceID  int    `form:"province_id" binding:"required"`
	CityID      int    `form:"city_id" binding:"required"`
	Description string `form:"description" binding:"required"`
	Amenities   string `form:"amenities"`
	Area        int    `form:"area" binding:"required"`
}

type HouseUpdateRequest struct {
	Title       string `form:"title" binding:"required"`
	Bedrooms    int    `form:"bedrooms" binding:"required"`
	Bathrooms   int    `form:"bathrooms" binding:"required"`
	TypeRent    string `form:"type_rent" binding:"required"`
	Price       int64  `form:"price" binding:"required"`
	ProvinceID  int    `form:"province_id" binding:"required"`
	CityID      int    `form:"city_id" binding:"required"`
	Description string `form:"description" binding:"required"`
	Amenities   string `form:"amenities"`
	Area        int    `form:"area" binding:"required"`
}
