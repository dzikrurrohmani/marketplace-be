package delivery

import (
	"store/config"
	"store/delivery/controller"
	"store/manager"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type appServer struct {
	usecaseManager manager.UsecaseManager
	engine         *gin.Engine
	host           string
}

func Server() *appServer {
	engine := gin.Default()
	engine.SetTrustedProxies([]string{"localhost"})
	engine.Use(corsConfiguration())
	appConfig := config.NewConfig()

	infra := manager.NewInfra(appConfig)
	repoManager := manager.NewRepositoryManager(infra)
	usecaseManager := manager.NewUsecaseManager(repoManager)

	host := appConfig.Url
	return &appServer{
		usecaseManager: usecaseManager,
		engine:         engine,
		host:           host,
	}
}

func (a *appServer) iniController() {
	router := a.engine.Group("api")
	controller.StartUserController(router, a.usecaseManager.UserUsecase())
	controller.StartProductController(router, a.usecaseManager.ProductUsecase(), a.usecaseManager.UserUsecase())
	controller.StartCategoryController(router, a.usecaseManager.CategoryUsecase(), a.usecaseManager.UserUsecase())
	controller.StartTransactionController(router, a.usecaseManager.TransactionUsecase(), a.usecaseManager.UserUsecase())

}

func (a *appServer) Run() {
	a.iniController()
	err := a.engine.Run(a.host)
	if err != nil {
		panic(err)
	}
}

func corsConfiguration() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowAllOrigins:  true,
		AllowHeaders:     []string{"Origin", "Date", "Content-Length", "Content-Type", "Content-Disposition", "Accept", "X-Requested-With", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Access-Control-Request-Method", "Access-Control-Request-Headers", "Authorization", "token"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	})
}
