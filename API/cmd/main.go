package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/bjs3118/HackZurich2921/API"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

func main() {

	fmt.Println("Loading Server")

	var httpPort = flag.String("httpPort", ":8080", "Port for serving http server")
	var httpServerTLSCertFileName = flag.String("httpServerTLSCertFileName", "cert/server.crt", "File path of TLS HTTP server certificate")
	var httpServerTLSKeyFileName = flag.String("httpServerTLSKeyFileName", "cert/server.key", "File path of TLS HTTP server key")
	var serverDBFilePath = flag.String("db", "serverDB.db", "SQLite DB file path")
	flag.Parse()

	serverDBDSN := "db/" + *serverDBFilePath

	// Context

	ctx := context.Background()

	// Logging

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("server: failed to create zap logger: %v\n", err)
	}
	defer logger.Sync()

	// SQL lite database

	serverDB, err := server.OpenSQLiteDB(ctx, logger, serverDBDSN)
	if err != nil {
		logger.Fatal("server: failed to open SQLite server database", zap.Error(err))
	}
	logger.Info("server: opened server sqlite3 DB")

	r := chi.NewRouter()

	httpServer := server.OpenHttpServer(ctx, logger, r, serverDB)
	defer httpServer.Close()

	logger.Info("server: opened http server")

	if err := httpServer.Serve(ctx, *httpPort, *httpServerTLSCertFileName, *httpServerTLSKeyFileName); err != nil {
		logger.Fatal("server: failed to serve http server", zap.Error(err))
	}
}
