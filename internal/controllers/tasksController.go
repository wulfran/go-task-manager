package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"task-manager/internal/contextkeys"
	"task-manager/internal/helpers"
	"task-manager/internal/models"
	"task-manager/internal/requests"
	"task-manager/internal/services"
	"time"

	"github.com/go-chi/chi/v5"
)

type TasksController interface {
	Index() func(w http.ResponseWriter, r *http.Request)
	Store(bodySizeLimit int64) func(w http.ResponseWriter, r *http.Request)
	Show() func(w http.ResponseWriter, r *http.Request)
	Update(bodySizeLimit int64) func(w http.ResponseWriter, r *http.Request)
	Delete() func(w http.ResponseWriter, r *http.Request)
}

type tasksController struct {
	ts services.TaskService
}

func NewTasksController(ts services.TaskService) TasksController {
	return &tasksController{ts: ts}
}

func (t tasksController) Index() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		uID, ok := r.Context().Value(contextkeys.UserID).(int64)
		if !ok {
			helpers.JsonResponse(w, http.StatusInternalServerError, fmt.Sprintf("invalid user data, please relog"))
			return
		}

		tl, err := t.ts.GetTasksList(uID)
		if err != nil {
			helpers.JsonResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to retreive tasks list: %v", err))
			return
		}

		helpers.JsonResponse(w, http.StatusOK, tl)
	}
}
func (t tasksController) Store(bodySizeLimit int64) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, bodySizeLimit)
		var req requests.CreateTasksRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			helpers.JsonResponse(w, http.StatusRequestEntityTooLarge, fmt.Sprintf("store task: request body too large: %v", err))
			return
		}

		v := req.Validate()
		if !v.Validated {
			helpers.JsonResponse(w, http.StatusUnprocessableEntity, fmt.Sprintf("store task: validation failed: %s", v.Message))
			return
		}
		now := time.Now()
		p := models.TaskPayload{
			Name:        req.Name,
			Priority:    req.Priority,
			Description: req.Description,
			DueDate:     req.DueDate,
			CreatedAt:   &now,
		}
		if err := t.ts.StoreTask(r.Context(), p); err != nil {
			helpers.JsonResponse(w, http.StatusInternalServerError, fmt.Sprintf("store task: failed to save the data: %v", err))
			return
		}
		helpers.JsonResponse(w, http.StatusOK, fmt.Sprintf("successfully created a new task"))
	}
}
func (t tasksController) Show() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		uID, ok := r.Context().Value(contextkeys.UserID).(int64)
		if !ok {
			helpers.JsonResponse(w, http.StatusInternalServerError, fmt.Sprintf("invalid user data, please relog"))
			return
		}

		rawID := chi.URLParam(r, "task_id")
		var id int
		id, err := strconv.Atoi(rawID)
		if err != nil {
			helpers.JsonResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid task id"))
		}

		t, err := t.ts.ShowTask(id)
		if err != nil {
			helpers.JsonResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to get task data: %v", err))
			return
		}

		if t.CreatedBy != uID {
			helpers.JsonResponse(w, http.StatusUnauthorized, fmt.Sprintf("you do not have the permission to access this data"))
			return
		}

		helpers.JsonResponse(w, http.StatusOK, t)
	}
}
func (t tasksController) Update(bodySizeLimit int64) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, bodySizeLimit)
		var req requests.UpdateTaskRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			helpers.JsonResponse(w, http.StatusRequestEntityTooLarge, fmt.Sprintf("update task: payload invalid: %v", err))
			return
		}

		v := req.Validate()
		if !v.Validated {
			helpers.JsonResponse(w, http.StatusUnprocessableEntity, fmt.Sprintf("update task: validation failed: %s", v.Message))
			return
		}

		rawId := chi.URLParam(r, "task_id")
		var id int64
		id, err = strconv.ParseInt(rawId, 10, 10)
		if err != nil {
			helpers.JsonResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid task id"))
			return
		}

		uID, ok := r.Context().Value(contextkeys.UserID).(int64)
		if !ok {
			helpers.JsonResponse(w, http.StatusInternalServerError, fmt.Sprintf("invalid user data, please relog"))
			return
		}

		isOwner, err := t.ts.IsTaskOwner(uID, int(id))
		if err != nil {
			helpers.JsonResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to check ownership: %v", err))
			return
		}

		if !isOwner {
			helpers.JsonResponse(w, http.StatusUnauthorized, fmt.Sprintf("you do not have the permissions for that action"))
			return
		}

		p := models.UpdateTask{
			ID:          id,
			Name:        req.Name,
			Priority:    req.Priority,
			Description: req.Description,
			DueDate:     req.DueDate,
		}

		t, err := t.ts.UpdateTask(r.Context(), p)
		if err != nil {
			helpers.JsonResponse(w, http.StatusInternalServerError, fmt.Sprintf("update: failed to update task: %v", err))
			return
		}

		helpers.JsonResponse(w, http.StatusOK, t)
	}
}
func (t tasksController) Delete() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		uID, ok := r.Context().Value(contextkeys.UserID).(int64)
		if !ok {
			helpers.JsonResponse(w, http.StatusInternalServerError, fmt.Sprintf("invalid user data please relog"))
			return
		}

		rawID := chi.URLParam(r, "task_id")
		var id int
		id, err := strconv.Atoi(rawID)
		if err != nil {
			helpers.JsonResponse(w, http.StatusBadRequest, fmt.Sprintf("invalid task id"))
			return
		}

		if err := t.ts.DeleteTask(id, uID); err != nil {
			helpers.JsonResponse(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete the task: %v", err))
			return
		}
		helpers.JsonResponse(w, http.StatusOK, fmt.Sprintf("task deleted successfully"))
	}
}
