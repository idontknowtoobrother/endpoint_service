package main

import (
	"context"

	"github.com/gin-gonic/gin"
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
	serviceService := service.NewEndpointService(ctx, serviceRepository)
	serviceHandler := handler.NewEndpointHandler(serviceService)

	// TODO: Implement router and server

	router := gin.Default()
	router.GET("/endpoints", serviceHandler.GetEndpoints)
	router.GET("/endpoints/:uuid", serviceHandler.GetEndpoint)

}
