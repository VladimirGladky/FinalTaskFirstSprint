package main

import (
	"context"
	"github.com/VladimirGladky/FinalTaskFirstSprint/internal/orchestrator/server"
	"github.com/VladimirGladky/FinalTaskFirstSprint/pkg/logger"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	ctx, _ = logger.New(ctx)
	_ = godotenv.Load("local.env")
	srv := server.New(ctx)
	logger.GetLoggerFromCtx(ctx).Info(ctx, "Orchestrator started")
	srv.Run()
}
