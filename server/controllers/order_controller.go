package controllers

import (
	"context"
	"net/http"
	"server/configs"
	"server/models"
	"server/responses"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var orderCollection *mongo.Collection = configs.GetCollection(configs.DB, "orders")
var validate = validator.New()

func AddOrder() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var order models.Order
		defer cancel()

		// validate the request body
		if err := c.BindJSON(&order); err != nil {
			c.JSON(http.StatusBadRequest, responses.OrderResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}

		// use the validate library to validate fields
		if validationErr := validate.Struct(&order); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.OrderResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data: map[string]interface{}{
					"data": validationErr.Error(),
				},
			})
			return
		}

		newOrder := models.Order{
			ID:     primitive.NewObjectID(),
			Dish:   order.Dish,
			Price:  order.Price,
			Server: order.Server,
			Table:  order.Table,
		}

		result, err := orderCollection.InsertOne(ctx, newOrder)
		if err != nil {
			c.JSON(http.StatusBadRequest, responses.OrderResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}

		c.JSON(http.StatusCreated, responses.OrderResponse{
			Status:  http.StatusCreated,
			Message: "success",
			Data: map[string]interface{}{
				"data": result,
			},
		})
	}
}

func GetAllOrders() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var orders []models.Order
		defer cancel()

		results, err := orderCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.OrderResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}

		// reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleOrder models.Order
			if err := results.Decode(&singleOrder); err != nil {
				c.JSON(http.StatusInternalServerError, responses.OrderResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data: map[string]interface{}{
						"data": err.Error(),
					},
				})
			}
			orders = append(orders, singleOrder)
		}

		c.JSON(http.StatusOK, responses.OrderResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data: map[string]interface{}{
				"data": orders,
			},
		})
	}
}
