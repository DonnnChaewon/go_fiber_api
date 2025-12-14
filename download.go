package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Download(route fiber.Router, db *gorm.DB) {
	route.Get("/", func(c *fiber.Ctx) error {
		format := c.Query("format", "json")
		movies := new([]Movie)

		if err := db.Find(&movies).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Fail to fetch movies",
			})
		}

		var fileName string
		
		switch format {
		case "json":
			fileName = "movie.json"

			file, err := os.Create(fileName)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Fail to create the JSON file",
				})
			}
			defer file.Close()

			encoder := json.NewEncoder(file)
			encoder.SetIndent("", " ")
			if err := encoder.Encode(movies); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Fail to write the JSON file",
				})
			}
		case "csv":
			fileName = "movie.csv"

			file, err := os.Create(fileName)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Fail to create the JSON file",
				})
			}
			defer file.Close()

			writer := csv.NewWriter(file)
			writer.Write([]string{"ID", "Title", "Status", "Director"})

			for _, movie := range *movies {
				writer.Write([]string{
					fmt.Sprintf("%d", movie.ID),
					movie.Title,
					string(movie.Status),
					movie.Director,
				})
			}

			writer.Flush()
			if err := writer.Error(); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Fail to write the CSV file",
				})
			}
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Only json and csv format are allowed",
			})
		}

		defer os.Remove(fileName)
		return c.Download(fileName)
	})
}