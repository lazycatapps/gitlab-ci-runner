package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

var (
	CONFIG_DIR = os.Getenv("CONFIG_DIR")
	configPath = strings.TrimRight(CONFIG_DIR, "/") + "/config.toml"
)

// Version information (set during build)
var (
	Version       = "dev"
	GitCommit     = "unknown"
	GitCommitFull = "unknown"
	GitBranch     = "unknown"
	BuildTime     = "unknown"
)

type RegisterRequest struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Token string `json:"token"`
}

type Runner struct {
	Name   string `json:"name"`
	URL    string `json:"url"`
	Token  string `json:"token"`
	Status string `json:"status"` // running, stopped, unknown
}

type ConfigToml struct {
	Concurrent int `toml:"concurrent"`
	Runners    []struct {
		Name     string   `toml:"name"`
		URL      string   `toml:"url"`
		Token    string   `toml:"token"`
		Executor string   `toml:"executor"`
		Shell    string   `toml:"shell,omitempty"`
		Builds   struct{} `toml:"builds_dir,omitempty"`
		Cache    struct{} `toml:"cache_dir,omitempty"`
	} `toml:"runners"`
}

const (
	logsDir = "/home/gitlab-runner/.gitlab-runner/logs"
	pidDir  = "/home/gitlab-runner/.gitlab-runner/pids"
)

func enableCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func main() {
	// Print version information
	log.Printf("GitLab CI Runner Manager")
	log.Printf("Version: %s", Version)
	log.Printf("Git Commit: %s (%s)", GitCommit, GitBranch)
	log.Printf("Build Time: %s", BuildTime)
	log.Printf("Config Path: %s", configPath)
	log.Println("========================================")

	http.HandleFunc("/api/runners/register", handleRegister)
	http.HandleFunc("/api/runners", handleRunners)
	http.HandleFunc("/api/runners/delete", handleDelete)
	http.HandleFunc("/api/runners/restart", handleRestart)
	http.HandleFunc("/api/runners/logs", handleLogs)
	http.HandleFunc("/api/version", handleVersion)
	http.Handle("/", http.FileServer(http.Dir("../frontend/static")))

	port := "8098"
	log.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleRegister(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Name == "" || req.URL == "" || req.Token == "" {
		http.Error(w, "Name, URL and Token are required", http.StatusBadRequest)
		return
	}

	// Execute gitlab-runner register command
	cmd := exec.Command("gitlab-runner", "register",
		"--non-interactive",
		"--url", req.URL,
		"--token", req.Token,
		"--name", req.Name,
		"--config", configPath,
		"--executor", "shell",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error registering runner: %v, output: %s", err, string(output))
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": fmt.Sprintf("Failed to register runner, please check the token and url: %s", string(output)),
			"output":  string(output),
		})
		return
	}

	log.Printf("Runner registered successfully: %s", req.Name)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Runner registered successfully",
		"output":  string(output),
	})
}

func handleRunners(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	runners, err := getRunners()
	if err != nil {
		log.Printf("Error getting runners: %v", err)
		http.Error(w, fmt.Sprintf("Failed to get runners: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(runners)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name  string `json:"name"`
		Token string `json:"token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Token == "" {
		http.Error(w, "Token is required", http.StatusBadRequest)
		return
	}

	if err := stopRunner(req.Name); err != nil {
		log.Printf("Warning: failed to stop runner %s: %v", req.Name, err)
	}

	// Execute gitlab-runner unregister command
	cmd := exec.Command("gitlab-runner", "unregister", "--token", req.Token)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Error unregistering runner: %v, output: %s", err, string(output))
		http.Error(w, fmt.Sprintf("Failed to unregister runner: %s", string(output)), http.StatusInternalServerError)
		return
	}

	log.Printf("Runner unregistered successfully: %s", req.Token)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Runner unregistered successfully",
	})
}

func handleRestart(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "Runner name is required", http.StatusBadRequest)
		return
	}

	// Stop the runner if it's running
	if err := stopRunner(req.Name); err != nil {
		log.Printf("Warning: failed to stop runner %s: %v", req.Name, err)
	}

	// Start the runner
	if err := startRunner(req.Name); err != nil {
		log.Printf("Error starting runner %s: %v", req.Name, err)
		http.Error(w, fmt.Sprintf("Failed to start runner: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Runner %s restarted successfully", req.Name)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Runner %s restarted successfully", req.Name),
	})
}

func handleLogs(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get runner name from query parameter
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "Runner name is required", http.StatusBadRequest)
		return
	}

	// Read logs from log file
	logContent, err := getRunnerLogs(name)
	if err != nil {
		logContent = fmt.Sprintf("Error reading logs for runner %s: %v", name, err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"name":    name,
		"logs":    logContent,
	})
}

func handleVersion(w http.ResponseWriter, r *http.Request) {
	enableCORS(w)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"version":       Version,
		"gitCommit":     GitCommit,
		"gitCommitFull": GitCommitFull,
		"gitBranch":     GitBranch,
		"buildTime":     BuildTime,
	})
}

func getRunners() ([]Runner, error) {
	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Return empty list if config doesn't exist yet
		return []Runner{}, nil
	}

	// Read the config.toml file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config ConfigToml
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Convert to Runner slice
	runners := make([]Runner, 0, len(config.Runners))
	for _, r := range config.Runners {
		runners = append(runners, Runner{
			Name:   r.Name,
			URL:    r.URL,
			Token:  r.Token,
			Status: getRunnerStatus(r.Name),
		})
	}

	return runners, nil
}

// getRunnerStatus checks if a specific runner process is running
func getRunnerStatus(name string) string {
	pidFile := filepath.Join(pidDir, name+".pid")

	// Check if PID file exists
	pidData, err := os.ReadFile(pidFile)
	if err != nil {
		return "stopped"
	}

	pid := strings.TrimSpace(string(pidData))

	// Check if process is still running
	cmd := exec.Command("ps", "-p", pid)
	if err := cmd.Run(); err != nil {
		// Process not running, clean up stale PID file
		os.Remove(pidFile)
		return "stopped"
	}

	return "running"
}

// startRunner starts a runner in the background with nohup
func startRunner(name string) error {
	// Ensure directories exist
	os.MkdirAll(logsDir, 0755)
	os.MkdirAll(pidDir, 0755)

	logFile := filepath.Join(logsDir, name+".log")
	pidFile := filepath.Join(pidDir, name+".pid")

	// Start runner with nohup in background
	cmd := exec.Command("nohup", "gitlab-runner", "run",
		"--config", configPath,
		"--working-directory", "/home/gitlab-runner",
		"-n", name)

	// Redirect output to log file
	outFile, err := os.Create(logFile)
	if err != nil {
		return fmt.Errorf("failed to create log file: %w", err)
	}
	defer outFile.Close()

	cmd.Stdout = outFile
	cmd.Stderr = outFile

	// Start the process
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start runner: %w", err)
	}

	// Save PID to file
	pid := fmt.Sprintf("%d", cmd.Process.Pid)
	if err := os.WriteFile(pidFile, []byte(pid), 0644); err != nil {
		cmd.Process.Kill()
		return fmt.Errorf("failed to write PID file: %w", err)
	}

	log.Printf("Started runner %s with PID %s", name, pid)
	return nil
}

// stopRunner stops a running runner by killing its process
func stopRunner(name string) error {
	pidFile := filepath.Join(pidDir, name+".pid")

	// Read PID from file
	pidData, err := os.ReadFile(pidFile)
	if err != nil {
		return fmt.Errorf("PID file not found: %w", err)
	}

	pid := strings.TrimSpace(string(pidData))

	// Kill the process
	cmd := exec.Command("kill", pid)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to kill process: %w", err)
	}

	// Remove PID file
	os.Remove(pidFile)

	log.Printf("Stopped runner %s (PID %s)", name, pid)
	return nil
}

// getRunnerLogs reads the log file for a runner
func getRunnerLogs(name string) (string, error) {
	logFile := filepath.Join(logsDir, name+".log")

	// Check if log file exists
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		return "No logs available yet. Runner may not have been started.", nil
	}

	// Read log file (last 1000 lines to avoid huge files)
	cmd := exec.Command("tail", "-n", "1000", logFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to read log file: %w", err)
	}

	return string(output), nil
}

func init() {
	// Ensure directories exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Printf("Warning: failed to create config directory: %v", err)
	}
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		log.Printf("Warning: failed to create logs directory: %v", err)
	}
	if err := os.MkdirAll(pidDir, 0755); err != nil {
		log.Printf("Warning: failed to create pids directory: %v", err)
	}
}
