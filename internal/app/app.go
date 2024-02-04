package app

import (
	"EWallet/internal/config"
	"EWallet/internal/repository/postgresql"
	"EWallet/internal/transport/http/handlers/get"
	"EWallet/internal/transport/http/handlers/post"
	"EWallet/internal/usecase"
	"EWallet/pkg/logger/handlers/slogpretty"
	"EWallet/pkg/server"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func Start() error {
	// you can set CONFIG_PATH like this or use set config_path in your terminal

	//err := os.Setenv("CONFIG_PATH", "config.yaml")
	//if err != nil {
	//	log.Fatal("did not set env: CONFIG_PATH")
	//}

	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatal("troubles with config initialization:", err)
	}

	logger, err := setupLogger(cfg.Logger.LogLevel)
	if err != nil {
		log.Fatal("troubles with logger initialization:", err)
	}

	router := mux.NewRouter()
	httpServer := server.New(
		router,
		cfg.HTTPServer.Host,
		cfg.HTTPServer.Port,
	)

	pgDatabase, err := postgresql.NewPGSQLRepo(
		cfg.PGSQL.UserName,
		cfg.PGSQL.Password,
		cfg.PGSQL.Host,
		cfg.PGSQL.Port,
		cfg.PGSQL.DbName,
		cfg.PGSQL.PingTime,
	)
	if err != nil {
		logger.Error("troubles with pgsql initialization:", err)
	}

	databaseUseCase := usecase.NewRepository(pgDatabase)

	router.HandleFunc("/api/v1/wallet", post.CreateBalance(databaseUseCase, logger)).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/wallet/{walletID}/send", post.SendMoney(databaseUseCase, logger)).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/wallet/{walletID}/history", get.WalletHistory(databaseUseCase, logger)).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/wallet/{walletID}", get.WalletBalance(databaseUseCase, logger)).Methods(http.MethodGet)

	wait(httpServer)

	return nil
}

func wait(s *server.Server) {
	log.Println("server is starting")

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	select {
	case i := <-ch:
		log.Println("shutdown signal:" + i.String())
	case err := <-s.Notify():
		log.Fatal("wait - server.Notify", err)
	}
	log.Println("App is stopping...")
}

func setupLogger(env string) (*slog.Logger, error) {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = setupPrettySlog()
	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		return nil, fmt.Errorf("wrong env value")
	}
	return logger, nil
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
