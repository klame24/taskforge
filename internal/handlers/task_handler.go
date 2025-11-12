package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"taskforge/internal/auth"
	"taskforge/internal/dto/request"
	"taskforge/internal/models"
	"taskforge/internal/repository"
)

func extractTaskIDFromURL(r *http.Request) (int, error) {
	// Извлекаем ID из URL пути: /api/v1/tasks/123
	path := r.URL.Path
	segments := strings.Split(path, "/")

	if len(segments) < 5 {
		return 0, errors.New("invalid URL format")
	}

	taskIDStr := segments[4] // /api/v1/tasks/123 → segments[4] = "123"
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		return 0, errors.New("task ID must be a number")
	}

	return taskID, nil
}

type TaskHandler struct {
	taskRepo repository.TaskRepository
}

func NewTaskHandler(taskRepo repository.TaskRepository) *TaskHandler {
	return &TaskHandler{
		taskRepo: taskRepo,
	}
}

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	userID, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		JSONError(w, http.StatusUnauthorized, "User authentication required")
		return
	}

	// maybe problem with userID.UserID
	tasks, err := h.taskRepo.GetByUserID(userID.UserID)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to get tasks")
		return
	}

	JSONSuccess(w, http.StatusOK, map[string]interface{}{
		"tasks": tasks,
	})
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		JSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	defer r.Body.Close()

	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		JSONError(w, http.StatusUnauthorized, "user authentication required")
		return
	}

	var req request.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if req.Title == "" {
		JSONError(w, http.StatusBadRequest, "Title is required")
		return
	}

	if req.Status == "" {
		req.Status = "todo"
	}
	if req.Priority == "" {
		req.Priority = "medium"
	}

	task := &models.Task{
		UserID:      claims.UserID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}

	if err := h.taskRepo.Create(task); err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to create task")
		return
	}

	JSONSuccess(w, http.StatusCreated, map[string]interface{}{
		"task": task,
	})
}

func (h *TaskHandler) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		JSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		JSONError(w, http.StatusUnauthorized, "user authentication required")
		return
	}

	taskID, err := extractTaskIDFromURL(r)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid Task ID")
		return
	}

	task, err := h.taskRepo.GetByID(taskID)
	if err != nil {
		if err == repository.ErrTaskNotFound {
			JSONError(w, http.StatusNotFound, "Task not found")
			return
		} else {
			JSONError(w, http.StatusInternalServerError, "Failed to get task")
		}
		return
	}

	if task.UserID != claims.UserID {
		JSONError(w, http.StatusForbidden, "Access denied")
		return
	}

	JSONSuccess(w, http.StatusOK, map[string]interface{}{
		"task": task,
	})
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		JSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	defer r.Body.Close()

	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		JSONError(w, http.StatusUnauthorized, "authenticate user required")
		return
	}

	taskID, err := extractTaskIDFromURL(r)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid Task ID")
		return
	}

	var req request.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	task, err := h.taskRepo.GetByID(taskID)
	if err != nil {
		if err == repository.ErrTaskNotFound {
			JSONError(w, http.StatusNotFound, "Task not found")
			return
		} else {
			JSONError(w, http.StatusInternalServerError, "Failed to get task")
		}
		return
	}

	if task.UserID != claims.UserID {
		JSONError(w, http.StatusForbidden, "Access denied")
		return
	}

	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Status != "" {
		task.Status = req.Status
	}
	if req.Priority != "" {
		task.Priority = req.Priority
	}
	if !req.DueDate.IsZero() {
		task.DueDate = req.DueDate
	}

	if err := h.taskRepo.Update(task); err != nil {
		JSONError(w, http.StatusInternalServerError, "Failed to update task")
		return
	}

	JSONSuccess(w, http.StatusOK, map[string]interface{}{
		"task": task,
	})
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		JSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	claims, ok := auth.GetUserFromContext(r.Context())
	if !ok {
		JSONError(w, http.StatusUnauthorized, "user authenticate required")
		return
	}

	taskID, err := extractTaskIDFromURL(r)
	if err != nil {
		JSONError(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := h.taskRepo.GetByID(taskID)
	if err != nil {
		if err == repository.ErrTaskNotFound {
			JSONError(w, http.StatusNotFound, "Task not found")
			return
		} else {
			JSONError(w, http.StatusInternalServerError, "Cant get task")
			return
		}
	}

	if task.UserID != claims.UserID {
		JSONError(w, http.StatusForbidden, "Access denied")
		return
	}

	if err := h.taskRepo.Delete(taskID); err != nil {
		JSONError(w, http.StatusInternalServerError, "Cant delete this task")
		return
	}

	JSONSuccess(w, http.StatusNoContent, map[string]interface{}{
		"message": "Task deleted successfully",
	})
}
