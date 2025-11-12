package router

import (
	"encoding/json"
	"net/http"
	"taskforge/internal/auth"
	"taskforge/internal/handlers"

	// "taskforge/internal/middleware"
	// "taskforge/internal/middleware"
	"github.com/rs/cors"
)

type Router struct {
	mux         *http.ServeMux
	userHandler *handlers.UserHandler
	taskHandler *handlers.TaskHandler
	jwtManager  *auth.JWTManager
}

func NewRouter(userHandler *handlers.UserHandler, taskHandler *handlers.TaskHandler, jwtManager *auth.JWTManager) *Router {
	return &Router{
		mux:         http.NewServeMux(),
		userHandler: userHandler,
		taskHandler: taskHandler,
		jwtManager:  jwtManager,
	}
}

// func (r *Router) SetupRoutes() *http.ServeMux {
// 	r.mux.HandleFunc("POST /api/v1/register", r.userHandler.Register)
// 	r.mux.HandleFunc("POST /api/v1/login", r.userHandler.Login)

// 	r.mux.Handle("GET /api/v1/protected", middleware.AuthMiddleware(jwtManager)(
// 		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			// Получаем user из context'а
// 			claims := middleware.GetUserFromContext(r.Context())
// 			if claims == nil {
// 				JSONError(w, http.StatusUnauthorized, "User not found")
// 				return
// 			}

// 			// Используем данные пользователя
// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(map[string]interface{}{
// 				"success": true,
// 				"data": map[string]interface{}{
// 					"user_id": claims.UserID,
// 					"email":   claims.Email,
// 					"message": "This is a protected endpoint",
// 				},
// 			})
// 		}),
// 	))

// 	return r.mux
// }

func (r *Router) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Public routes
	mux.HandleFunc("POST /api/v1/register", r.userHandler.Register)
	mux.HandleFunc("POST /api/v1/login", r.userHandler.Login)
	mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		JSONResponse(w, http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"message": "Service is running",
		})
	})

	// Protected routes - TASKS
	mux.Handle("GET /api/v1/tasks", r.jwtManager.AuthMiddleware(
		http.HandlerFunc(r.taskHandler.GetTasks),
	))
	mux.Handle("POST /api/v1/tasks", r.jwtManager.AuthMiddleware(
		http.HandlerFunc(r.taskHandler.CreateTask),
	))
	mux.Handle("GET /api/v1/tasks/{id}", r.jwtManager.AuthMiddleware(
		http.HandlerFunc(r.taskHandler.GetTaskByID),
	))
	mux.Handle("PUT /api/v1/tasks/{id}", r.jwtManager.AuthMiddleware(
		http.HandlerFunc(r.taskHandler.UpdateTask),
	))
	mux.Handle("DELETE /api/v1/tasks/{id}", r.jwtManager.AuthMiddleware(
		http.HandlerFunc(r.taskHandler.DeleteTask),
	))

	// Добавляем CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := c.Handler(mux)
	return handler
}

// Обработчик для защищенного маршрута
// func (r *Router) protectedHandler(w http.ResponseWriter, req *http.Request) {
// 	// Получаем user из context'а

// 	claims, ok := auth.GetUserFromContext(req.Context())
// 	if !ok || claims == nil {
// 		JSONError(w, http.StatusUnauthorized, "User not found in context")
// 		return
// 	}

// 	// Используем данные пользователя
// 	JSONResponse(w, http.StatusOK, map[string]interface{}{
// 		"success": true,
// 		"data": map[string]interface{}{
// 			"user_id": claims.UserID,
// 			"email":   claims.Email,
// 			"message": "This is a protected endpoint",
// 		},
// 	})
// }

func JSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func JSONError(w http.ResponseWriter, status int, message string) {
	JSONResponse(w, status, map[string]interface{}{
		"error":   true,
		"message": message,
	})
}
