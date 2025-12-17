package health

import (
	"context"
	"sync"
	"time"

	"gorm.io/gorm"
)

// Status represents the health status of a component
type Status string

const (
	StatusUp      Status = "up"
	StatusDown    Status = "down"
	StatusUnknown Status = "unknown"
)

// Check represents a single health check result
type Check struct {
	Name    string `json:"name"`
	Status  Status `json:"status"`
	Message string `json:"message,omitempty"`
	Latency string `json:"latency,omitempty"`
}

// Response represents the overall health response
type Response struct {
	Status    Status            `json:"status"`
	Service   string            `json:"service"`
	Version   string            `json:"version"`
	Timestamp string            `json:"timestamp"`
	Checks    []Check           `json:"checks"`
	Details   map[string]string `json:"details,omitempty"`
}

// Checker provides health check functionality
type Checker struct {
	serviceName string
	version     string
	db          *gorm.DB
	checks      []func(context.Context) Check
	mu          sync.RWMutex
}

// NewChecker creates a new health checker
func NewChecker(serviceName, version string, db *gorm.DB) *Checker {
	return &Checker{
		serviceName: serviceName,
		version:     version,
		db:          db,
		checks:      make([]func(context.Context) Check, 0),
	}
}

// AddCheck adds a custom health check
func (c *Checker) AddCheck(check func(context.Context) Check) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.checks = append(c.checks, check)
}

// Check performs all health checks and returns the response
func (c *Checker) Check(ctx context.Context) Response {
	response := Response{
		Status:    StatusUp,
		Service:   c.serviceName,
		Version:   c.version,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Checks:    make([]Check, 0),
	}

	// Check database
	if c.db != nil {
		dbCheck := c.checkDatabase(ctx)
		response.Checks = append(response.Checks, dbCheck)
		if dbCheck.Status != StatusUp {
			response.Status = StatusDown
		}
	}

	// Run custom checks
	c.mu.RLock()
	for _, check := range c.checks {
		result := check(ctx)
		response.Checks = append(response.Checks, result)
		if result.Status != StatusUp {
			response.Status = StatusDown
		}
	}
	c.mu.RUnlock()

	return response
}

// Liveness returns basic liveness status (is the service running?)
func (c *Checker) Liveness() Response {
	return Response{
		Status:    StatusUp,
		Service:   c.serviceName,
		Version:   c.version,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}
}

// Readiness returns readiness status (is the service ready to accept traffic?)
func (c *Checker) Readiness(ctx context.Context) Response {
	return c.Check(ctx)
}

func (c *Checker) checkDatabase(ctx context.Context) Check {
	start := time.Now()
	check := Check{Name: "database"}

	sqlDB, err := c.db.DB()
	if err != nil {
		check.Status = StatusDown
		check.Message = "Failed to get database connection: " + err.Error()
		return check
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		check.Status = StatusDown
		check.Message = "Failed to ping database: " + err.Error()
		return check
	}

	check.Status = StatusUp
	check.Latency = time.Since(start).String()
	return check
}

// GRPCCheck creates a check function for gRPC connection
func GRPCCheck(name string, checkFn func() error) func(context.Context) Check {
	return func(ctx context.Context) Check {
		start := time.Now()
		check := Check{Name: name}

		if err := checkFn(); err != nil {
			check.Status = StatusDown
			check.Message = err.Error()
		} else {
			check.Status = StatusUp
			check.Latency = time.Since(start).String()
		}

		return check
	}
}
