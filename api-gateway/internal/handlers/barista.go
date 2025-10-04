package handlers

import (
	pb_barista "api-gateway/gen/barista"
	"context"
	"github.com/gofiber/fiber/v2"
	"time"
)

func StartPreparingHandler(baristaClient pb_barista.BaristaServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orderID := c.Params("order_id")

		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		resp, err := baristaClient.StartPreparing(ctx, &pb_barista.StartPreparingRequest{OrderId: orderID})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if resp.Error != "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": resp.Error})
		}

		return c.JSON(fiber.Map{"status": resp.Status})
	}
}

func OrderReadyHandler(baristaClient pb_barista.BaristaServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		orderID := c.Params("order_id")

		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		resp, err := baristaClient.OrderReady(ctx, &pb_barista.OrderReadyRequest{OrderId: orderID})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if resp.Error != "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": resp.Error})
		}

		return c.JSON(fiber.Map{"status": resp.Status})
	}
}

func GetPendingOrdersHandler(baristaClient pb_barista.BaristaServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		resp, err := baristaClient.GetPendingOrders(ctx, &pb_barista.GetPendingOrdersRequest{})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if resp.Error != "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": resp.Error})
		}

		return c.JSON(resp.Orders)
	}

}
