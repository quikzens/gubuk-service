package house

import (
	"context"
	"errors"
	"gubuk-service/media"
	"gubuk-service/util"
	"strconv"
	"strings"

	db "gubuk-service/db"
	sqlc "gubuk-service/db/sqlc"

	sq "github.com/Masterminds/squirrel"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateHouse(c *gin.Context) {
	payload, _ := c.Get("user")
	userPayload, _ := payload.(*util.UserPayload)
	userID := userPayload.UserID

	ownerID, err := uuid.Parse(userID)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	var req HouseCreateRequest
	err = c.Bind(&req)
	if err != nil {
		util.SendBadRequest(c, err)
		return
	}

	featuredImage, err := c.FormFile("featured_image")
	if err != nil {
		util.SendBadRequest(c, err)
		return
	}

	err = media.ValidateImage(featuredImage)
	if err != nil {
		util.SendBadRequest(c, err)
		return
	}

	newFeaturedImageURL, err := media.UploadMedia("house", featuredImage)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	newHouse, err := db.Queries.CreateHouse(context.TODO(), sqlc.CreateHouseParams{
		ID:            uuid.New(),
		OwnerID:       ownerID,
		Title:         req.Title,
		FeaturedImage: newFeaturedImageURL,
		Bedrooms:      int32(req.Bedrooms),
		Bathrooms:     int32(req.Bathrooms),
		TypeRent:      req.TypeRent,
		Price:         req.Price,
		ProvinceID:    int32(req.ProvinceID),
		CityID:        int32(req.CityID),
		Amenities:     req.Amenities,
		Description:   req.Description,
		Area:          int32(req.Area),
	})
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	util.SendSuccess(c, newHouse)
}

func UpdateHouse(c *gin.Context) {
	payload, _ := c.Get("user")
	userPayload, _ := payload.(*util.UserPayload)
	userID := userPayload.UserID

	houseID := c.Param("id")
	id, err := uuid.Parse(houseID)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	var req HouseUpdateRequest
	err = c.Bind(&req)
	if err != nil {
		util.SendBadRequest(c, err)
		return
	}

	updatedHouse, err := db.Queries.GetHouseById(context.TODO(), id)
	if err != nil {
		util.SendBadRequest(c, errors.New("house with the provided id is not exist"))
		return
	}

	if userID != updatedHouse.OwnerID.String() {
		util.SendBadRequest(c, errors.New("you are not own this house, you could not delete it"))
		return
	}

	updateHouseParams := sqlc.UpdateHouseParams{
		ID:            id,
		Title:         req.Title,
		FeaturedImage: updatedHouse.FeaturedImage,
		Bedrooms:      int32(req.Bedrooms),
		Bathrooms:     int32(req.Bathrooms),
		TypeRent:      req.TypeRent,
		Price:         req.Price,
		ProvinceID:    int32(req.ProvinceID),
		CityID:        int32(req.CityID),
		Description:   req.Description,
		Amenities:     req.Amenities,
		Area:          int32(req.Area),
	}

	featuredImage, err := c.FormFile("featured_image")
	if err == nil {
		newFeaturedImage, err := media.UpdateMedia("house", updatedHouse.FeaturedImage, featuredImage)
		if err != nil {
			util.SendServerError(c, err)
			return
		}

		updateHouseParams.FeaturedImage = newFeaturedImage
	}

	updatedHouseData, err := db.Queries.UpdateHouse(context.TODO(), updateHouseParams)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	util.SendSuccess(c, updatedHouseData)
}

func DeleteHouse(c *gin.Context) {
	payload, _ := c.Get("user")
	userPayload, _ := payload.(*util.UserPayload)
	userID := userPayload.UserID

	houseID := c.Param("id")
	id, err := uuid.Parse(houseID)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	deletedHouse, err := db.Queries.GetHouseById(context.TODO(), id)
	if err != nil {
		util.SendBadRequest(c, errors.New("house with the provided id is not exist"))
		return
	}

	if userID != deletedHouse.OwnerID.String() {
		util.SendBadRequest(c, errors.New("you are not own this house, you could not delete it"))
		return
	}

	err = media.DestroyMedia(deletedHouse.FeaturedImage)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	err = db.Queries.DeleteHouse(context.TODO(), id)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	util.SendSuccess(c, nil)
}

func GetHouseList(c *gin.Context) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	listHouseQueryBuilder := psql.Select("id", "title", "featured_image", "bedrooms", "bathrooms", "type_rent", "price", "province_id", "city_id", "description", "amenities", "area", "created_at", "updated_at").From("homes")

	typeRentFilter := c.Query("type_rent")
	if typeRentFilter != "" {
		listHouseQueryBuilder = listHouseQueryBuilder.Where(sq.Eq{
			"type_rent": typeRentFilter,
		})
	}

	minBedroomFilter, _ := strconv.Atoi(c.Query("bedrooms"))
	if minBedroomFilter > 0 {
		listHouseQueryBuilder = listHouseQueryBuilder.Where(sq.GtOrEq{
			"bedrooms": minBedroomFilter,
		})
	}

	minBathroomFilter, _ := strconv.Atoi(c.Query("bathrooms"))
	if minBathroomFilter > 0 {
		listHouseQueryBuilder = listHouseQueryBuilder.Where(sq.GtOrEq{
			"bathrooms": minBathroomFilter,
		})
	}

	maxPriceFilter, _ := strconv.Atoi(c.Query("price"))
	if maxPriceFilter > 0 {
		listHouseQueryBuilder = listHouseQueryBuilder.Where(sq.LtOrEq{
			"price": maxPriceFilter,
		})
	}

	provinceFilter, _ := strconv.Atoi(c.Query("province_id"))
	if provinceFilter > 0 {
		listHouseQueryBuilder = listHouseQueryBuilder.Where(sq.Eq{
			"province_id": provinceFilter,
		})
	}

	cityFilter, _ := strconv.Atoi(c.Query("city_id"))
	if cityFilter > 0 {
		listHouseQueryBuilder = listHouseQueryBuilder.Where(sq.Eq{
			"city_id": cityFilter,
		})
	}

	amenitiesFilter := c.Query("amenities")
	if amenitiesFilter != "" {
		amenitiesFilters := strings.Split(amenitiesFilter, ",")
		for _, v := range amenitiesFilters {
			listHouseQueryBuilder = listHouseQueryBuilder.Where(sq.ILike{"amenities": "%" + v + "%"})
		}
	}

	limitFilter, _ := strconv.Atoi(c.Query("limit"))
	if limitFilter > 0 {
		listHouseQueryBuilder = listHouseQueryBuilder.Limit(uint64(limitFilter))
	}

	offsetFilter, _ := strconv.Atoi(c.Query("offset"))
	if offsetFilter > 0 {
		listHouseQueryBuilder = listHouseQueryBuilder.Offset(uint64(offsetFilter))
	}

	listHouseQueryBuilder = listHouseQueryBuilder.OrderBy("created_at DESC")
	listHouseQuery, args, err := listHouseQueryBuilder.ToSql()
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	rows, err := db.DB.QueryContext(context.TODO(), listHouseQuery, args...)
	if err != nil {
		util.SendServerError(c, err)
		return
	}
	defer rows.Close()
	houseList := make([]sqlc.ListHouseRow, 0)
	for rows.Next() {
		var i sqlc.ListHouseRow
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.FeaturedImage,
			&i.Bedrooms,
			&i.Bathrooms,
			&i.TypeRent,
			&i.Price,
			&i.ProvinceID,
			&i.CityID,
			&i.Description,
			&i.Amenities,
			&i.Area,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			util.SendServerError(c, err)
			return
		}
		houseList = append(houseList, i)
	}
	if err := rows.Close(); err != nil {
		util.SendServerError(c, err)
		return
	}
	if err := rows.Err(); err != nil {
		util.SendServerError(c, err)
		return
	}

	util.SendSuccess(c, houseList)
}

func GetMyHouseList(c *gin.Context) {
	payload, _ := c.Get("user")
	userPayload, _ := payload.(*util.UserPayload)
	userID := userPayload.UserID

	ownerID, err := uuid.Parse(userID)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	myHouseList, err := db.Queries.ListMyHouse(context.TODO(), ownerID)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	util.SendSuccess(c, myHouseList)
}

func GetHouseDetail(c *gin.Context) {
	houseID := c.Param("id")
	id, err := uuid.Parse(houseID)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	house, err := db.Queries.GetHouseById(context.TODO(), id)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	util.SendSuccess(c, house)
}

func GetHouseCount(c *gin.Context) {
	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	countHouseQueryBuilder := psql.Select("count(*)").From("homes")

	typeRentFilter := c.Query("type_rent")
	if typeRentFilter != "" {
		countHouseQueryBuilder = countHouseQueryBuilder.Where(sq.Eq{
			"type_rent": typeRentFilter,
		})
	}

	minBedroomFilter, _ := strconv.Atoi(c.Query("bedrooms"))
	if minBedroomFilter > 0 {
		countHouseQueryBuilder = countHouseQueryBuilder.Where(sq.GtOrEq{
			"bedrooms": minBedroomFilter,
		})
	}

	minBathroomFilter, _ := strconv.Atoi(c.Query("bathrooms"))
	if minBathroomFilter > 0 {
		countHouseQueryBuilder = countHouseQueryBuilder.Where(sq.GtOrEq{
			"bathrooms": minBathroomFilter,
		})
	}

	maxPriceFilter, _ := strconv.Atoi(c.Query("price"))
	if maxPriceFilter > 0 {
		countHouseQueryBuilder = countHouseQueryBuilder.Where(sq.LtOrEq{
			"price": maxPriceFilter,
		})
	}

	provinceFilter, _ := strconv.Atoi(c.Query("province_id"))
	if provinceFilter > 0 {
		countHouseQueryBuilder = countHouseQueryBuilder.Where(sq.Eq{
			"province_id": provinceFilter,
		})
	}

	cityFilter, _ := strconv.Atoi(c.Query("city_id"))
	if cityFilter > 0 {
		countHouseQueryBuilder = countHouseQueryBuilder.Where(sq.Eq{
			"city_id": cityFilter,
		})
	}

	amenitiesFilter := c.Query("amenities")
	if amenitiesFilter != "" {
		amenitiesFilters := strings.Split(amenitiesFilter, ",")
		for _, v := range amenitiesFilters {
			countHouseQueryBuilder = countHouseQueryBuilder.Where(sq.ILike{"amenities": "%" + v + "%"})
		}
	}

	countHouseQuery, args, err := countHouseQueryBuilder.ToSql()
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	row := db.DB.QueryRowContext(context.TODO(), countHouseQuery, args...)
	var houseCount int64
	err = row.Scan(&houseCount)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	util.SendSuccess(c, houseCount)
}
