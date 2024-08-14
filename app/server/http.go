package server

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/krisdioles/ppr-wallet/app/provider"
	"github.com/krisdioles/ppr-wallet/app/server/controller"
	"github.com/krisdioles/ppr-wallet/config"
)

func InitHttpServer(cfg *config.ServerConfig, usecase *provider.Usecase) {
	router := gin.Default()

	userBalanceController := controller.NewUserBalanceController(usecase.UserBalanceUsecase)
	router.GET("/api/user-balance/:id", userBalanceController.GetUserBalanceByID)
	router.PATCH("/api/user-balance/:id/disburse", userBalanceController.DisburseBalance)

	router.Run(fmt.Sprintf("localhost:%d", cfg.Port))
}
