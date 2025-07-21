package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/ShiftyX1/Facade/internal/config"
	"github.com/ShiftyX1/Facade/internal/handlers"
	"github.com/ShiftyX1/Facade/internal/router"
	"github.com/ShiftyX1/Facade/internal/server"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

var (
	configPath string
	port       int
	host       string
	verbose    bool

	// Build info - set via ldflags
	Version   string
	BuildTime = "unknown"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "facade",
		Short: "Simple and flexible mock HTTP server",
		Long: `Facade is a simple and flexible mock HTTP server for developers.
It allows you to quickly create mock APIs through YAML configuration without writing code.`,
		Version: Version,
		Run:     runServer,
	}

	rootCmd.Flags().StringVarP(&configPath, "config", "c", "configs/example.yaml", "path to config file")
	rootCmd.Flags().IntVarP(&port, "port", "p", 0, "server port (overrides config)")
	rootCmd.Flags().StringVar(&host, "host", "", "server host (overrides config)")
	rootCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")

	if err := rootCmd.Execute(); err != nil {
		slog.Error("Failed to execute command", "error", err)
		os.Exit(1)
	}
}

func runServer(cmd *cobra.Command, args []string) {

	setupLogging()

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		slog.Error("Failed to load config", "error", err, "path", configPath)
		os.Exit(1)
	}

	if port > 0 {
		cfg.Server.Port = port
	}
	if host != "" {
		cfg.Server.Host = host
	}

	slog.Info("Facade starting",
		"version", Version,
		"build_time", BuildTime,
		"config_path", configPath,
		"host", cfg.Server.Host,
		"port", cfg.Server.Port,
		"routes", len(cfg.Routes),
	)

	appRouter := router.New(cfg)
	if err := appRouter.Setup(); err != nil {
		slog.Error("Failed to setup router", "error", err)
		os.Exit(1)
	}

	mainRouter := mux.NewRouter()

	stateHandler := handlers.NewStateHandler(appRouter.GetStateManager())
	mainRouter.HandleFunc("/_facade/health", handlers.HealthHandler).Methods("GET")
	mainRouter.HandleFunc("/_facade/state", stateHandler.GetState).Methods("GET")
	mainRouter.HandleFunc("/_facade/state", stateHandler.ClearState).Methods("DELETE")
	mainRouter.HandleFunc("/_facade/state/keys", stateHandler.GetStateKeys).Methods("GET")

	mainRouter.PathPrefix("/").Handler(appRouter.Handler())

	srv := server.New(cfg, mainRouter)
	if err := srv.Start(); err != nil {
		slog.Error("Server error", "error", err)
		os.Exit(1)
	}
}

func setupLogging() {
	level := slog.LevelInfo
	if verbose {
		level = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	slog.SetDefault(logger)
}

func init() {

	flag.StringVar(&configPath, "config", "configs/example.yaml", "path to config file")
	flag.IntVar(&port, "port", 0, "server port")
	flag.StringVar(&host, "host", "", "server host")
	flag.BoolVar(&verbose, "verbose", false, "enable verbose logging")
}
