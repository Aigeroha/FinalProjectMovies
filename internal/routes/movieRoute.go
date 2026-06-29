package routes

import (
	"github.com/gofiber/fiber/v3"
	"final-project/internal/handlers"
	"final-project/internal/services"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	
	MovieRoutes(api)
}

func MovieRoutes(router fiber.Router) {
	movieService := services.NewMovieService()
	movieHandler := handlers.NewMovieHandler(movieService)

	router.Get("/movies", movieHandler.GetAllMovies)            
	router.Get("/movies/search", movieHandler.GetMoviesFilter)   
	router.Get("/movies/page", movieHandler.GetMoviesPaginated)  
	router.Get("/movies/stats", movieHandler.GetMovieStats)     
	
	// Классический CRUD по ID
	router.Get("/movies/:id", movieHandler.GetMovieByID)    
	router.Post("/movies", movieHandler.CreateMovie)        
	router.Put("/movies/:id", movieHandler.UpdateMovie)     
	router.Patch("/movies/:id", movieHandler.PatchMovie)   
	router.Delete("/movies/:id", movieHandler.DeleteMovie)  
}