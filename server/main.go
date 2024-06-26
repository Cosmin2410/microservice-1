package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/microservice/server/db"
	"github.com/microservice/server/domain"
)

var (
	DBConn *sql.DB
)

func setupRouter(app *fiber.App) {
	app.Get("/api/v1/users", func(c *fiber.Ctx) error {
		users, err := db.Find(DBConn)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
		return c.JSON(users)
	})

	app.Post("/api/v1/post", func(c *fiber.Ctx) error {
		var user domain.User
		if err := c.BodyParser(&user); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		id, err := db.Insert(DBConn, user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		user.ID = id
		return c.Status(fiber.StatusCreated).JSON(user)
	})
}

func main() {
	var err error

	dsn := os.Getenv("DB_PASS")
	if dsn == "" {
		log.Fatal("DB_PASS environment variable is not set")
	}

	DBConn, err = db.SetupDBConn(dsn)
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow all origins for development
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	setupRouter(app)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	log.Fatal(app.Listen(":8000"))
}
