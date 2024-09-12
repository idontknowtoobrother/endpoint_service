package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	conn "github.com/idontknowtoobrother/practice_go_hexagonal/connections"
	handler "github.com/idontknowtoobrother/practice_go_hexagonal/handlers/endpoint"
	repo "github.com/idontknowtoobrother/practice_go_hexagonal/repositories/endpoint"
	service "github.com/idontknowtoobrother/practice_go_hexagonal/services/endpoint"
)

func main() {
	ctx := context.Background()
	mongoClient := conn.NewDatabaseConnection(ctx)
	serviceDatabase := mongoClient.Database("endpoint_service")
	serviceCollection := serviceDatabase.Collection("endpoints")

	serviceRepository := repo.NewEndpointRepository(ctx, serviceCollection)
	// micro client
	restyClient := resty.New()
	serviceService := service.NewEndpointService(ctx, serviceRepository, restyClient)
	serviceHandler := handler.NewEndpointHandler(serviceService)

	// TODO: Implement router and server

	router := gin.Default()
	apiV1 := router.Group("/api/v1")
	apiV1.GET("/endpoints", serviceHandler.GetEndpoints)
	apiV1.POST("/endpoints.new", serviceHandler.CreateEndpoint)
	apiV1.PUT("/endpoints.update", serviceHandler.UpdateEndpoint)
	apiV1.DELETE("/endpoints.delete", serviceHandler.DeleteEndpoint)
	apiV1.POST("/endpoints.send/:path", serviceHandler.SendToWebhook)

	// Set up the HTTP server
	srv := &http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	// Channel to listen for interrupt or terminate signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start the server in a goroutine so it doesn't block
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for the signal
	<-quit
	log.Println("shutting down server...")

}
