package channel

import (
	"fmt"
	"sync/atomic"
	"tasks-service-demo/internal/entities"
	apperrors "tasks-service-demo/internal/errors"
)

// Operation types
const (
	OpCreate   = "create"
	OpRead     = "read"
	OpUpdate   = "update"
	OpDelete   = "delete"
	OpGetAll   = "getall"
	OpShutdown = "shutdown"
)

// Operation represents a request to the channel store
type Operation struct {
	Type     string
	TaskID   int
	Task     *entities.Task
	Response chan Result
}

// Result represents the response from an operation
type Result struct {
	Task  *entities.Task
	Tasks []*entities.Task
	Error error
}

// ChannelStore implements simple single-worker channel-based storage
type ChannelStore struct {
	operations chan Operation
	nextID     int64 // atomic counter for ID generation
	shutdown   chan struct{}
}

// NewChannelStore creates a simple single-worker channel-based store
func NewChannelStore(numWorkers int) *ChannelStore {
	cs := &ChannelStore{
		operations: make(chan Operation, 1000),
		nextID:     0,
		shutdown:   make(chan struct{}),
	}

	// Start single worker
	go cs.worker()

	return cs
}

// worker processes operations from the channel (simple single-worker)
func (cs *ChannelStore) worker() {
	// Single worker with local storage - no locks needed!
	localStorage := make(map[int]*entities.Task)

	for {
		select {
		case op := <-cs.operations:
			switch op.Type {
			case OpCreate:
				localStorage[op.Task.ID] = op.Task
				op.Response <- Result{Task: op.Task, Error: nil}

			case OpRead:
				if task, exists := localStorage[op.TaskID]; exists {
					// Return a copy to avoid race conditions
					taskCopy := *task
					op.Response <- Result{Task: &taskCopy, Error: nil}
				} else {
					op.Response <- Result{Error: apperrors.ErrTaskNotFound}
				}

			case OpUpdate:
				if _, exists := localStorage[op.TaskID]; exists {
					op.Task.ID = op.TaskID
					localStorage[op.TaskID] = op.Task
					op.Response <- Result{Task: op.Task, Error: nil}
				} else {
					op.Response <- Result{Error: apperrors.ErrTaskNotFound}
				}

			case OpDelete:
				if _, exists := localStorage[op.TaskID]; exists {
					delete(localStorage, op.TaskID)
					op.Response <- Result{Error: nil}
				} else {
					op.Response <- Result{Error: apperrors.ErrTaskNotFound}
				}

			case OpGetAll:
				// Collect all tasks from this worker's storage
				tasks := make([]*entities.Task, 0, len(localStorage))
				for _, task := range localStorage {
					taskCopy := *task
					tasks = append(tasks, &taskCopy)
				}
				op.Response <- Result{Tasks: tasks, Error: nil}

			case OpShutdown:
				return
			}

		case <-cs.shutdown:
			return
		}
	}
}

// Create adds a new task to the store
func (cs *ChannelStore) Create(task *entities.Task) *apperrors.AppError {
	// Generate unique ID atomically
	id := int(atomic.AddInt64(&cs.nextID, 1))
	task.ID = id

	response := make(chan Result, 1)

	op := Operation{
		Type:     OpCreate,
		Task:     task,
		Response: response,
	}

	cs.operations <- op
	result := <-response
	if result.Error != nil {
		return apperrors.ErrStorageError.WithCause(fmt.Errorf("Create failed from channel result: %v", result.Error))
	}
	return nil
}

// GetByID retrieves a task by its ID
func (cs *ChannelStore) GetByID(id int) (*entities.Task, *apperrors.AppError) {
	response := make(chan Result, 1)

	op := Operation{
		Type:     OpRead,
		TaskID:   id,
		Response: response,
	}

	cs.operations <- op
	result := <-response
	if result.Error != nil {
		return nil, apperrors.ErrStorageError.WithCause(fmt.Errorf("GetByID failed from channel result: %v", result.Error))
	}

	return result.Task, nil
}

// Update modifies an existing task
func (cs *ChannelStore) Update(id int, updatedTask *entities.Task) *apperrors.AppError {
	response := make(chan Result, 1)

	op := Operation{
		Type:     OpUpdate,
		TaskID:   id,
		Task:     updatedTask,
		Response: response,
	}

	cs.operations <- op
	result := <-response
	if result.Error != nil {
		return apperrors.ErrStorageError.WithCause(fmt.Errorf("Update failed from channel result: %v", result.Error))
	}
	return nil
}

// Delete removes a task from the store
func (cs *ChannelStore) Delete(id int) *apperrors.AppError {
	response := make(chan Result, 1)

	op := Operation{
		Type:     OpDelete,
		TaskID:   id,
		Response: response,
	}

	cs.operations <- op
	result := <-response
	if result.Error != nil {
		return apperrors.ErrStorageError.WithCause(fmt.Errorf("Delete failed from channel result: %v", result.Error))
	}
	return nil
}

// GetAll retrieves all tasks
func (cs *ChannelStore) GetAll() []*entities.Task {
	response := make(chan Result, 1)

	op := Operation{
		Type:     OpGetAll,
		Response: response,
	}

	cs.operations <- op
	result := <-response

	if result.Error != nil {
		return []*entities.Task{}
	}

	return result.Tasks
}

// Shutdown gracefully shuts down the storage manager
func (cs *ChannelStore) Shutdown() {
	close(cs.shutdown)
}
