package services

import (
	"errors"
	"sync"
	"time"

	"github.com/nurcholisnanda/hotel-management-system/models"
	"github.com/nurcholisnanda/hotel-management-system/repositories"
)

type HotelMgmtService interface {
	FindAvailableRooms(checkinDate string, checkoutDate string, roomQty int, roomTypeID int) (availableRooms *models.HotelAvailableRoomsResponse, err error)
	FindPromoRooms(req *models.PromoRoomsRequest) (res *models.PromoRoomsResponse, err error)
}

type hotelMgmtService struct {
	repository repositories.HotelMgmtRepo
}

func NewHotelMgmtService(repository repositories.HotelMgmtRepo) HotelMgmtService {
	return &hotelMgmtService{
		repository: repository,
	}
}

//function to find available rooms
func (service *hotelMgmtService) FindAvailableRooms(checkinDate string, checkoutDate string, roomQty int, roomTypeID int) (availableRooms *models.HotelAvailableRoomsResponse, err error) {
	var wg sync.WaitGroup
	totalPrice := 0

	//specify fields for fetching it to the repo
	fields := "date >= ? AND date < ?"
	ids, err := service.repository.FindRoomIDFromStayRoomByFields(fields, checkinDate, checkoutDate)
	if err != nil {
		ids = append(ids, 0)
	}

	//specify fields for fetching it to the repo
	fields = "hotel_id = ? AND room_type_id = ?"
	rooms, err := service.repository.FindAvailableRoomsByFields(ids, fields, 1, roomTypeID)
	if err != nil {
		return nil, err
	}

	//specify fields for fetching it to the repo
	fields = "hotel_id = ? AND room_type_id = ? AND date >= ? AND date < ?"
	prices, err := service.repository.FindPricesByFields(fields, 1, roomTypeID, checkinDate, checkoutDate)
	if err != nil {
		return nil, err
	}

	if len(rooms) < roomQty {
		return nil, errors.New("Total available rooms exceeded your requirement.")
	}

	//sum the total values from price list and change date format
	wg.Add(len(prices))
	for _, price := range prices {
		go func(price *models.Price) {
			defer wg.Done()
			totalPrice = totalPrice + price.Price
			price.Date = price.Date[0:10]
		}(price)
	}
	wg.Wait()

	//assign price list data to rooms data
	wg.Add(len(rooms))
	for _, room := range rooms {
		go func(room *models.Room) {
			defer wg.Done()
			room.Price = prices
		}(room)
	}
	wg.Wait()

	totalPrice = roomQty * totalPrice

	//assign response model with processed data
	availableRooms = &models.HotelAvailableRoomsResponse{
		RoomQty:        roomQty,
		RoomTypeID:     roomTypeID,
		CheckinDate:    checkinDate,
		CheckoutDate:   checkoutDate,
		TotalPrice:     totalPrice,
		AvailableRooms: rooms,
	}

	return availableRooms, nil
}

//function to find room promotion price
func (service *hotelMgmtService) FindPromoRooms(req *models.PromoRoomsRequest) (res *models.PromoRoomsResponse, err error) {
	//initialize promo rules and promo prices
	var wg sync.WaitGroup
	promoPrice, totalPrice := 0, 0
	allowedForAnyStayDay := false
	allowedForBookingDayPromo := false
	allowedForBookingHourPromo := false

	//fetching required promo and promo rules
	promo, err := service.repository.FindPromoByID(req.PromoID)
	if err != nil {
		return nil, err
	}
	promoStay, err := service.repository.FindStayPromoByID(promo.StayDayPromoID)
	if err != nil {
		allowedForAnyStayDay = true
	}
	promoBook, err := service.repository.FindBookingPromoByID(promo.BookingDayPromoID)
	if err != nil {
		allowedForBookingDayPromo = true
	}

	//validate promo rules data from database
	if promo.BookingHourFirst < 0 || promo.BookingHourLast > 23 ||
		promo.MinimumNights < 0 || promo.MinimumRooms < 0 ||
		(promo.BookingHourFirst == promo.BookingHourLast && promo.BookingHourFirst > 0) ||
		(promo.IsPercentage && promo.Percentage < 0) ||
		(!promo.IsPercentage && promo.Currency < 0) ||
		len(req.AvailableRooms) == 0 {
		return nil, errors.New("sorry, this promo currently unavailable")
	}

	//validate booking day promo rules data
	if !allowedForBookingDayPromo {
		bookingWeekday := time.Now().Local().Weekday()
		switch int(bookingWeekday) {
		case 0:
			allowedForBookingDayPromo = promoBook.IsSunPromo
		case 1:
			allowedForBookingDayPromo = promoBook.IsMonPromo
		case 2:
			allowedForBookingDayPromo = promoBook.IsTuePromo
		case 3:
			allowedForBookingDayPromo = promoBook.IsWedPromo
		case 4:
			allowedForBookingDayPromo = promoBook.IsThuPromo
		case 5:
			allowedForBookingDayPromo = promoBook.IsFriPromo
		case 6:
			allowedForBookingDayPromo = promoBook.IsSatPromo
		}
	}

	//validate booking hour promo rules data
	bookingHour := time.Now().Local().Hour()
	if promo.BookingHourFirst == 0 {
		allowedForBookingHourPromo = true
	} else if promo.BookingHourFirst > promo.BookingHourLast {
		allowedForBookingHourPromo = bookingHour >= promo.BookingHourFirst || bookingHour < promo.BookingHourLast
	} else {
		allowedForBookingHourPromo = bookingHour >= promo.BookingHourFirst && bookingHour < promo.BookingHourLast
	}

	//get promo prices with specific rules
	totalNights := len(req.AvailableRooms[0].Price)
	if totalNights >= promo.MinimumNights && req.RoomQty >= promo.MinimumRooms &&
		allowedForBookingDayPromo && allowedForBookingHourPromo {

		roomPrices := req.AvailableRooms[0].Price
		//get promo prices if no stay day rules
		if allowedForAnyStayDay {
			wg.Add(len(roomPrices))
			for _, price := range roomPrices {
				go func(price *models.Price) {
					defer wg.Done()
					if promo.IsPercentage {
						temp := promo.Percentage * price.Price / 100
						promoPrice = promoPrice + temp
						price.Price = price.Price - temp
						totalPrice = totalPrice + price.Price
					} else {
						promoPrice = promoPrice + promo.Currency
						price.Price = price.Price - promo.Currency
						totalPrice = totalPrice + price.Price
					}
				}(price)
			}
			wg.Wait()
		} else {
			//get promo prices if there are stay day rules
			wg.Add(len(roomPrices))
			for _, price := range roomPrices {
				go func(price *models.Price) {
					defer wg.Done()
					allowedForPromo := false
					day, err := time.Parse("2006-01-02", price.Date)
					if err != nil {
						totalPrice = totalPrice + price.Price
						return
					}
					weekDay := day.Weekday()
					switch int(weekDay) {
					case 0:
						allowedForPromo = promoStay.IsSunPromo
					case 1:
						allowedForPromo = promoStay.IsMonPromo
					case 2:
						allowedForPromo = promoStay.IsTuePromo
					case 3:
						allowedForPromo = promoStay.IsWedPromo
					case 4:
						allowedForPromo = promoStay.IsThuPromo
					case 5:
						allowedForPromo = promoStay.IsFriPromo
					case 6:
						allowedForPromo = promoStay.IsSatPromo
					}
					if allowedForPromo {
						if promo.IsPercentage {
							temp := promo.Percentage * price.Price / 100
							promoPrice = promoPrice + temp
							price.Price = price.Price - temp
							totalPrice = totalPrice + price.Price
						} else {
							promoPrice = promoPrice + promo.Currency
							price.Price = price.Price - promo.Currency
							totalPrice = totalPrice + price.Price
						}
					} else {
						totalPrice = totalPrice + price.Price
					}
				}(price)
			}
			wg.Wait()
		}

		//assign back price list data after promo pricesz
		wg.Add(len(req.AvailableRooms))
		for _, room := range req.AvailableRooms {
			go func(room *models.Room) {
				wg.Done()
				room.Price = roomPrices
			}(room)
		}
		wg.Wait()

		//assign response model with processed data
		res = &models.PromoRoomsResponse{
			PromoID:        req.PromoID,
			RoomQty:        req.RoomQty,
			RoomTypeID:     req.RoomTypeID,
			CheckinDate:    req.CheckinDate,
			CheckoutDate:   req.CheckoutDate,
			PromoPrice:     promoPrice,
			TotalPrice:     totalPrice,
			AvailableRooms: req.AvailableRooms,
		}
	} else {
		return nil, errors.New("sorry, this promo can't be used for your request")
	}
	return res, nil
}
