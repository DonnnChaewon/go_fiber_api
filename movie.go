package main

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func MovieHandlers(route fiber.Router, db *gorm.DB) {
	route.Get("/", func(c *fiber.Ctx) error {
		title := c.Query("title")
		status := c.Query("status")
		director := c.Query("director")

		userId := int(c.Locals("userId").(float64))
		movies := new([]Movie)

		query := db.Where("user_id = ?", userId)

		if title != "" {
			query.Where("title LIKE ?", "%"+title+"%")
		}
		if status != "" {
			query.Where("status = ?", status)
		}
		if director != "" {
			query.Where("director = ?", director)
		}

		if err := query.Find(&movies).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Movie not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(movies)
	})

	route.Get("/:id", func(c *fiber.Ctx) error {
		movieId, _ := c.ParamsInt("id")
		userId := int(c.Locals("userId").(float64))
		movie := new(Movie)

		if err := db.Where("id = ? AND user_id = ?", movieId, userId).First(&movie).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Movie not found",
			})
		}

		return c.Status(fiber.StatusOK).JSON(movie)
	})

	route.Post("/", func(c *fiber.Ctx) error {
		movie := new(Movie)
		movie.UserID = int(c.Locals("userId").(float64))

		if err := c.BodyParser(movie); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := db.Create(&movie).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(movie)
	})

	route.Put("/:id", func(c *fiber.Ctx) error {
		movieId, _ := c.ParamsInt("id")
		userId := int(c.Locals("userId").(float64))
		movie := new(Movie)

		if err := db.Where("id = ? AND user_id = ?", movieId, userId).First(&movie).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Movie not found",
			})
		}

		if err:= c.BodyParser(movie); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		if err := db.Save(&movie).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(movie)
	})

	route.Delete("/:id", func(c *fiber.Ctx) error {
		movieId, _ := c.ParamsInt("id")
		userId := int(c.Locals("userId").(float64))
		movie := new(Movie)

		if err := db.Where("id = ? AND user_id = ?", movieId, userId).First(&movie).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Movie not found",
			})
		}

		if err := db.Delete(&movie).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})
}