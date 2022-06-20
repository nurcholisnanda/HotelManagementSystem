package models

type Hotel struct {
	ID        int    `gorm:"primary_key" json:"-"`
	HotelName string `gorm:"type:varchar(100)" json:"hotel_name"`
	Address   string `gorm:"type:varchar(500)" json:"address"`
}

type Room struct {
	ID         int      `gorm:"primary_key" json:"room_id"`
	HotelID    int      `gorm:"hotel_id" json:"-"`
	RoomTypeID int      `gorm:"room_type_id" json:"-"`
	RoomNumber int      `gorm:"default:0" json:"room_number"`
	RoomStatus string   `json:"-"`
	Price      []*Price `gorm:"foreignkey:RoomTypeID" json:"price"`
}

type RoomType struct {
	ID   int    `gorm:"primary_key" json:"id"`
	Name string `gorm:"type:varchar(50)" json:"name"`
}

type Price struct {
	Date       string `gorm:"type:date" json:"date"`
	HotelID    int    `gorm:"hotel_id" json:"-"`
	RoomTypeID int    `gorm:"room_type_id" json:"-"`
	Price      int    `gorm:"default:0" json:"price"`
}

type Reservation struct {
	ID              int    `gorm:"primary_key" json:"id"`
	Order           Order  `gorm:"foreignkey:OrderID"`
	OrderID         int    `gorm:"order_id" json:"-"`
	CustomerName    string `gorm:"type:varchar(50)" json:"customerName" binding:"required"`
	BookedRoomCount int    `gorm:"default:0" json:"bookedRoomCount"`
	CheckinDate     string `gorm:"type:date" json:"checkinDate"`
	CheckoutDate    string `gorm:"type:date" json:"checkoutDate"`
	Hotel           Hotel  `gorm:"foreignkey:HotelID"`
	HotelID         int    `gorm:"hotel_id" json:"-"`
}

type OrderStatus struct {
	ID     int    `gorm:"primary_key" json:"id"`
	Status string `gorm:"type:varchar(20)" json:"status"`
}

type Order struct {
	ID            int         `gorm:"primary_key" json:"id"`
	FinalPrice    int         `gorm:"default:0" json:"finalPrice"`
	OrderStatus   OrderStatus `gorm:"foreignkey:OrderStatusID"`
	OrderStatusID int         `gorm:"order_status_id" json:"-"`
}

type Stay struct {
	ID            int         `gorm:"primary_key" json:"id"`
	Reservation   Reservation `gorm:"foreignkey:ReservationID"`
	ReservationID int         `gorm:"reservation_id" json:"-"`
	GuestName     string      `gorm:"type:varchar(50)" json:"guestName" binding:"required"`
	Room          Room        `gorm:"foreignkey:RoomID"`
	RoomID        int         `gorm:"room_id" json:"-"`
}

type StayRoom struct {
	ID     int    `gorm:"primary_key" json:"id"`
	Stay   Stay   `gorm:"foreignkey:StayID"`
	StayID int    `gorm:"stay_id" json:"-"`
	Room   Room   `gorm:"foreignkey:RoomID"`
	RoomID int    `gorm:"room_id" json:"-"`
	Date   string `gorm:"type:date" json:"date"`
}

type Promo struct {
	ID                int  `gorm:"primary_key" json:"id"`
	MinimumNights     int  `gorm:"default:0" json:"minimum_nigths"`
	MinimumRooms      int  `gorm:"default:0" json:"minimum_rooms"`
	StayDayPromoID    int  `gorm:"stay_day_promo_id" json:"-"`
	BookingDayPromoID int  `gorm:"booking_day_promo_id" json:"-"`
	BookingHourFirst  int  `gorm:"default:0" json:"booking_hour_first"`
	BookingHourLast   int  `gorm:"default:0" json:"booking_hour_last"`
	IsPercentage      bool `gorm:"default:false" json:"is_percentage"`
	Percentage        int  `gorm:"default:0" json:"percentage"`
	Currency          int  `gorm:"default:0" json:"currency"`
}

type StayDayPromo struct {
	ID         int  `gorm:"primary_key" json:"id"`
	IsMonPromo bool `gorm:"default:false" json:"is_mon_promo"`
	IsTuePromo bool `gorm:"default:false" json:"is_tue_promo"`
	IsWedPromo bool `gorm:"default:false" json:"is_wed_promo"`
	IsThuPromo bool `gorm:"default:false" json:"is_thu_promo"`
	IsFriPromo bool `gorm:"default:false" json:"is_fri_promo"`
	IsSatPromo bool `gorm:"default:false" json:"is_sat_promo"`
	IsSunPromo bool `gorm:"default:false" json:"is_sun_promo"`
}

type BookingDayPromo struct {
	ID         int  `gorm:"primary_key" json:"id"`
	IsMonPromo bool `gorm:"default:false" json:"is_mon_promo"`
	IsTuePromo bool `gorm:"default:false" json:"is_tue_promo"`
	IsWedPromo bool `gorm:"default:false" json:"is_wed_promo"`
	IsThuPromo bool `gorm:"default:false" json:"is_thu_promo"`
	IsFriPromo bool `gorm:"default:false" json:"is_fri_promo"`
	IsSatPromo bool `gorm:"default:false" json:"is_sat_promo"`
	IsSunPromo bool `gorm:"default:false" json:"is_sun_promo"`
}

type HotelAvailableRoomsResponse struct {
	RoomQty        int     `json:"room_qty"`
	RoomTypeID     int     `json:"room_type_id"`
	CheckinDate    string  `json:"checkin_date"`
	CheckoutDate   string  `json:"checkout_date"`
	TotalPrice     int     `json:"total_price"`
	AvailableRooms []*Room `json:"available_rooms"`
}

type PromoRoomsRequest struct {
	PromoID        int     `json:"promo_id"`
	RoomQty        int     `json:"room_qty"`
	RoomTypeID     int     `json:"room_type_id"`
	CheckinDate    string  `json:"checkin_date"`
	CheckoutDate   string  `json:"checkout_date"`
	TotalPrice     int     `json:"total_price"`
	AvailableRooms []*Room `json:"available_rooms"`
}

type PromoRoomsResponse struct {
	PromoID        int     `json:"promo_id"`
	RoomQty        int     `json:"room_qty"`
	RoomTypeID     int     `json:"room_type_id"`
	CheckinDate    string  `json:"checkin_date"`
	CheckoutDate   string  `json:"checkout_date"`
	PromoPrice     int     `json:"promo_price"`
	TotalPrice     int     `json:"total_price"`
	AvailableRooms []*Room `json:"available_rooms"`
}
