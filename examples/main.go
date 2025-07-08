package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/brennoo/go-gmp"
	"github.com/brennoo/go-gmp/client"
	"github.com/brennoo/go-gmp/connections"
)

// Configuration constants
const (
	gvmdHost     = "localhost:9390"
	username     = "admin"
	password     = "admin"
	scannerName  = "OpenVAS Default"
	configName   = "Full and fast"
	portListName = "All IANA assigned TCP"
	targetName   = "localhost"
	targetHosts  = "127.0.0.1"
	taskName     = "New Task"
	pollInterval = 10 * time.Second
)

// GMPClient wraps the GMP client with additional functionality
type GMPClient struct {
	client gmp.Client
	conn   gmp.Connection
}

// NewGMPClient creates a new GMP client with TLS connection
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

// Close closes the underlying connection
func (g *GMPClient) Close() error {
	return g.conn.Close()
}

// Authenticate authenticates with the GMP server with retry logic
func (g *GMPClient) Authenticate(ctx context.Context, username, password string) error {
	log.Printf("Authenticating as user: %s", username)

	auth := &gmp.AuthenticateCommand{
		Credentials: gmp.AuthenticateCredentials{
			Username: username,
			Password: password,
		},
	}

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

// GetScanner retrieves a scanner by name
func (g *GMPClient) GetScanner(ctx context.Context, name string) (*gmp.GetScannersResponseScanner, error) {
	log.Printf("Getting scanner: %s", name)

	cmd := &gmp.GetScannersCommand{
		Filter: fmt.Sprintf(`name="%s"`, name),
	}

	resp, err := g.client.GetScanners(cmd)
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

// GetConfig retrieves a configuration by name
func (g *GMPClient) GetConfig(ctx context.Context, name string) (string, error) {
	log.Printf("Getting configuration: %s", name)

	cmd := &gmp.GetConfigsCommand{
		Filter: fmt.Sprintf(`name="%s"`, name),
	}

	resp, err := g.client.GetConfigs(cmd)
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

// GetPortList retrieves a port list by name
func (g *GMPClient) GetPortList(ctx context.Context, name string) (string, error) {
	log.Printf("Getting port list: %s", name)

	cmd := &gmp.GetPortListsCommand{
		Filter: fmt.Sprintf(`name="%s"`, name),
	}

	resp, err := g.client.GetPortLists(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to get port lists: %w", err)
	}

	if resp.Status != "200" {
		return "", fmt.Errorf("get port lists failed with status %s: %s", resp.Status, resp.StatusText)
	}

	if len(resp.PortList) == 0 {
		return "", fmt.Errorf("port list '%s' not found", name)
	}

	log.Printf("Found port list: %s (ID: %s)", resp.PortList[0].Name, resp.PortList[0].ID)
	return resp.PortList[0].ID, nil
}

// GetTarget retrieves a target by name
func (g *GMPClient) GetTarget(ctx context.Context, name string) (*gmp.GetTargetsResponseTarget, error) {
	log.Printf("Getting target: %s", name)

	cmd := &gmp.GetTargetsCommand{
		Filter: fmt.Sprintf(`name="%s"`, name),
	}

	resp, err := g.client.GetTargets(cmd)
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

// DeleteTarget deletes a target by ID
func (g *GMPClient) DeleteTarget(ctx context.Context, targetID string) error {
	log.Printf("Deleting target: %s", targetID)

	cmd := &gmp.DeleteTargetCommand{
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

// DeleteTargetWithDependencies deletes a target and all tasks that use it
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

// GetTargetByID retrieves a target by ID
func (g *GMPClient) GetTargetByID(ctx context.Context, targetID string) (*gmp.GetTargetsResponseTarget, error) {
	log.Printf("Getting target by ID: %s", targetID)

	cmd := &gmp.GetTargetsCommand{
		TargetID: targetID,
		Tasks:    true, // Include tasks that use this target
	}

	resp, err := g.client.GetTargets(cmd)
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

// GetTask retrieves a task by name
func (g *GMPClient) GetTask(ctx context.Context, name string) (*gmp.GetTasksResponseTask, error) {
	log.Printf("Getting task: %s", name)

	cmd := &gmp.GetTasksCommand{
		Filter: fmt.Sprintf(`name="%s"`, name),
	}

	resp, err := g.client.GetTasks(cmd)
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

// DeleteTask deletes a task by ID
func (g *GMPClient) DeleteTask(ctx context.Context, taskID string) error {
	log.Printf("Deleting task: %s", taskID)

	cmd := &gmp.DeleteTaskCommand{
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

// CreateTarget creates a new target, deleting existing one if it exists
func (g *GMPClient) CreateTarget(ctx context.Context, name, hosts, portListID string) (*gmp.CreateTargetResponse, error) {
	log.Printf("Creating target: %s", name)

	// Check if target already exists
	if existingTarget, err := g.GetTarget(ctx, name); err == nil {
		log.Printf("Target '%s' already exists (ID: %s), deleting it first", name, existingTarget.ID)
		if err := g.DeleteTargetWithDependencies(ctx, existingTarget.ID); err != nil {
			return nil, fmt.Errorf("failed to delete existing target: %w", err)
		}
	}

	cmd := &gmp.CreateTargetCommand{
		Name:     name,
		Hosts:    hosts,
		PortList: &gmp.CreateTargetPortList{ID: portListID},
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

// CreateTask creates a new task, deleting existing one if it exists
func (g *GMPClient) CreateTask(ctx context.Context, name, configID, targetID, scannerID string) (*gmp.CreateTaskResponse, error) {
	log.Printf("Creating task: %s", name)

	// Check if task already exists
	if existingTask, err := g.GetTask(ctx, name); err == nil {
		log.Printf("Task '%s' already exists (ID: %s), deleting it first", name, existingTask.ID)
		if err := g.DeleteTask(ctx, existingTask.ID); err != nil {
			return nil, fmt.Errorf("failed to delete existing task: %w", err)
		}
	}

	cmd := &gmp.CreateTaskCommand{
		Name: name,
		Config: &gmp.CreateTaskConfig{
			ID: configID,
		},
		Target: &gmp.CreateTaskTarget{
			ID: targetID,
		},
		Scanner: &gmp.CreateTaskScanner{
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

// StartTask starts a task
func (g *GMPClient) StartTask(ctx context.Context, taskID string) error {
	log.Printf("Starting task: %s", taskID)

	cmd := &gmp.StartTaskCommand{
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

// WaitForTaskCompletion waits for a task to complete
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

// getTaskProgress gets the current progress of a task
func (g *GMPClient) getTaskProgress(ctx context.Context, taskID string) (int, error) {
	cmd := &gmp.GetTasksCommand{
		TaskID: taskID,
	}

	resp, err := g.client.GetTasks(cmd)
	if err != nil {
		return 0, err
	}

	if resp.Status != "200" {
		return 0, fmt.Errorf("get tasks failed with status %s: %s", resp.Status, resp.StatusText)
	}

	if len(resp.Task) == 0 {
		return 0, fmt.Errorf("task not found")
	}

	progressStr := resp.Task[0].Progress
	if progressStr == "" {
		return 0, fmt.Errorf("progress value is empty")
	}
	progress, err := strconv.Atoi(progressStr)
	if err != nil {
		return 0, fmt.Errorf("invalid progress value: %s", progressStr)
	}

	return progress, nil
}

// GetResults retrieves results for a task
func (g *GMPClient) GetResults(ctx context.Context, taskID string) ([]gmp.Result, error) {
	log.Printf("Getting results for task: %s", taskID)

	cmd := &gmp.GetResultsCommand{
		TaskID: taskID,
		Filter: "min_qod=0 rows=1000",
	}

	resp, err := g.client.GetResults(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to get results: %w", err)
	}

	if resp.Status != "200" {
		return nil, fmt.Errorf("get results failed with status %s: %s", resp.Status, resp.StatusText)
	}

	log.Printf("Retrieved %d results", len(resp.Results))
	return resp.Results, nil
}

// PrintResults prints the results in a formatted way
func PrintResults(results []gmp.Result) {
	if len(results) == 0 {
		log.Println("No results found")
		return
	}

	log.Printf("Found %d results:", len(results))
	log.Println("=" + strings.Repeat("=", 80))

	for i, result := range results {
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
	defer gmpClient.Close()

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

	// Get results
	results, err := gmpClient.GetResults(ctx, task.ID)
	if err != nil {
		log.Fatalf("Failed to get results: %v", err)
	}

	// Print results
	PrintResults(results)
}
