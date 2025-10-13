package handlers

import (
	pb_notification "api-gateway/gen/notification"
	"context"
	"github.com/gofiber/fiber/v2"
	"time"
)

func GetSentNotificationsHandler(notificationClient pb_notification.NotificationServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		since := c.Query("since", "")

		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		resp, err := notificationClient.GetSentNotifications(ctx, &pb_notification.GetSentNotificationsRequest{Since: since})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if resp.Error != "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": resp.Error})
		}

		return c.JSON(resp.Notifications)
	}
}

func TriggerManualNotificationHandler(notificationClient pb_notification.NotificationServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req struct {
			OrderID string `json:"order_id"`
			Item    string `json:"item"`
			Status  string `json:"status"`
		}

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		resp, err := notificationClient.TriggerManualNotification(ctx, &pb_notification.TriggerManualNotificationRequest{
			OrderId: req.OrderID,
			Item:    req.Item,
			Status:  req.Status,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		if resp.Error != "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": resp.Error})
		}

		return c.JSON(fiber.Map{"status": resp.Status})
	}

}
