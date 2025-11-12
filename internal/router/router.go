package router

import (
	"encoding/json"
	"net/http"
	"taskforge/internal/auth"
	"taskforge/internal/handlers"
	// "taskforge/internal/middleware"
	// "taskforge/internal/middleware"
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

func (r *Router) SetupRoutes() *http.ServeMux {
	// Public routes (доступны без авторизации)
	r.mux.HandleFunc("POST /api/v1/register", r.userHandler.Register)
	r.mux.HandleFunc("POST /api/v1/login", r.userHandler.Login)

	// Health check
	r.mux.HandleFunc("GET /api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		JSONResponse(w, http.StatusOK, map[string]interface{}{
			"status":  "ok",
			"message": "Service is running",
		})
	})

	// Protected routes - TASKS
    r.mux.Handle("GET /api/v1/tasks", r.jwtManager.AuthMiddleware(
        http.HandlerFunc(r.taskHandler.GetTasks),
    ))
    
    r.mux.Handle("POST /api/v1/tasks", r.jwtManager.AuthMiddleware(
        http.HandlerFunc(r.taskHandler.CreateTask),
    ))
    
    r.mux.Handle("GET /api/v1/tasks/{id}", r.jwtManager.AuthMiddleware(
        http.HandlerFunc(r.taskHandler.GetTaskByID),
    ))
    
    r.mux.Handle("PUT /api/v1/tasks/{id}", r.jwtManager.AuthMiddleware(
        http.HandlerFunc(r.taskHandler.UpdateTask),
    ))
    
    r.mux.Handle("DELETE /api/v1/tasks/{id}", r.jwtManager.AuthMiddleware(
        http.HandlerFunc(r.taskHandler.DeleteTask),
    ))

	// Protected routes (требуют авторизации)
	protected := http.NewServeMux()

	// Пример защищенного маршрута
	protected.Handle("GET /api/v1/protected", r.jwtManager.AuthMiddleware(
		http.HandlerFunc(r.protectedHandler),
	))

	// Подключаем защищенные маршруты к основному роутеру
	r.mux.Handle("/api/v1/", protected)

	return r.mux
}

// Обработчик для защищенного маршрута
func (r *Router) protectedHandler(w http.ResponseWriter, req *http.Request) {
	// Получаем user из context'а

	claims, ok := auth.GetUserFromContext(req.Context())
	if !ok || claims == nil {
		JSONError(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	// Используем данные пользователя
	JSONResponse(w, http.StatusOK, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"user_id": claims.UserID,
			"email":   claims.Email,
			"message": "This is a protected endpoint",
		},
	})
}

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
