package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func main() {
	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/run",runHandler)

	fmt.Println("goboxd running on :8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

type RunRequest struct {
	Language string `json:"language"`
	Source string `json:"source"`
	ExpectedOutput string `json:"expected_output"`
}
type RunResponse struct {
	Status string `json:"status"`
	Stdout string `json:"stdout"`
}
func runHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type","application/json")
	
	var req RunRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}

	tempDir, err := os.MkdirTemp("", "goboxd-*")
	if err != nil {
		http.Error(w, "failed to create temp dir", http.StatusInternalServerError)
		return
	}
	defer os.RemoveAll(tempDir)
	fmt.Println("Temp directory:", tempDir)
	sourcePath := filepath.Join(tempDir, "main.py")
	err = os.WriteFile(sourcePath, []byte(req.Source), 0644)
	if err != nil {
		http.Error(w, "failed to write source file", http.StatusInternalServerError)
		return
	}
	fmt.Println("Source file written:", sourcePath)
	cmd := exec.Command("python3", sourcePath)

	output, err := cmd.CombinedOutput()
	
	fmt.Println("Execution output:")
	fmt.Println(string(output))

	if err != nil {
		fmt.Println("Execution error:", err)
	}
	fmt.Println("Language:", req.Language)
	fmt.Println("Source:", req.Source)

	status := "wrong answer"

	if string(output) == req.ExpectedOutput {
		status = "accepted"
	}

	response := RunResponse{
		Status: status,
		Stdout: string(output),
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}
