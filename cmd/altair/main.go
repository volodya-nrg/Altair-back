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

/*
ЗАМЕТКИ:
 - фото хранятся в табл images и в value_props (id записей в поле value). POST/PUT "files" - зарез-но
 - sort - https://gobyexample.com/sorting
 - INSERT INTO cats_props (cat_id, prop_id, pos, is_require, is_can_as_filter, `comment`) VALUES
 	- (x, 91, 3, 1, 0, '5'),\n
 - SHOW VARIABLES WHERE variable_name = 'max_user_connections' (10)
 - https://github.com/golang/go/wiki/SliceTricks
*/

func init() {
	// gin.SetMode(gin.ReleaseMode)

	logger.Init(os.Stdout, os.Stdout, os.Stderr)

	if err := configs.Load("./config.json"); err != nil {
		logger.Error.Fatalln(err.Error())
	}

	confDB := configs.Cfg.DB
	if err := server.InitDB(confDB.User, confDB.Password, confDB.Host, confDB.Port, confDB.Name); err != nil {
		logger.Error.Fatalln(err.Error())
	}
}
func main() {
	route := setupRouter()
	if err := route.Run("127.0.0.1:8080"); err != nil {
		logger.Error.Fatalln(err.Error())
	}
}
func setupRouter() *gin.Engine {
	configCors := cors.DefaultConfig()
	configCors.AllowOrigins = []string{"http://localhost:8080", "http://localhost:4200"}
	configCors.AllowCredentials = true
	configCors.AllowMethods = []string{"GET", "POST", "PUT", "DELETE"}

	route := gin.Default()
	route.MaxMultipartMemory = 16 << 20 // 16 MiB. Lower memory limit for multipart forms (default is 32 MiB)

	route.LoadHTMLGlob("./web/templates/*")
	route.StaticFile("/favicon.ico", "./web/assets/img/favicon.ico")
	route.Static("/assets", "./web/assets")
	route.Use(static.Serve("/images", static.LocalFile("./web/images", true)))
	route.Use(cors.New(configCors))

	route.GET("/", func(c *gin.Context) {
		output := c.DefaultQuery("format", "html")
		serviceCats := service.NewCatService()
		serviceKindProps := service.NewKindPropService()
		serviceProps := service.NewPropService()

		cats, err := serviceCats.GetCats()
		if err != nil {
			logger.Warning.Println(err)
			c.JSON(500, err.Error())
			return
		}

		catsTree := serviceCats.GetCatsAsTree(cats)

		kindProps, err := serviceKindProps.GetKindProps("kind_prop_id asc")
		if err != nil {
			logger.Warning.Println(err)
			c.JSON(500, err.Error())
			return
		}

		props, err := serviceProps.GetProps("title asc")
		if err != nil {
			logger.Warning.Println(err)
			c.JSON(500, err.Error())
			return
		}

		data := gin.H{
			"catsTree":  catsTree,
			"kindProps": kindProps,
			"props":     props,
		}

		if output == "json" {
			c.JSON(http.StatusOK, data)
			return
		}

		c.HTML(http.StatusOK, "index.html", data)
	})

	v1 := route.Group("/api/v1")

	v1.GET("/cats", controller.GetCats)
	v1.GET("/cats/:catId", controller.GetCatsCatId)
	v1.POST("/cats", controller.PostCats)
	v1.PUT("/cats/:catId", controller.PutCatsCatId)
	v1.DELETE("/cats/:catId", controller.DeleteCatsCatId)

	v1.GET("/users", controller.GetUsers)
	v1.GET("/users/:userId", controller.GetUsersUserId)
	v1.POST("/users", controller.PostUsers)
	v1.PUT("/users/:userId", controller.PutUsersUserId)
	v1.DELETE("/users/:userId", controller.DeleteUsersUserId)

	v1.GET("/ads", controller.GetAds)
	v1.GET("/ads/:adId", controller.GetAdsAdId)
	v1.POST("/ads", controller.PostAds)
	v1.PUT("/ads/:adId", controller.PutAdsAdId)
	v1.DELETE("/ads/:adId", controller.DeleteAdsAdId)

	v1.GET("/props", controller.GetProps)
	v1.GET("/props/:propId", controller.GetPropsPropId)
	v1.POST("/props", controller.PostProps)
	v1.PUT("/props/:propId", controller.PutPropsPropId)
	v1.DELETE("/props/:propId", controller.DeletePropsPropId)

	v1.GET("/kind_props", controller.GetKindProps)
	v1.GET("/kind_props/:kindPropId", controller.GetKindPropsKindPropId)
	v1.POST("/kind_props", controller.PostKindProps)
	v1.PUT("/kind_props/:kindPropId", controller.PutKindPropsKindPropId)
	v1.DELETE("/kind_props/:kindPropId", controller.DeleteKindPropsKindPropId)

	v1.GET("/search/ads", controller.GetSearchAds)
	v1.GET("/pages/ad/:adId", controller.GetPagesAdAdId)
	v1.GET("/pages/main", controller.GetPagesMain)
	v1.GET("/test", controller.GetTest)

	route.NoRoute(func(c *gin.Context) {
		c.String(404, "404 page not found")
	})

	return route
}
