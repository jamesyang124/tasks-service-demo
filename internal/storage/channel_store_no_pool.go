package storage

import (
	"errors"
	"sync/atomic"
	"tasks-service-demo/internal/models"
)

// ChannelStoreNoPool implements the same logic but without sync.Pool for comparison
type ChannelStoreNoPool struct {
	workers    []chan Operation
	numWorkers int
	nextID     int64 // atomic counter for ID generation
	shutdown   chan struct{}
}

// NewChannelStoreNoPool creates a new channel-based store without sync.Pool
func NewChannelStoreNoPool(numWorkers int) *ChannelStoreNoPool {
	if numWorkers <= 0 {
		numWorkers = 4 // Default worker pool size
	}

	cs := &ChannelStoreNoPool{
		workers:    make([]chan Operation, numWorkers),
		numWorkers: numWorkers,
		nextID:     0,
		shutdown:   make(chan struct{}),
	}

	// Start worker pool - each worker has its own channel and storage
	for i := 0; i < numWorkers; i++ {
		cs.workers[i] = make(chan Operation, 1000) // Much larger buffers
		go cs.worker(i, cs.workers[i])
	}

	return cs
}

// worker processes operations from its own channel (true lock-free per worker)
func (cs *ChannelStoreNoPool) worker(workerID int, operations chan Operation) {
	// Each worker has its own local storage - no locks needed!
	localStorage := make(map[int]*models.Task)

	for {
		select {
		case op := <-operations:
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

// getWorker determines which worker should handle the operation (hash-based routing)
func (cs *ChannelStoreNoPool) getWorker(id int) int {
	return id % cs.numWorkers
}

// Create adds a new task to the store
func (cs *ChannelStoreNoPool) Create(task *models.Task) error {
	// Generate unique ID atomically
	id := int(atomic.AddInt64(&cs.nextID, 1))
	task.ID = id

	// Route to appropriate worker
	workerIndex := cs.getWorker(id)
	response := make(chan Result, 1) // No pooling - fresh allocation

	op := Operation{
		Type:     OpCreate,
		Task:     task,
		Response: response,
	}

	cs.workers[workerIndex] <- op
	result := <-response
	return result.Error
}

// GetByID retrieves a task by its ID
func (cs *ChannelStoreNoPool) GetByID(id int) (*models.Task, error) {
	workerIndex := cs.getWorker(id)
	response := make(chan Result, 1) // No pooling - fresh allocation

	op := Operation{
		Type:     OpRead,
		TaskID:   id,
		Response: response,
	}

	cs.workers[workerIndex] <- op
	result := <-response
	return result.Task, result.Error
}

// Update modifies an existing task
func (cs *ChannelStoreNoPool) Update(id int, updatedTask *models.Task) error {
	workerIndex := cs.getWorker(id)
	response := make(chan Result, 1) // No pooling - fresh allocation

	op := Operation{
		Type:     OpUpdate,
		TaskID:   id,
		Task:     updatedTask,
		Response: response,
	}

	cs.workers[workerIndex] <- op
	result := <-response
	return result.Error
}

// Delete removes a task from the store
func (cs *ChannelStoreNoPool) Delete(id int) error {
	workerIndex := cs.getWorker(id)
	response := make(chan Result, 1) // No pooling - fresh allocation

	op := Operation{
		Type:     OpDelete,
		TaskID:   id,
		Response: response,
	}

	cs.workers[workerIndex] <- op
	result := <-response
	return result.Error
}

// GetAll retrieves all tasks from all workers
func (cs *ChannelStoreNoPool) GetAll() []*models.Task {
	responses := make([]chan Result, cs.numWorkers)
	
	// Send GetAll operation to all workers
	for i := 0; i < cs.numWorkers; i++ {
		responses[i] = make(chan Result, 1) // No pooling - fresh allocation
		op := Operation{
			Type:     OpGetAll,
			Response: responses[i],
		}
		cs.workers[i] <- op
	}

	// Collect results from all workers
	var allTasks []*models.Task
	for i := 0; i < cs.numWorkers; i++ {
		result := <-responses[i]
		if result.Error == nil {
			allTasks = append(allTasks, result.Tasks...)
		}
	}

	return allTasks
}

// Shutdown gracefully shuts down the storage manager
func (cs *ChannelStoreNoPool) Shutdown() {
	close(cs.shutdown)
}