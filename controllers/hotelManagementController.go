package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nurcholisnanda/hotel-management-system/models"
	"github.com/nurcholisnanda/hotel-management-system/services"
)

type HotelMgmtController interface {
	GetAvailableRooms(ctx *gin.Context)
	GetPromoPriceRooms(ctx *gin.Context)
}

type hotelMgmtController struct {
	service services.HotelMgmtService
}

func NewHotelMgmtController(service services.HotelMgmtService) HotelMgmtController {
	return &hotelMgmtController{
		service: service,
	}
}

// GetAvailableRooms godoc
// @Summary Get available rooms
// @Tags Hotel Management
// @Description Get available rooms
// @ID get-available-rooms
// @Accept  json
// @Produce  json
// @Param checkin_date query string true "Checkin date" example("2022-12-31")
// @Param checkout_date query string true "Checkout date" example("2022-12-31")
// @Param room_qty query int true "Room Qty" default(1)
// @Param room_type_id query int true "Room Type ID" default(1)
// @Success 200 {object} models.HotelAvailableRoomsResponse
// @Failure 404 {object} models.ErrResponse
// @Failure 400 {object} models.ErrResponse
// @Router /available-rooms [get]
func (c *hotelMgmtController) GetAvailableRooms(ctx *gin.Context) {
	queryParam := ctx.Request.URL.Query()
	dateForm := "2006-01-02"

	//date validation
	checkinDate := queryParam.Get("checkin_date")
	_, err := time.Parse(dateForm, checkinDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Success: false,
		})
		return
	}

	//date validation
	checkoutDate := queryParam.Get("checkout_date")
	_, err = time.Parse(dateForm, checkoutDate)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Success: false,
		})
		return
	}

	//room quantity validation
	roomQty := queryParam.Get("room_qty")
	qty, err := strconv.ParseUint(roomQty, 10, 32)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Success: false,
		})
		return
	}

	//room type id validation
	roomTypeID := queryParam.Get("room_type_id")
	rTypeId, err := strconv.Atoi(roomTypeID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Success: false,
		})
		return
	}

	//call function to get available rooms
	availableRooms, err := c.service.FindAvailableRooms(checkinDate, checkoutDate, int(qty), rTypeId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, models.ErrResponse{
			Status:  http.StatusNotFound,
			Message: err.Error(),
			Success: false,
		})
		return
	} else {
		ctx.JSON(http.StatusOK, availableRooms)
	}
}

// GetPromoPriceRooms godoc
// @Summary Get rooms with promo prices
// @Tags Hotel Management
// @Description Get rooms with promo prices
// @ID get-promo-rooms
// @Accept  json
// @Produce  json
// @Param body body models.PromoRoomsRequest true "Models of PromoRoomsRequest type"
// @Success 200 {object} models.PromoRoomsResponse
// @Failure 404 {object} models.ErrResponse
// @Failure 400 {object} models.ErrResponse
// @Router /promo-rooms [post]
func (c *hotelMgmtController) GetPromoPriceRooms(ctx *gin.Context) {
	var req models.PromoRoomsRequest

	err := ctx.BindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, models.ErrResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Success: false,
		})
		return
	} else {
		//call function to get promo rooms
		PromoRooms, err := c.service.FindPromoRooms(&req)
		if err != nil {
			ctx.JSON(http.StatusNotFound, models.ErrResponse{
				Status:  http.StatusNotFound,
				Message: err.Error(),
				Success: false,
			})
			return
		} else {
			ctx.JSON(http.StatusOK, PromoRooms)
		}
	}
}
