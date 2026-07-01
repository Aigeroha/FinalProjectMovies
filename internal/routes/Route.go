package routes

import (
	"final-project/internal/handlers"
	"final-project/internal/middleware"
	"final-project/internal/repository"
	"final-project/internal/services"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")

	MovieRoutes(api)
	ScheduleRoutes(api)
	TicketRoutes(api)
	CustomerRoutes(api)

	app.Use("/ws", func(c fiber.Ctx) error {
		if c.Get("Upgrade") == "websocket" {
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Get("/ws/halls", websocket.New(func(c *websocket.Conn) {
		_ = handlers.HandlerWebSocket(c)
	}))

}

func MovieRoutes(router fiber.Router) {
    movieRepo := repository.NewMovieRepository()
    movieService := services.NewMovieService(movieRepo)
    movieHandler := handlers.NewMovieHandler(movieService)

    // --- ПУБЛИЧНЫЕ (доступны всем без токена и ключа) ---
    router.Get("/movies", movieHandler.GetAllMovies)
    router.Get("/movies/search", movieHandler.GetMoviesFilter)
    router.Get("/movies/page", movieHandler.GetMoviesPaginated)
    router.Get("/movies/stats", movieHandler.GetMovieStats)
    router.Get("/movies/:id", movieHandler.GetMovieByID)

    // --- АДМИНСКИЕ (только с секретным ключом) ---
    // Используем middleware.AdminOnly(), который мы создадим ниже
    admin := router.Group("/movies")
    admin.Use(middleware.AdminOnly()) 

    admin.Post("/", movieHandler.CreateMovie)
    admin.Put("/:id", movieHandler.UpdateMovie)
    admin.Patch("/:id", movieHandler.PatchMovie)
    admin.Delete("/:id", movieHandler.DeleteMovie)
}

func ScheduleRoutes(router fiber.Router) {
	scheduleRepo := repository.NewScheduleRepository()
	scheduleService := services.NewScheduleService(scheduleRepo)
	scheduleHandler := handlers.NewScheduleHandler(scheduleService)

	router.Get("/schedules", scheduleHandler.GetSchedules)
	router.Get("/schedules/page", scheduleHandler.GetSchedulesPaginated)

	router.Post("/schedules", scheduleHandler.CreateSchedule)
	router.Put("/schedules/:id", scheduleHandler.UpdateSchedule)
	router.Patch("/schedules/:id", scheduleHandler.PatchSchedule)
	router.Delete("/schedules/:id", scheduleHandler.DeleteSchedule)
}

func TicketRoutes(router fiber.Router) {
	ticketRepo := repository.NewTicketRepository()
	scheduleRepo := repository.NewScheduleRepository()
	ticketService := services.NewTicketService(ticketRepo, scheduleRepo)
	ticketHandler := handlers.NewTicketHandler(ticketService)

	router.Get("/tickets", ticketHandler.GetTickets)

	protected := router.Group("", middleware.Protected())
	protected.Post("/tickets", ticketHandler.BuyTicket)
	protected.Post("/tickets/refund", ticketHandler.RefundTicket)
}

func CustomerRoutes(router fiber.Router) {
	customerRepo := repository.NewCustomerRepository()
	customerService := services.NewCustomerService(customerRepo)
	customerHandler := handlers.NewCustomerHandler(customerService)

	router.Post("/customers/register", customerHandler.Register)
	router.Post("/customers/login", customerHandler.Login)

	protected := router.Group("/customers")
	protected.Use(middleware.Protected())

	protected.Get("/", customerHandler.GetAllCustomers)
	protected.Get("/:id", customerHandler.GetProfile)
	protected.Patch("/:id", customerHandler.PatchCustomer)
	protected.Post("/:id/topup", customerHandler.TopUpWallet)
	protected.Delete("/:id", customerHandler.DeleteCustomer)
}
