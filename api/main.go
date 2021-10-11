package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/tiago123456789/monitor_background_job/config"
	"github.com/tiago123456789/monitor_background_job/models"
	"github.com/tiago123456789/monitor_background_job/queue"
	"github.com/tiago123456789/monitor_background_job/repositories"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	client := config.NewConnection()
	cache := config.NewCache(rdb)

	producer := queue.NewProducer()

	eventNotificationRepository := repositories.NewEventNotificationRepository(cache, producer)
	companyRepository := repositories.NewCompanyRepostory(client)
	jobRepository := repositories.NewJobRepository(client, companyRepository)

	app := fiber.New()

	app.Post("/companies", func(c *fiber.Ctx) error {
		company := new(models.Company)
		if err := c.BodyParser(company); err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": err,
			})
		}
		err := companyRepository.Create(models.Company{
			Name:     company.Name,
			Email:    company.Email,
			Password: company.Password,
		})
		if err != nil {
			fmt.Print(err)
			return c.Status(409).JSON(&fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		return c.SendStatus(201)
	})

	app.Post("/companies/:companyId/jobs", func(c *fiber.Ctx) error {
		job := new(models.JobModel)
		if err := c.BodyParser(job); err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		err := jobRepository.Create(models.JobModel{
			Name:      job.Name,
			Interval:  job.Interval,
			Companyid: c.Params("companyId"),
		})
		if err != nil {
			fmt.Print(err)
			return c.Status(409).JSON(&fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		return c.SendStatus(201)
	})

	app.Get("/companies/:companyId/jobs", func(c *fiber.Ctx) error {
		jobs, err := jobRepository.FindByCompanyId(c.Params("companyId"))
		if err != nil {
			return c.SendStatus(404)
		}
		return c.JSON(jobs)
	})

	app.Get("/event-notifications/:id/:jobId", func(c *fiber.Ctx) error {
		err := eventNotificationRepository.StoreLast(c.Params("id"), c.Params("jobId"))
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.SendStatus(200)
	})

	app.Get("/event-notifications/:id", func(c *fiber.Ctx) error {
		data, err := eventNotificationRepository.Get(c.Params("id"))
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.SendString(data)
	})

	app.Listen(":4000")
}
