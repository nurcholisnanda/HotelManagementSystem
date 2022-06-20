package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/nurcholisnanda/hotel-management-system/controllers"
	"github.com/nurcholisnanda/hotel-management-system/docs"
	"github.com/nurcholisnanda/hotel-management-system/repositories"
	"github.com/nurcholisnanda/hotel-management-system/services"
	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

var (
	repository repositories.HotelMgmtRepo      = repositories.NewHotelMgmtRepo()
	service    services.HotelMgmtService       = services.NewHotelMgmtService(repository)
	controller controllers.HotelMgmtController = controllers.NewHotelMgmtController(service)
)

const applicationBasePath = "/"

//SetupRouter ... Configure routes
func SetupRouter() *gin.Engine {

	// Swagger 2.0 Meta Information
	docs.SwaggerInfo.Title = "Hotel Management API"
	docs.SwaggerInfo.Description = "API for managing your hotel"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "hotel-management-system-bbx.et.r.appspot.com"
	docs.SwaggerInfo.BasePath = applicationBasePath
	docs.SwaggerInfo.Schemes = []string{"https"}

	server := gin.Default()

	SetHotelManagementRoutes(server)

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return server
}
