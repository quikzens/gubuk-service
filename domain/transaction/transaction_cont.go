package transaction

import (
	"context"
	"errors"
	"gubuk-service/media"
	"gubuk-service/util"
	"strconv"
	"strings"
	"time"

	db "gubuk-service/db"
	sqlc "gubuk-service/db/sqlc"

	sq "github.com/Masterminds/squirrel"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateTransaction(c *gin.Context) {
	payload, _ := c.Get("user")
	userPayload, _ := payload.(*util.UserPayload)
	userID := userPayload.UserID

	var req TransactionCreateRequest
	err := c.Bind(&req)
	if err != nil {
		util.SendBadRequest(c, err)
		return
	}

	tenantID, err := uuid.Parse(userID)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	houseID, err := uuid.Parse(req.HouseID)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	house, err := db.Queries.GetHouseById(context.TODO(), houseID)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	var checkOut time.Time
	switch house.TypeRent {
	case "day":
		checkOut = req.CheckIn.AddDate(0, 0, req.TimeRent)
	case "month":
		checkOut = req.CheckIn.AddDate(0, req.TimeRent, 0)
	case "year":
		checkOut = req.CheckIn.AddDate(req.TimeRent, 0, 0)
	}

	newTransaction, err := db.Queries.CreateTransaction(context.TODO(), sqlc.CreateTransactionParams{
		ID:            uuid.New(),
		TenantID:      tenantID,
		OwnerID:       house.OwnerID,
		HouseID:       houseID,
		PaymentStatus: "waiting-payment",
		PaymentProof:  "",
		TotalPayment:  int64(req.TimeRent) * house.Price,
		CheckIn:       req.CheckIn,
		CheckOut:      checkOut,
		TimeRent:      strconv.Itoa(req.TimeRent),
	})
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	util.SendSuccess(c, newTransaction)
}

func ListTransaction(c *gin.Context) {
	payload, _ := c.Get("user")
	userPayload, _ := payload.(*util.UserPayload)
	userID := userPayload.UserID
	userRole := userPayload.UserRole

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	listTransactionQueryBuilder := psql.Select("transactions.id", "tenant.fullname AS tenant_fullname", "tenant.gender AS tenant_gender", "tenant.phone_number AS tenant_phone_number", "house.title AS house_title", "house.province_id AS house_province_id", "house.city_id AS house_city_id", "house.amenities AS house_amenities", "house.type_rent AS house_type_rent", "payment_status", "payment_proof", "total_payment", "check_in", "check_out", "time_rent", "transactions.created_at", "transactions.updated_at").From("transactions").Join("users AS tenant ON tenant.id = transactions.tenant_id").Join("homes AS house ON house.id = transactions.house_id")

	if userRole == "tenant" {
		listTransactionQueryBuilder = listTransactionQueryBuilder.Where(sq.Eq{"tenant_id": userID})
	} else if userRole == "owner" {
		listTransactionQueryBuilder = listTransactionQueryBuilder.Where(sq.Eq{"transactions.owner_id": userID})
	}

	statusFilter := c.Query("status")
	if statusFilter != "" {
		orQuery := sq.Or{}
		for _, v := range strings.Split(statusFilter, ",") {
			orQuery = append(orQuery, sq.Eq{"payment_status": v})
		}
		listTransactionQueryBuilder = listTransactionQueryBuilder.Where(orQuery)
	}

	listTransactionQuery, args, err := listTransactionQueryBuilder.ToSql()
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	rows, err := db.DB.QueryContext(context.TODO(), listTransactionQuery, args...)
	if err != nil {
		util.SendServerError(c, err)
		return
	}
	defer rows.Close()
	transactionList := make([]TransactionListRow, 0)
	for rows.Next() {
		var i TransactionListRow
		if err := rows.Scan(
			&i.ID,
			&i.TenantFullname,
			&i.TenantGender,
			&i.TenantPhoneNumber,
			&i.HouseTitle,
			&i.HouseProvinceID,
			&i.HouseCityID,
			&i.HouseAmenities,
			&i.HouseTypeRent,
			&i.PaymentStatus,
			&i.PaymentProof,
			&i.TotalPayment,
			&i.CheckIn,
			&i.CheckOut,
			&i.TimeRent,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			util.SendServerError(c, err)
			return
		}
		transactionList = append(transactionList, i)
	}
	if err := rows.Close(); err != nil {
		util.SendServerError(c, err)
		return
	}
	if err := rows.Err(); err != nil {
		util.SendServerError(c, err)
		return
	}

	util.SendSuccess(c, transactionList)
}

func PayTransaction(c *gin.Context) {
	transactionID := c.Param("id")
	id, err := uuid.Parse(transactionID)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	paymentProof, err := c.FormFile("payment_proof")
	if err != nil {
		util.SendBadRequest(c, err)
		return
	}

	err = media.ValidateImage(paymentProof)
	if err != nil {
		util.SendBadRequest(c, err)
		return
	}

	newPaymentProof, err := media.UploadMedia("transaction", paymentProof)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	err = db.Queries.UpdateTransactionPaymentProofById(context.TODO(), sqlc.UpdateTransactionPaymentProofByIdParams{
		ID:            id,
		PaymentStatus: "waiting-approve",
		PaymentProof:  newPaymentProof,
		UpdatedAt:     time.Now(),
	})
	if err != nil {
		util.SendServerError(c, err)
	}

	util.SendSuccess(c, gin.H{
		"new_image": newPaymentProof,
	})
}

func UpdateTransactionStatus(c *gin.Context) {
	transactionID := c.Param("id")
	id, err := uuid.Parse(transactionID)
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	status := c.Query("status")
	if status == "" {
		util.SendBadRequest(c, errors.New("status query is required"))
		return
	}

	err = db.Queries.UpdateTransactionStatusById(context.TODO(), sqlc.UpdateTransactionStatusByIdParams{
		ID:            id,
		PaymentStatus: status,
		UpdatedAt:     time.Now(),
	})
	if err != nil {
		util.SendServerError(c, err)
		return
	}

	util.SendSuccess(c, nil)
}
