package main

import (
	"altair/configs"
	"altair/pkg/controller"
	"altair/pkg/logger"
	"altair/pkg/manager"
	"altair/pkg/middleware"
	"altair/server"
	"flag"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	var fileConfig string
	flag.StringVar(&fileConfig, "config", "config-debug.json", "work conf file")
	flag.Parse()

	if err := configs.Load(fileConfig); err != nil {
		log.Fatalln(err.Error())
	}

	confDB := configs.Cfg.DB
	if err := server.InitDB(confDB.User, confDB.Password, confDB.Host, confDB.Port, confDB.Name); err != nil {
		log.Fatalln(err.Error())
	}

	ioWriterLogInfo, ioWriterLogWarn, ioWriterLogError := os.Stdout, os.Stdout, os.Stderr

	if configs.Cfg.Mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)

		f, err := os.OpenFile("./errors.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalln(err.Error())
		}
		defer f.Close()
		ioWriterLogWarn, ioWriterLogError = f, f
	}
	// if configs.Cfg.Mode == gin.DebugMode {
	//	 server.Db.LogMode(true)
	// }

	logger.Init(ioWriterLogInfo, ioWriterLogWarn, ioWriterLogError, configs.Cfg.Mode == gin.DebugMode)
	route := setupRouter()

	if err := route.Run("0.0.0.0:8080"); err != nil {
		log.Println(err.Error())
	}
}
func setupRouter() *gin.Engine {
	// проверим наличие папок, иначе картинки сохранятся не будут
	if !manager.FolderOrFileExists(manager.DirImages) {
		if err := os.MkdirAll(manager.DirImages, os.ModePerm); err != nil {
			logger.Error.Fatalln(err.Error())
		}
	}
	if !manager.FolderOrFileExists(manager.DirResample) {
		if err := os.MkdirAll(manager.DirResample, os.ModePerm); err != nil {
			logger.Error.Fatalln(err.Error())
		}
	}

	route := gin.Default()
	// route.MaxMultipartMemory = 16 << 20 // 16 MiB. Lower memory limit for multipart forms (default is 32 MiB)

	route.Static("/images", manager.DirImages)
	route.Static("/resample", manager.DirResample)
	route.Use(CORSMiddleware())
	route.Use(middleware.RoleIs())

	v1 := route.Group("/api/v1")
	onlyAuth := route.Group("/api/v1").Use(middleware.Auth())
	onlyAdmin := route.Group("/api/v1").Use(middleware.Auth(), middleware.Admin())

	v1.GET("", controller.GetMain)

	v1.GET("/cats", controller.GetCats)
	v1.GET("/cats/:catID", controller.GetCatsCatID)
	onlyAdmin.POST("/cats", controller.PostCats)
	onlyAdmin.PUT("/cats/:catID", controller.PutCatsCatID)
	onlyAdmin.DELETE("/cats/:catID", controller.DeleteCatsCatID)

	onlyAdmin.GET("/users", controller.GetUsers)
	onlyAdmin.GET("/users/:userID", controller.GetUsersUserID)
	onlyAdmin.POST("/users", controller.PostUsers)
	onlyAdmin.PUT("/users/:userID", controller.PutUsersUserID)
	onlyAdmin.DELETE("/users/:userID", controller.DeleteUsersUserID)

	v1.GET("/ads", controller.GetAds)
	v1.GET("/ads/:adID", controller.GetAdsAdID)
	onlyAuth.POST("/ads", controller.PostAds)
	onlyAuth.PUT("/ads/:adID", controller.PutAdsAdID)
	onlyAuth.DELETE("/ads/:adID", controller.DeleteAdsAdID)

	v1.GET("/props", controller.GetProps) // нужно при доб./изменении объявления
	onlyAdmin.GET("/props/:propID", controller.GetPropsPropID)
	onlyAdmin.POST("/props", controller.PostProps)
	onlyAdmin.PUT("/props/:propID", controller.PutPropsPropID)
	onlyAdmin.DELETE("/props/:propID", controller.DeletePropsPropID)

	onlyAdmin.GET("/kind_props", controller.GetKindProps)
	onlyAdmin.GET("/kind_props/:kindPropID", controller.GetKindPropsKindPropID)
	onlyAdmin.POST("/kind_props", controller.PostKindProps)
	onlyAdmin.PUT("/kind_props/:kindPropID", controller.PutKindPropsKindPropID)
	onlyAdmin.DELETE("/kind_props/:kindPropID", controller.DeleteKindPropsKindPropID)

	onlyAuth.GET("/auth/logout", controller.GetAuthLogout)
	v1.POST("/auth/login", controller.PostAuthLogin)
	onlyAuth.POST("/auth/refresh-tokens", controller.PostAuthRefreshTokens)

	onlyAuth.GET("/profile/ads/:adID", controller.GetProfileAdsAdID)
	onlyAuth.GET("/profile/ads", controller.GetProfileAds)
	onlyAuth.GET("/profile/settings", controller.GetProfileSettings)
	onlyAuth.GET("/profile/info", controller.GetProfileInfo)
	onlyAuth.PUT("/profile/phone/:number/:code", controller.PutProfilePhoneNumberCode)
	onlyAuth.POST("/profile/phone", controller.PostProfilePhone)
	onlyAuth.DELETE("/profile/phone/:number", controller.DeleteProfilePhoneNumber)
	onlyAuth.PUT("/profile", controller.PutProfile)
	onlyAuth.DELETE("/profile", controller.DeleteProfile)
	v1.POST("/profile", controller.PostProfile)
	v1.GET("/profile/check-email-through/:hash", controller.GetProfileCheckEmailThroughHash)

	v1.GET("/pages/ad/:adID", controller.GetPagesAdAdID)
	v1.GET("/pages/main", controller.GetPagesMain)

	v1.POST("/recover/send-hash", controller.PostRecoverSendHash)     // отправитель хеша на почту
	v1.POST("/recover/change-pass", controller.PostRecoverChangePass) // приемник нового пароля

	v1.GET("/search/ads", controller.GetSearchAds)
	onlyAdmin.GET("/test", controller.GetTest)
	// v1.GET("/test", controller.GetTest)
	v1.GET("/resample/:width/:height/*path", controller.GetResampleWidthHeightPath)
	v1.GET("/phones/:phoneID", controller.GetPhonesPhoneID)

	route.NoRoute(func(c *gin.Context) {
		c.String(404, manager.Err404.Error())
	})

	return route
}
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//	"https://www.altair.uz", "http://localhost:4200",
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		//[]string{"Content-Type", "Authorization"}
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "OPTIONS, GET, POST, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
