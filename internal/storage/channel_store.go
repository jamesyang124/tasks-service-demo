package storage

import (
	"errors"
	"sync/atomic"
	"tasks-service-demo/internal/models"
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
	Task     *models.Task
	Response chan Result
}

// Result represents the response from an operation
type Result struct {
	Task  *models.Task
	Tasks []*models.Task
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
	localStorage := make(map[int]*models.Task)

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
					op.Response <- Result{Error: errors.New("task not found")}
				}

			case OpUpdate:
				if _, exists := localStorage[op.TaskID]; exists {
					op.Task.ID = op.TaskID
					localStorage[op.TaskID] = op.Task
					op.Response <- Result{Task: op.Task, Error: nil}
				} else {
					op.Response <- Result{Error: errors.New("task not found")}
				}

			case OpDelete:
				if _, exists := localStorage[op.TaskID]; exists {
					delete(localStorage, op.TaskID)
					op.Response <- Result{Error: nil}
				} else {
					op.Response <- Result{Error: errors.New("task not found")}
				}

			case OpGetAll:
				// Collect all tasks from this worker's storage
				tasks := make([]*models.Task, 0, len(localStorage))
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
func (cs *ChannelStore) Create(task *models.Task) error {
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
	return result.Error
}

// GetByID retrieves a task by its ID
func (cs *ChannelStore) GetByID(id int) (*models.Task, error) {
	response := make(chan Result, 1)

	op := Operation{
		Type:     OpRead,
		TaskID:   id,
		Response: response,
	}

	cs.operations <- op
	result := <-response
	return result.Task, result.Error
}

// Update modifies an existing task
func (cs *ChannelStore) Update(id int, updatedTask *models.Task) error {
	response := make(chan Result, 1)

	op := Operation{
		Type:     OpUpdate,
		TaskID:   id,
		Task:     updatedTask,
		Response: response,
	}

	cs.operations <- op
	result := <-response
	return result.Error
}

// Delete removes a task from the store
func (cs *ChannelStore) Delete(id int) error {
	response := make(chan Result, 1)

	op := Operation{
		Type:     OpDelete,
		TaskID:   id,
		Response: response,
	}

	cs.operations <- op
	result := <-response
	return result.Error
}

// GetAll retrieves all tasks
func (cs *ChannelStore) GetAll() []*models.Task {
	response := make(chan Result, 1)
	
	op := Operation{
		Type:     OpGetAll,
		Response: response,
	}
	
	cs.operations <- op
	result := <-response
	
	if result.Error != nil {
		return []*models.Task{}
	}
	
	return result.Tasks
}

// Shutdown gracefully shuts down the storage manager
func (cs *ChannelStore) Shutdown() {
	close(cs.shutdown)
}
