package api

import (
	fiber "github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	model "github.com/adzpm/tg-clicker/internal/model"
)

func (a *API) Login(c *fiber.Ctx) error {
	// get telegram_id from body

	tid := c.QueryInt("telegram_id")

	if tid == 0 {
		if err := c.Status(400).JSON(fiber.Map{"error": "telegram_id is required"}); err != nil {
			return err
		}

		return nil
	}

	a.lgr.Info("try to login", zap.Int("telegram_id", tid))

	// get user by telegram_id
	// if user not exists, create new user
	// return user

	u, err := a.str.GetUserByTelegramID(uint64(tid))

	if err != nil {
		u, err = a.str.CreateUser(&model.User{
			TelegramID: uint64(tid),
			Coins:      5,
		})

		if err != nil {
			if err = c.Status(500).JSON(fiber.Map{"error": err.Error()}); err != nil {
				return err
			}

			return err
		}
	}

	if u, err = a.str.Login(uint64(tid)); err != nil {
		if err = c.Status(500).JSON(fiber.Map{"error": err.Error()}); err != nil {
			return err
		}

		return err
	}

	return c.Status(200).JSON(u)
}

func (a *API) Click(c *fiber.Ctx) error {
	// get telegram_id from body
	// get user by telegram_id
	// add coins to user
	// return user

	tid := c.QueryInt("telegram_id")

	if tid == 0 {
		if err := c.Status(400).JSON(fiber.Map{"error": "telegram_id is required"}); err != nil {
			return err
		}

		return nil
	}

	a.lgr.Info("try to click", zap.Int("telegram_id", tid))

	u, err := a.str.AddCoinsByTelegramID(uint64(tid), uint64(1))

	if err != nil {
		if err = c.Status(500).JSON(fiber.Map{"error": err.Error()}); err != nil {
			return err
		}

		return err
	}

	return c.Status(200).JSON(u)
}
