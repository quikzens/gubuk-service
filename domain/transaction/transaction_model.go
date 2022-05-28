package transaction

import (
	"time"

	"github.com/google/uuid"
)

type TransactionCreateRequest struct {
	HouseID  string    `form:"house_id" binding:"required"`
	CheckIn  time.Time `form:"check_in" binding:"required"`
	TimeRent int       `form:"time_rent" binding:"required"`
}

type TransactionListRow struct {
	ID                uuid.UUID `json:"id"`
	TenantFullname    string    `json:"tenant_fullname"`
	TenantGender      string    `json:"tenant_gender"`
	TenantPhoneNumber string    `json:"tenant_phone_number"`
	HouseTitle        string    `json:"house_title"`
	HouseProvinceID   int32     `json:"house_province_id"`
	HouseCityID       int32     `json:"house_city_id"`
	HouseAmenities    string    `json:"house_amenities"`
	HouseTypeRent     string    `json:"house_type_rent"`
	PaymentStatus     string    `json:"payment_status"`
	PaymentProof      string    `json:"payment_proof"`
	TotalPayment      int64     `json:"total_payment"`
	CheckIn           time.Time `json:"check_in"`
	CheckOut          time.Time `json:"check_out"`
	TimeRent          string    `json:"time_rent"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
