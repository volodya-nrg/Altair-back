package main

import (
	"altair/configs"
	"altair/pkg/controller"
	"altair/pkg/logger"
	"altair/pkg/service"
	"altair/server"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
	"os"
)

// sort - https://gobyexample.com/sorting

func init() {
	// gin.SetMode(gin.ReleaseMode)

	logger.Init(os.Stdout, os.Stdout, os.Stderr)

	if err := configs.Load("./config.json"); err != nil {
		logger.Error.Fatalln(err.Error())
	}

	pDbConf := configs.Cfg.DB
	if err := server.InitDB(pDbConf.User, pDbConf.Password, pDbConf.Host, pDbConf.Port, pDbConf.Name); err != nil {
		logger.Error.Fatalln(err.Error())
	}
}
func main() {
	r := setupRouter()
	if err := r.Run("127.0.0.1:8080"); err != nil {
		logger.Error.Fatalln(err.Error())
	}
}
func setupRouter() *gin.Engine {
	configCors := cors.DefaultConfig()
	configCors.AllowOrigins = []string{"http://localhost:8080"}
	configCors.AllowCredentials = true
	configCors.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}

	pEngine := gin.Default()
	pEngine.MaxMultipartMemory = 16 << 20 // 16 MiB. Lower memory limit for multipart forms (default is 32 MiB)

	pEngine.LoadHTMLGlob("./web/templates/*")
	pEngine.StaticFile("/favicon.ico", "./web/assets/img/favicon.ico")
	pEngine.Static("/assets", "./web/assets")
	pEngine.Use(static.Serve("/images", static.LocalFile("./web/images", true)))
	pEngine.Use(cors.New(configCors))

	pEngine.GET("/", func(c *gin.Context) {
		serviceCats := service.NewCatService()
		serviceKindProperties := service.NewKindPropertyService()
		serviceProperties := service.NewPropertyService()

		catsTree, err := serviceCats.GetCatsAsTree()
		if err != nil {
			logger.Warning.Println(err)
			c.JSON(500, err.Error())
			return
		}

		kindProperties, err := serviceKindProperties.GetKindProperties()
		if err != nil {
			logger.Warning.Println(err)
			c.JSON(500, err.Error())
			return
		}

		// выгрузить все property для категорий
		properties, err := serviceProperties.GetProperties(false)
		if err != nil {
			logger.Warning.Println(err)
			c.JSON(500, err.Error())
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"catsTree":       catsTree,
			"kindProperties": kindProperties,
			"properties":     properties,
		})
	})

	v1 := pEngine.Group("/api/v1")

	v1.GET("/cats", controller.GetCats)
	v1.GET("/cats/:catId", controller.GetCatsCatId)
	v1.POST("/cats", controller.PostCats)
	v1.PUT("/cats/:catId", controller.PutCatsCatId)
	v1.DELETE("/cats/:catId", controller.DeleteCatsCatId)

	v1.GET("/users", controller.GetUsers)
	v1.GET("/users/:userId", controller.GetUsersUserId)
	v1.POST("/users", controller.PostUsers)
	v1.PUT("/users/:userId", controller.PutUsersUserId)

	v1.GET("/ads", controller.GetAds)
	v1.GET("/ads/:adId", controller.GetAdsAdId)
	v1.POST("/ads", controller.PostAds)
	v1.PUT("/ads/:adId", controller.PutAdsAdId)
	v1.DELETE("/ads/:adId", controller.DeleteAdsAdId)

	v1.GET("/properties", controller.GetProperties)
	v1.GET("/properties/:propertyId", controller.GetPropertiesPropertyId)
	v1.POST("/properties", controller.PostProperties)
	v1.PUT("/properties/:propertyId", controller.PutPropertiesPropertyId)
	v1.DELETE("/properties/:propertyId", controller.DeletePropertiesPropertyId)

	v1.GET("/kind_properties", controller.GetKindProperties)
	v1.GET("/kind_properties/:kindPropertyId", controller.GetKindPropertiesKindPropertyId)
	v1.POST("/kind_properties", controller.PostKindProperties)
	v1.PUT("/kind_properties/:kindPropertyId", controller.PutKindPropertiesKindPropertyId)
	v1.DELETE("/kind_properties/:kindPropertyId", controller.DeleteKindPropertiesKindPropertyId)

	v1.GET("/test", controller.GetTest)

	pEngine.NoRoute(func(c *gin.Context) {
		c.String(404, "404 page not found")
	})

	return pEngine
}
