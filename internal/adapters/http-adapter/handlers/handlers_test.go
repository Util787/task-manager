package handlers

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Util787/task-manager/internal/domain"
	"github.com/Util787/task-manager/internal/infrastructure/repo/inmemory"
	"github.com/Util787/task-manager/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

const timeDelta = time.Second

func setupTestRouter(handlers *Handlers) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// middleware
	router.Use(func(c *gin.Context) {
		c.Set("op", "test")
		c.Next()
	})

	// routes init
	router.POST("/tasks", handlers.createTask)
	router.GET("/tasks/:id/state", handlers.getTaskStateByID)
	router.GET("/tasks/:id/result", handlers.getTaskResultByID)
	router.DELETE("/tasks/:id", handlers.deleteTask)

	return router
}

func createTestHandlers() (*Handlers, *inmemory.TaskRepository) {
	logger := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	repo := inmemory.NewTaskRepository(logger)
	usecase := usecase.NewTaskUsecase(repo)
	handlers := New(logger, usecase)
	return handlers, repo
}

// create task tests

func TestCreateTask_OK(t *testing.T) {
	handlers, _ := createTestHandlers()
	router := setupTestRouter(handlers)

	// test data
	requestBody := createTaskRequest{
		Title:       "Test Task",
		Description: "Test Description",
	}

	// request
	jsonBody, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// response check
	assert.Equal(t, http.StatusCreated, w.Code)

	var response createTaskResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "task created successfully with id")

	message := response.Message
	idStr := message[len("task created successfully with id "):]
	_, err = uuid.Parse(idStr)
	assert.NoError(t, err)
}

func TestCreateTask_InvalidJSON(t *testing.T) {
	handlers, _ := createTestHandlers()
	router := setupTestRouter(handlers)

	// request
	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// response check
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "invalid request body")
}

func TestCreateTask_MissingTitle(t *testing.T) {
	handlers, _ := createTestHandlers()
	router := setupTestRouter(handlers)

	// request
	requestBody := map[string]string{
		"description": "Test Description",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// response check
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "invalid request body")
}

func TestCreateTask_EmptyTitle(t *testing.T) {
	handlers, _ := createTestHandlers()
	router := setupTestRouter(handlers)

	// request
	requestBody := createTaskRequest{
		Title:       "",
		Description: "Test Description",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// response check
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "invalid request body")
}

func TestCreateTask_TitleTooLong(t *testing.T) {
	handlers, _ := createTestHandlers()
	router := setupTestRouter(handlers)

	// request
	longTitle := ""
	for i := 0; i < 256; i++ {
		longTitle += "a"
	}

	requestBody := createTaskRequest{
		Title:       longTitle,
		Description: "Test Description",
	}
	jsonBody, _ := json.Marshal(requestBody)

	req, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// response check
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response.Message, "invalid request body")
}

// get task state tests

func TestGetTaskStateByID_OK(t *testing.T) {
	handlers, repo := createTestHandlers()
	router := setupTestRouter(handlers)

	// create task
	task := &domain.Task{
		Title:       "Test Task",
		Description: "Test Description",
		TaskState: domain.TaskState{
			Status:       domain.StatusCompleted,
			WorkDuration: time.Second * 10,
		},
	}

	taskID := repo.CreateTask(task)

	// request
	req, _ := http.NewRequest("GET", "/tasks/"+taskID.String()+"/state", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// response check
	assert.Equal(t, http.StatusOK, w.Code)

	var response getTaskStateResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, task.TaskState, response.State)
	assert.WithinDuration(t, task.CreatedAt, response.CreatedAt, timeDelta)
}

func TestGetTaskStateByID_InvalidUUID(t *testing.T) {
	handlers, _ := createTestHandlers()
	router := setupTestRouter(handlers)

	// request
	req, _ := http.NewRequest("GET", "/tasks/invalid-uuid/state", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// response check
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid task id", response.Message)
}

func TestGetTaskStateByID_TaskNotFound(t *testing.T) {
	handlers, _ := createTestHandlers()
	router := setupTestRouter(handlers)

	// request
	nonExistentID := uuid.New()
	req, _ := http.NewRequest("GET", "/tasks/"+nonExistentID.String()+"/state", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// response check
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "task not found", response.Message)
}

// get task result tests

func TestGetTaskResultByID_OK(t *testing.T) {
	handlers, repo := createTestHandlers()
	router := setupTestRouter(handlers)

	// create task
	task := &domain.Task{
		Title:       "Test Task",
		Description: "Test Description",
		Result:      "task completed successfully",
	}

	taskID := repo.CreateTask(task)

	// request
	req, _ := http.NewRequest("GET", "/tasks/"+taskID.String()+"/result", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// response check
	assert.Equal(t, http.StatusOK, w.Code)

	var response getTaskResultResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "task result: "+task.Result, response.Message)
}

func TestGetTaskResultByID_InvalidUUID(t *testing.T) {
	handlers, _ := createTestHandlers()
	router := setupTestRouter(handlers)

	// request
	req, _ := http.NewRequest("GET", "/tasks/invalid-uuid/result", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// response check
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid task id", response.Message)
}

func TestGetTaskResultByID_TaskNotFound(t *testing.T) {
	handlers, _ := createTestHandlers()
	router := setupTestRouter(handlers)

	// request
	nonExistentID := uuid.New()
	req, _ := http.NewRequest("GET", "/tasks/"+nonExistentID.String()+"/result", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// response check
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "task not found", response.Message)
}

// delete task tests

func TestDeleteTask_OK(t *testing.T) {
	handlers, repo := createTestHandlers()
	router := setupTestRouter(handlers)

	// create task
	task := &domain.Task{
		Title:       "Test Task",
		Description: "Test Description",
	}

	taskID := repo.CreateTask(task)

	// request
	req, _ := http.NewRequest("DELETE", "/tasks/"+taskID.String(), nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// response check
	assert.Equal(t, http.StatusOK, w.Code)

	var response deleteTaskResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "task deleted successfully", response.Message)

	// check if task is deleted
	_, _, err = repo.GetTaskStateByID(taskID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "task not found")
}

func TestDeleteTask_InvalidUUID(t *testing.T) {
	handlers, _ := createTestHandlers()
	router := setupTestRouter(handlers)

	// request
	req, _ := http.NewRequest("DELETE", "/tasks/invalid-uuid", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// response check
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "invalid task id", response.Message)
}

func TestDeleteTask_TaskNotFound(t *testing.T) {
	handlers, _ := createTestHandlers()
	router := setupTestRouter(handlers)

	// request
	nonExistentID := uuid.New()
	req, _ := http.NewRequest("DELETE", "/tasks/"+nonExistentID.String(), nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// response check
	assert.Equal(t, http.StatusNotFound, w.Code)

	var response errorResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "task not found", response.Message)
}

// Integration test
func TestTaskFullLifecycle(t *testing.T) {
	handlers, _ := createTestHandlers()
	router := setupTestRouter(handlers)

	// 1. create task
	createRequestBody := createTaskRequest{
		Title:       "Integration Test Task",
		Description: "This is a test task for integration testing",
	}

	createJSON, _ := json.Marshal(createRequestBody)
	createReq, _ := http.NewRequest("POST", "/tasks", bytes.NewBuffer(createJSON))
	createReq.Header.Set("Content-Type", "application/json")

	createW := httptest.NewRecorder()
	router.ServeHTTP(createW, createReq)

	assert.Equal(t, http.StatusCreated, createW.Code)

	var createResponse createTaskResponse
	err := json.Unmarshal(createW.Body.Bytes(), &createResponse)
	assert.NoError(t, err)

	// extract id from response
	message := createResponse.Message
	idStr := message[len("task created successfully with id "):]
	taskID, err := uuid.Parse(idStr)
	assert.NoError(t, err)

	// 2. get task state
	stateReq, _ := http.NewRequest("GET", "/tasks/"+taskID.String()+"/state", nil)
	stateW := httptest.NewRecorder()
	router.ServeHTTP(stateW, stateReq)

	assert.Equal(t, http.StatusOK, stateW.Code)

	var stateResponse getTaskStateResponse
	err = json.Unmarshal(stateW.Body.Bytes(), &stateResponse)
	assert.NoError(t, err)
	assert.Equal(t, domain.StatusInProgress, stateResponse.State.Status)

	// 3. get task result(empty, because task is just created)
	resultReq, _ := http.NewRequest("GET", "/tasks/"+taskID.String()+"/result", nil)
	resultW := httptest.NewRecorder()
	router.ServeHTTP(resultW, resultReq)

	assert.Equal(t, http.StatusOK, resultW.Code)

	var resultResponse getTaskResultResponse
	err = json.Unmarshal(resultW.Body.Bytes(), &resultResponse)
	assert.NoError(t, err)
	assert.Equal(t, "task result: ", resultResponse.Message)

	// 4. delete task
	deleteReq, _ := http.NewRequest("DELETE", "/tasks/"+taskID.String(), nil)
	deleteW := httptest.NewRecorder()
	router.ServeHTTP(deleteW, deleteReq)

	assert.Equal(t, http.StatusOK, deleteW.Code)

	var deleteResponse deleteTaskResponse
	err = json.Unmarshal(deleteW.Body.Bytes(), &deleteResponse)
	assert.NoError(t, err)
	assert.Equal(t, "task deleted successfully", deleteResponse.Message)

	// 5. check if task is deleted
	stateReq2, _ := http.NewRequest("GET", "/tasks/"+taskID.String()+"/state", nil)
	stateW2 := httptest.NewRecorder()
	router.ServeHTTP(stateW2, stateReq2)

	assert.Equal(t, http.StatusNotFound, stateW2.Code)
}
