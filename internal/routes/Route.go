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

	router.Get("/movies", movieHandler.GetAllMovies)
	router.Get("/movies/search", movieHandler.GetMoviesFilter)
	router.Get("/movies/page", movieHandler.GetMoviesPaginated)
	router.Get("/movies/stats", movieHandler.GetMovieStats)
	router.Get("/movies/:id", movieHandler.GetMovieByID)

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

	admin := router.Group("/schedules")
	admin.Use(middleware.AdminOnly())

	admin.Post("/", scheduleHandler.CreateSchedule)
	admin.Put("/:id", scheduleHandler.UpdateSchedule)
	admin.Patch("/:id", scheduleHandler.PatchSchedule)
	admin.Delete("/:id", scheduleHandler.DeleteSchedule)
}

func TicketRoutes(router fiber.Router) {
	ticketRepo := repository.NewTicketRepository()
	scheduleRepo := repository.NewScheduleRepository()
	ticketService := services.NewTicketService(ticketRepo, scheduleRepo)
	ticketHandler := handlers.NewTicketHandler(ticketService)

	admin := router.Group("/tickets")
	admin.Use(middleware.AdminOnly())
	admin.Get("/", ticketHandler.GetTickets)

	protected := router.Group("/tickets", middleware.Protected())
	protected.Post("/buy", ticketHandler.BuyTicket)
	protected.Post("/refund", ticketHandler.RefundTicket)
}

func CustomerRoutes(router fiber.Router) {
	customerRepo := repository.NewCustomerRepository()
	customerService := services.NewCustomerService(customerRepo)
	customerHandler := handlers.NewCustomerHandler(customerService)

	router.Post("/customers/register", customerHandler.Register)
	router.Post("/customers/login", customerHandler.Login)

	admin := router.Group("/customers")
	admin.Use(middleware.AdminOnly())
	admin.Get("/", customerHandler.GetAllCustomers)

	protected := router.Group("/customers", middleware.Protected())
	protected.Get("/:id", customerHandler.GetProfile)
	protected.Patch("/:id", customerHandler.PatchCustomer)
	protected.Post("/:id/topup", customerHandler.TopUpWallet)
	protected.Delete("/:id", customerHandler.DeleteCustomer)
}
