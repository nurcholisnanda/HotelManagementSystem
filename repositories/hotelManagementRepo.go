package repositories

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/nurcholisnanda/hotel-management-system/configs"
	"github.com/nurcholisnanda/hotel-management-system/models"
)

type HotelMgmtRepo interface {
	FindAvailableRoomsByFields(unavailableRoomId []int, fields string, values ...interface{}) (rooms []*models.Room, err error)
	FindPricesByFields(fields string, values ...interface{}) (prices []*models.Price, err error)
	FindRoomIDFromStayRoomByFields(fields string, values ...interface{}) (ids []int, err error)
	FindPromoByID(id int) (*models.Promo, error)
	FindStayPromoByID(id int) (*models.StayDayPromo, error)
	FindBookingPromoByID(id int) (*models.BookingDayPromo, error)
}

type hotelMgmtRepo struct {
	connection *gorm.DB
}

func NewHotelMgmtRepo() HotelMgmtRepo {
	//connect to mysql database
	db, err := gorm.Open("mysql", configs.DbURL(configs.BuildDBConfig()))
	if err != nil {
		fmt.Println("Status:", err)
		fmt.Println("connection error")
	}
	//auto migrate table by model
	db.AutoMigrate(&models.Hotel{}, &models.Room{}, &models.RoomType{}, &models.Price{}, &models.Order{},
		&models.OrderStatus{}, &models.Stay{}, &models.StayRoom{}, &models.Reservation{},
		&models.Promo{}, &models.BookingDayPromo{}, &models.StayDayPromo{})

	//add foreign key that next can be used for preload gorm func
	db.Model(&models.Room{}).AddForeignKey("hotel_id", "hotels(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Room{}).AddForeignKey("room_type_id", "room_types(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Price{}).AddForeignKey("hotel_id", "hotels(id)", "RESTRICT", "RESTRICT")
	db.Model(&models.Price{}).AddForeignKey("room_type_id", "room_types(id)", "RESTRICT", "RESTRICT")
	return &hotelMgmtRepo{
		connection: db,
	}
}

//fetching available rooms by defined fields and value
func (repo *hotelMgmtRepo) FindAvailableRoomsByFields(unavailableRoomId []int, fields string, values ...interface{}) (rooms []*models.Room, err error) {
	if err = repo.connection.Debug().Not(unavailableRoomId).Where(fields, values...).Find(&rooms).Error; err != nil {
		return nil, err
	}
	return rooms, nil
}

//fetching booked room id by defined fields and value
func (repo *hotelMgmtRepo) FindRoomIDFromStayRoomByFields(fields string, values ...interface{}) (ids []int, err error) {
	var stayRooms []*models.StayRoom
	if err = repo.connection.Debug().Where(fields, values...).Select("DISTINCT room_id").Find(&stayRooms).Error; err != nil {
		return nil, err
	}
	for _, stayRoom := range stayRooms {
		ids = append(ids, stayRoom.RoomID)
	}
	return ids, nil
}

//fetching prices by defined fields and value
func (repo *hotelMgmtRepo) FindPricesByFields(fields string, values ...interface{}) (prices []*models.Price, err error) {
	if err = repo.connection.Debug().Where(fields, values...).Find(&prices).Error; err != nil {
		return nil, err
	}
	return prices, nil
}

//fetching promos by id
func (repo *hotelMgmtRepo) FindPromoByID(id int) (*models.Promo, error) {
	var promo models.Promo
	if err := repo.connection.Debug().Where("id = ?", id).First(&promo).Error; err != nil {
		return nil, err
	}
	return &promo, nil
}

//fetching stay promo by id
func (repo *hotelMgmtRepo) FindStayPromoByID(id int) (*models.StayDayPromo, error) {
	var promo models.StayDayPromo
	if err := repo.connection.Debug().Where("id = ?", id).First(&promo).Error; err != nil {
		return nil, err
	}
	return &promo, nil
}

//fetching booking promo by id
func (repo *hotelMgmtRepo) FindBookingPromoByID(id int) (*models.BookingDayPromo, error) {
	var promo models.BookingDayPromo
	if err := repo.connection.Debug().Where("id = ?", id).First(&promo).Error; err != nil {
		return nil, err
	}
	return &promo, nil
}
