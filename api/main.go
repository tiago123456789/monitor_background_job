package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/tiago123456789/monitor_background_job/config"
	"github.com/tiago123456789/monitor_background_job/models"
	"github.com/tiago123456789/monitor_background_job/queue"
	"github.com/tiago123456789/monitor_background_job/repositories"
	"github.com/tiago123456789/monitor_background_job/services"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	redisClient := config.NewRedisClient()
	mongoClient := config.NewConnection()
	cache := config.NewCache(redisClient)

	producer := queue.NewProducer()

	eventNotificationRepository := repositories.NewEventNotificationRepository(
		cache, producer, mongoClient)
	companyRepository := repositories.NewCompanyRepostory(mongoClient)
	jobRepository := repositories.NewJobRepository(mongoClient, companyRepository)
	alertRepository := repositories.NewAlertRepository(mongoClient)

	authService := services.NewAuth(companyRepository)

	app := fiber.New()

	app.Use(cors.New(cors.ConfigDefault))

	hasPermission := func(c *fiber.Ctx) error {
		type ForbiddenResponse struct {
			Message string
		}

		accessToken := c.Get("Authorization")
		if accessToken == "" {
			return c.Status(403).JSON(ForbiddenResponse{
				Message: "Is necessary informate accessToken",
			})
		}

		accessTokenWithoutPrefix := strings.ReplaceAll(accessToken, "Bearer ", "")
		error := authService.IsAuthenticated(accessTokenWithoutPrefix)
		if error != nil {
			return c.Status(403).JSON(ForbiddenResponse{
				Message: "Is necessary informate accessToken",
			})
		}

		return c.Next()
	}

	app.Post("/auth/login", func(c *fiber.Ctx) error {
		credential := new(models.Credential)
		if err := c.BodyParser(credential); err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": err,
			})
		}
		accessToken, err := authService.Login(*credential)
		if err != nil {
			fmt.Print(err)
			return c.Status(401).JSON(&fiber.Map{
				"success": false,
				"message": "Credentials invalid",
			})
		}
		return c.JSON(accessToken)
	})

	app.Post("/companies", func(c *fiber.Ctx) error {
		company := new(models.Company)
		if err := c.BodyParser(company); err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": err,
			})
		}

		err := companyRepository.Create(*company)
		if err != nil {
			return c.Status(409).JSON(&fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		return c.SendStatus(201)
	})

	app.Post("/companies/:companyId/jobs", hasPermission, func(c *fiber.Ctx) error {
		job := new(models.JobModel)
		if err := c.BodyParser(job); err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		err := jobRepository.Create(models.JobModel{
			Name:      job.Name,
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

	app.Post("/companies/:companyId/jobs/:jobId/alerts", hasPermission, func(c *fiber.Ctx) error {
		alert := new(models.Alert)
		if err := c.BodyParser(alert); err != nil {
			return c.Status(400).JSON(&fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		alert.CompanyId = c.Params("companyId")
		err := alertRepository.Create(*alert)
		if err != nil {
			fmt.Print(err)
			return c.Status(409).JSON(&fiber.Map{
				"success": false,
				"message": err.Error(),
			})
		}
		return c.SendStatus(201)
	})

	app.Get("/companies/:companyId/jobs", hasPermission, func(c *fiber.Ctx) error {
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

	app.Get("/job-notifications-received/:jobId", hasPermission, func(c *fiber.Ctx) error {
		notifications, err := eventNotificationRepository.GetByJobID(c.Params("jobId"))
		if err != nil {
			return c.Status(200).SendString("[]")
		}
		return c.JSON(notifications)
	})

	app.Get("/event-notifications/:id", hasPermission, func(c *fiber.Ctx) error {
		data, err := eventNotificationRepository.Get(c.Params("id"))
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.SendString(data)
	})

	app.Listen(":4000")
}
