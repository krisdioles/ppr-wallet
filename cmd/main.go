package main

import (
	"github.com/krisdioles/ppr-wallet/app/infrastructure/database"
	"github.com/krisdioles/ppr-wallet/app/provider"
	"github.com/krisdioles/ppr-wallet/app/server"
	"github.com/krisdioles/ppr-wallet/config"
)

func main() {
	cfg := config.All()

	db := database.Init()
	repo := provider.InitRepositories(db)
	usecase := provider.InitUsecases(cfg, repo)

	server.InitHttpServer(&cfg.Server, usecase)
}
