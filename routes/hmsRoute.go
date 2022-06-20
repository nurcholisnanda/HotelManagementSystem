package routes

import "github.com/gin-gonic/gin"

func SetHotelManagementRoutes(route *gin.Engine) {
	grp1 := route.Group(applicationBasePath)
	{
		grp1.GET("available-rooms", func(ctx *gin.Context) {
			controller.GetAvailableRooms(ctx)
		})
		grp1.POST("promo-rooms", func(ctx *gin.Context) {
			controller.GetPromoPriceRooms(ctx)
		})
	}
}
