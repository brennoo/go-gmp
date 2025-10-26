//nolint:dupl,errcheck,unparam // Example patterns are intentionally repetitive

package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/brennoo/go-gmp"
	"github.com/brennoo/go-gmp/client"
	"github.com/brennoo/go-gmp/commands"
	"github.com/brennoo/go-gmp/commands/filtering"
	"github.com/brennoo/go-gmp/connections"
)

// Configuration constants.
const (
	gvmdHost     = "localhost:9390"
	username     = "admin"
	password     = "admin"
	scannerName  = "OpenVAS Default"
	configName   = "Full and fast"
	portListName = "All IANA assigned TCP"
	targetName   = "scanme.nmap.org"
	targetHosts  = "scanme.nmap.org"
	taskName     = "Scan scanme.nmap.org"
	pollInterval = 10 * time.Second
)

// GMPClient wraps the GMP client with additional functionality.
type GMPClient struct {
	client gmp.Client
	conn   gmp.Connection
}

// NewGMPClient creates a new GMP client with TLS connection.
func NewGMPClient(host string) (*GMPClient, error) {
	log.Printf("Connecting to GVMD at %s...", host)

	conn, err := connections.NewTLSConnection(host, true)
	if err != nil {
		return nil, fmt.Errorf("failed to create TLS connection: %w", err)
	}

	gmpClient := client.New(conn)

	return &GMPClient{
		client: gmpClient,
		conn:   conn,
	}, nil
}

// Close closes the underlying connection.
func (g *GMPClient) Close() error {
	return g.conn.Close()
}

// Authenticate authenticates with the GMP server with retry logic.
func (g *GMPClient) Authenticate(ctx context.Context, username, password string) error {
	log.Printf("Authenticating as user: %s", username)

	auth := &commands.Authenticate{Credentials: &commands.AuthenticateCredentials{Username: username, Password: password}}

	resp, err := g.client.Authenticate(auth)
	if err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}

	if resp.Status != "200" {
		return fmt.Errorf("authentication failed with status %s: %s", resp.Status, resp.StatusText)
	}

	log.Printf("Authentication successful (role: %s, timezone: %s)", resp.Role, resp.Timezone)
	return nil
}

// GetScanner retrieves a scanner by name.
func (g *GMPClient) GetScanner(ctx context.Context, name string) (*commands.Scanner, error) {
	log.Printf("Getting scanner: %s", name)

	resp, err := g.client.GetScanners(ctx, fmt.Sprintf(`name="%s"`, name))
	if err != nil {
		return nil, fmt.Errorf("failed to get scanners: %w", err)
	}

	if resp.Status != "200" {
		return nil, fmt.Errorf("get scanners failed with status %s: %s", resp.Status, resp.StatusText)
	}

	if len(resp.Scanner) == 0 {
		return nil, fmt.Errorf("scanner '%s' not found", name)
	}

	log.Printf("Found scanner: %s (ID: %s)", resp.Scanner[0].Name, resp.Scanner[0].ID)
	return &resp.Scanner[0], nil
}

// GetConfig retrieves a configuration by name.
func (g *GMPClient) GetConfig(ctx context.Context, name string) (string, error) {
	log.Printf("Getting configuration: %s", name)

	resp, err := g.client.GetConfigs(ctx, fmt.Sprintf(`name="%s"`, name))
	if err != nil {
		return "", fmt.Errorf("failed to get configs: %w", err)
	}

	if resp.Status != "200" {
		return "", fmt.Errorf("get configs failed with status %s: %s", resp.Status, resp.StatusText)
	}

	if len(resp.Config) == 0 {
		return "", fmt.Errorf("configuration '%s' not found", name)
	}

	log.Printf("Found configuration: %s (ID: %s)", resp.Config[0].Name, resp.Config[0].ID)
	return resp.Config[0].ID, nil
}

// GetPortList retrieves a port list by name.
func (g *GMPClient) GetPortList(ctx context.Context, name string) (string, error) {
	log.Printf("Getting port list: %s", name)

	resp, err := g.client.GetPortLists(ctx, fmt.Sprintf(`name="%s"`, name))
	if err != nil {
		return "", fmt.Errorf("failed to get port lists: %w", err)
	}

	if resp.Status != "200" {
		return "", fmt.Errorf("get port lists failed with status %s: %s", resp.Status, resp.StatusText)
	}

	if len(resp.PortLists) == 0 {
		return "", fmt.Errorf("port list '%s' not found", name)
	}

	log.Printf("Found port list: %s (ID: %s)", resp.PortLists[0].Name, resp.PortLists[0].ID)
	return resp.PortLists[0].ID, nil
}

// GetTarget retrieves a target by name.
func (g *GMPClient) GetTarget(ctx context.Context, name string) (*commands.Target, error) {
	log.Printf("Getting target: %s", name)

	resp, err := g.client.GetTargets(ctx, fmt.Sprintf(`name="%s"`, name))
	if err != nil {
		return nil, fmt.Errorf("failed to get targets: %w", err)
	}

	if resp.Status != "200" {
		return nil, fmt.Errorf("get targets failed with status %s: %s", resp.Status, resp.StatusText)
	}

	if len(resp.Target) == 0 {
		return nil, fmt.Errorf("target '%s' not found", name)
	}

	log.Printf("Found target: %s (ID: %s)", resp.Target[0].Name, resp.Target[0].ID)
	return &resp.Target[0], nil
}

// DeleteTarget deletes a target by ID.
func (g *GMPClient) DeleteTarget(ctx context.Context, targetID string) error {
	log.Printf("Deleting target: %s", targetID)

	cmd := &commands.DeleteTarget{
		TargetID: targetID,
		Ultimate: true, // Delete permanently, not just to trash
	}

	resp, err := g.client.DeleteTarget(cmd)
	if err != nil {
		return fmt.Errorf("failed to delete target: %w", err)
	}

	if resp.Status != "200" {
		return fmt.Errorf("delete target failed with status %s: %s", resp.Status, resp.StatusText)
	}

	log.Printf("Deleted target: %s", targetID)
	return nil
}

// DeleteTargetWithDependencies deletes a target and all tasks that use it.
func (g *GMPClient) DeleteTargetWithDependencies(ctx context.Context, targetID string) error {
	log.Printf("Deleting target with dependencies: %s", targetID)

	// First, get the target to see if it's in use
	target, err := g.GetTargetByID(ctx, targetID)
	if err != nil {
		return fmt.Errorf("failed to get target info: %w", err)
	}

	// If target has tasks, delete them first
	if target.Tasks != nil && len(target.Tasks.Task) > 0 {
		log.Printf("Target has %d associated tasks, deleting them first", len(target.Tasks.Task))
		for _, task := range target.Tasks.Task {
			if err := g.DeleteTask(ctx, task.ID); err != nil {
				log.Printf("Warning: failed to delete task %s: %v", task.ID, err)
				// Continue with other tasks even if one fails
			}
		}
	}

	// Now delete the target
	return g.DeleteTarget(ctx, targetID)
}

// GetTargetByID retrieves a target by ID.
func (g *GMPClient) GetTargetByID(ctx context.Context, targetID string) (*commands.Target, error) {
	log.Printf("Getting target by ID: %s", targetID)

	resp, err := g.client.GetTargets(ctx, fmt.Sprintf(`target_id="%s"`, targetID))
	if err != nil {
		return nil, fmt.Errorf("failed to get target: %w", err)
	}

	if resp.Status != "200" {
		return nil, fmt.Errorf("get target failed with status %s: %s", resp.Status, resp.StatusText)
	}

	if len(resp.Target) == 0 {
		return nil, fmt.Errorf("target with ID '%s' not found", targetID)
	}

	log.Printf("Found target: %s (ID: %s)", resp.Target[0].Name, resp.Target[0].ID)
	return &resp.Target[0], nil
}

// GetTask retrieves a task by name.
func (g *GMPClient) GetTask(ctx context.Context, name string) (*commands.GetTasksResponseTask, error) {
	log.Printf("Getting task: %s", name)

	resp, err := g.client.GetTasks(ctx, fmt.Sprintf(`name="%s"`, name))
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	if resp.Status != "200" {
		return nil, fmt.Errorf("get tasks failed with status %s: %s", resp.Status, resp.StatusText)
	}

	if len(resp.Task) == 0 {
		return nil, fmt.Errorf("task '%s' not found", name)
	}

	log.Printf("Found task: %s (ID: %s)", resp.Task[0].Name, resp.Task[0].ID)
	return &resp.Task[0], nil
}

// DeleteTask deletes a task by ID.
func (g *GMPClient) DeleteTask(ctx context.Context, taskID string) error {
	log.Printf("Deleting task: %s", taskID)

	cmd := &commands.DeleteTask{
		TaskID:   taskID,
		Ultimate: true,
	}

	resp, err := g.client.DeleteTask(cmd)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	if resp.Status != "200" {
		return fmt.Errorf("delete task failed with status %s: %s", resp.Status, resp.StatusText)
	}

	log.Printf("Deleted task: %s", taskID)
	return nil
}

// CreateTarget creates a new target, deleting existing one if it exists.
func (g *GMPClient) CreateTarget(ctx context.Context, name, hosts, portListID string) (*commands.CreateTargetResponse, error) {
	log.Printf("Creating target: %s", name)

	// Check if target already exists
	if existingTarget, err := g.GetTarget(ctx, name); err == nil {
		log.Printf("Target '%s' already exists (ID: %s), deleting it first", name, existingTarget.ID)
		if err := g.DeleteTargetWithDependencies(ctx, existingTarget.ID); err != nil {
			return nil, fmt.Errorf("failed to delete existing target: %w", err)
		}
	}

	cmd := &commands.CreateTarget{
		Name:     name,
		Hosts:    hosts,
		PortList: &commands.TargetPortList{ID: portListID},
	}

	resp, err := g.client.CreateTarget(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to create target: %w", err)
	}

	if resp.Status != "201" {
		return nil, fmt.Errorf("create target failed with status %s: %s", resp.Status, resp.StatusText)
	}

	log.Printf("Created target: %s (ID: %s)", name, resp.ID)
	return resp, nil
}

// CreateTask creates a new task, deleting existing one if it exists.
func (g *GMPClient) CreateTask(ctx context.Context, name, configID, targetID, scannerID string) (*commands.CreateTaskResponse, error) {
	log.Printf("Creating task: %s", name)

	// Check if task already exists
	if existingTask, err := g.GetTask(ctx, name); err == nil {
		log.Printf("Task '%s' already exists (ID: %s), deleting it first", name, existingTask.ID)
		if err := g.DeleteTask(ctx, existingTask.ID); err != nil {
			return nil, fmt.Errorf("failed to delete existing task: %w", err)
		}
	}

	cmd := &commands.CreateTask{
		Name: name,
		Config: &commands.TaskConfig{
			ID: configID,
		},
		Target: &commands.TaskTarget{
			ID: targetID,
		},
		Scanner: &commands.TaskScanner{
			ID: scannerID,
		},
	}

	resp, err := g.client.CreateTask(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	if resp.Status != "201" {
		return nil, fmt.Errorf("create task failed with status %s: %s", resp.Status, resp.StatusText)
	}

	log.Printf("Created task: %s (ID: %s)", name, resp.ID)
	return resp, nil
}

// StartTask starts a task.
func (g *GMPClient) StartTask(ctx context.Context, taskID string) error {
	log.Printf("Starting task: %s", taskID)

	cmd := &commands.StartTask{
		TaskID: taskID,
	}

	resp, err := g.client.StartTask(cmd)
	if err != nil {
		return fmt.Errorf("failed to start task: %w", err)
	}

	if resp.Status != "202" {
		return fmt.Errorf("start task failed with status %s: %s", resp.Status, resp.StatusText)
	}

	log.Printf("Task started successfully (report ID: %s)", resp.ReportID)
	return nil
}

// WaitForTaskCompletion waits for a task to complete.
func (g *GMPClient) WaitForTaskCompletion(ctx context.Context, taskID string) error {
	log.Println("Monitoring task progress...")

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			progress, err := g.getTaskProgress(ctx, taskID)
			if err != nil {
				log.Printf("Warning: failed to get task progress: %v", err)
				continue
			}

			if progress == -1 {
				log.Println("Task completed")
				return nil
			}

			log.Printf("Task progress: %d%%", progress)

			if progress >= 100 {
				log.Println("Task completed")
				return nil
			}
		}
	}
}

// getTaskProgress gets the current progress of a task.
func (g *GMPClient) getTaskProgress(ctx context.Context, taskID string) (int, error) {
	resp, err := g.client.GetTasks(ctx, fmt.Sprintf(`task_id="%s"`, taskID))
	if err != nil {
		return 0, err
	}

	if resp.Status != "200" {
		return 0, fmt.Errorf("get tasks failed with status %s: %s", resp.Status, resp.StatusText)
	}

	if len(resp.Task) == 0 {
		return 0, fmt.Errorf("task not found")
	}

	progress := resp.Task[0].Progress
	return progress, nil
}

// PrintResultsWithIterator prints all results using an iterator.
func (g *GMPClient) PrintResultsWithIterator(ctx context.Context, taskID string) error {
	log.Printf("Getting results for task: %s using iterator", taskID)

	// Create an iterator for results
	resultIter := g.client.Results(ctx, 50, taskID)
	defer resultIter.Close()

	resultCount := 0
	log.Println("Scanning results...")
	log.Println("=" + strings.Repeat("=", 80))

	// Iterate through all results
	for resultIter.Next() {
		result := resultIter.Current()
		resultCount++

		fmt.Printf("[%d] %s (Severity: %.1f)\n", resultCount, result.Name, result.Severity)
		if result.Description != "" {
			fmt.Printf("     Description: %s\n", result.Description)
		}
		if result.Host.Value != "" {
			fmt.Printf("     Host: %s\n", result.Host.Value)
		}
		if result.Port != "" {
			fmt.Printf("     Port: %s\n", result.Port)
		}
		fmt.Println()
	}

	// Check for errors
	if err := resultIter.Err(); err != nil {
		return fmt.Errorf("error during iteration: %w", err)
	}

	if resultCount == 0 {
		log.Println("No results found")
	} else {
		log.Printf("Total results found: %d", resultCount)
		log.Printf("Total available: %d", resultIter.Total())
	}

	return nil
}

// PrintResultsWithFiltering demonstrates the new filtering system.
func (g *GMPClient) PrintResultsWithFiltering(ctx context.Context, taskID string) error {
	log.Printf("Getting results for task: %s using filtering system", taskID)

	// Create filter options for high and critical severity results
	opts := []filtering.Option{
		filtering.WithTask(taskID),
		filtering.WithSeverityGreaterThanOrEqual(7.0), // High and Critical
		filtering.WithSortDesc(filtering.SortByResultSeverity),
		filtering.WithLimit(20), // Limit to top 20 results
	}

	// Use the new filtering method
	var args []filtering.FilterArg
	for _, opt := range opts {
		args = append(args, opt)
	}
	resp, err := g.client.GetResults(ctx, args...)
	if err != nil {
		return fmt.Errorf("failed to get filtered results: %w", err)
	}

	if resp.Status != "200" {
		return fmt.Errorf("get results failed with status %s: %s", resp.Status, resp.StatusText)
	}

	log.Printf("Found %d high/critical results (showing top 20)", len(resp.Results))
	log.Println("=" + strings.Repeat("=", 80))

	for i, result := range resp.Results {
		fmt.Printf("[%d] %s (Severity: %.1f)\n", i+1, result.Name, result.Severity)
		if result.Description != "" {
			fmt.Printf("     Description: %s\n", result.Description)
		}
		if result.Host.Value != "" {
			fmt.Printf("     Host: %s\n", result.Host.Value)
		}
		if result.Port != "" {
			fmt.Printf("     Port: %s\n", result.Port)
		}
		fmt.Println()
	}

	return nil
}

// PrintTasksWithFiltering demonstrates filtering tasks.
func (g *GMPClient) PrintTasksWithFiltering(ctx context.Context) error {
	log.Println("Getting tasks using filtering system")

	// Create filter options for running tasks
	opts := []filtering.Option{
		filtering.WithStatus(filtering.StatusRunning),
		filtering.WithSortDesc(filtering.SortByCreated),
		filtering.WithLimit(10), // Show last 10 running tasks
	}

	// Use the new filtering method
	var args []filtering.FilterArg
	for _, opt := range opts {
		args = append(args, opt)
	}
	resp, err := g.client.GetTasks(ctx, args...)
	if err != nil {
		return fmt.Errorf("failed to get filtered tasks: %w", err)
	}

	if resp.Status != "200" {
		return fmt.Errorf("get tasks failed with status %s: %s", resp.Status, resp.StatusText)
	}

	log.Printf("Found %d running tasks", len(resp.Task))
	log.Println("=" + strings.Repeat("=", 80))

	for i, task := range resp.Task {
		fmt.Printf("[%d] %s (Status: %s, Progress: %d%%)\n", i+1, task.Name, task.Status, task.Progress)
		if task.Comment != "" {
			fmt.Printf("     Comment: %s\n", task.Comment)
		}
		fmt.Printf("     Created: %s\n", task.CreationTime)
		fmt.Println()
	}

	return nil
}

func main() {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Create GMP client
	gmpClient, err := NewGMPClient(gvmdHost)
	if err != nil {
		log.Fatalf("Failed to create GMP client: %v", err)
	}
	defer gmpClient.Close() //nolint:errcheck // Ignoring close error

	// Authenticate
	if err := gmpClient.Authenticate(ctx, username, password); err != nil {
		log.Fatalf("Authentication failed: %v", err)
	}

	// Get scanner
	scanner, err := gmpClient.GetScanner(ctx, scannerName)
	if err != nil {
		log.Fatalf("Failed to get scanner: %v", err)
	}

	// Get configuration
	configID, err := gmpClient.GetConfig(ctx, configName)
	if err != nil {
		log.Fatalf("Failed to get configuration: %v", err)
	}

	// Get port list
	portListID, err := gmpClient.GetPortList(ctx, portListName)
	if err != nil {
		log.Fatalf("Failed to get port list: %v", err)
	}

	// Create target
	target, err := gmpClient.CreateTarget(ctx, targetName, targetHosts, portListID)
	if err != nil {
		log.Fatalf("Failed to create target: %v", err)
	}

	// Create task
	task, err := gmpClient.CreateTask(ctx, taskName, configID, target.ID, scanner.ID)
	if err != nil {
		log.Fatalf("Failed to create task: %v", err)
	}

	// Start task
	if err := gmpClient.StartTask(ctx, task.ID); err != nil {
		log.Fatalf("Failed to start task: %v", err)
	}

	// Wait for completion
	if err := gmpClient.WaitForTaskCompletion(ctx, task.ID); err != nil {
		log.Fatalf("Task failed or timed out: %v", err)
	}

	// Get and print results using iterator (old approach)
	if err := gmpClient.PrintResultsWithIterator(ctx, task.ID); err != nil {
		log.Fatalf("Failed to get results: %v", err)
	}

	// Demonstrate the new filtering system
	log.Println("\n" + strings.Repeat("=", 80))
	log.Println("DEMONSTRATING NEW FILTERING SYSTEM")
	log.Println(strings.Repeat("=", 80))

	// Show filtered results (high/critical severity only)
	if err := gmpClient.PrintResultsWithFiltering(ctx, task.ID); err != nil {
		log.Printf("Warning: Failed to get filtered results: %v", err)
	}

	// Show running tasks using filtering
	if err := gmpClient.PrintTasksWithFiltering(ctx); err != nil {
		log.Printf("Warning: Failed to get filtered tasks: %v", err)
	}
}
