package handlers

import (
	"api-gateway/gen/order"
	"context"
	"github.com/gofiber/fiber/v2"
	"time"
)

func CreateOrderHandler(orderClient order.OrderServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			Item string `json:"item"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		resp, err := orderClient.CreateOrder(ctx, &order.CreateOrderRequest{Item: req.Item})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(fiber.Map{"order_id": resp.OrderId, "status": resp.Status})
	}
}
