package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
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
	Stderr string `json:"stderr"`
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
	var sourcePath string

	if req.Language == "py3" {
		sourcePath = filepath.Join(tempDir, "main.py")
	} else if  req.Language == "cpp" {
		sourcePath = filepath.Join(tempDir, "main.cpp")
	} else {
		http.Error(w, "unsupported language", http.StatusBadRequest)
		return
	}
	err = os.WriteFile(sourcePath, []byte(req.Source), 0644)
	if err != nil {
		http.Error(w, "failed to write source file", http.StatusInternalServerError)
		return
	}
	fmt.Println("Source file written:", sourcePath)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var cmd *exec.Cmd

	if req.Language == "py3" {
		cmd = exec.CommandContext(ctx, "python3", sourcePath)
	} else if req.Language == "cpp" {
		binaryPath := filepath.Join(tempDir,"main")
		compileCmd := exec.Command("g++", sourcePath, "-o", binaryPath)
		compileOutput, compileErr := compileCmd.CombinedOutput()
		if compileErr != nil {
			response := RunResponse{
				Status: "compile error",
				Stdout: "",
				Stderr: string(compileOutput),
			}
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(response)
			return
		}
		cmd = exec.CommandContext(ctx, binaryPath)
	}
	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer

	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()

	fmt.Println("Execution output:")
	fmt.Println(stdoutBuf.String())

	if err != nil {
		fmt.Println("Execution error:", err)
	}
	fmt.Println("Language:", req.Language)
	fmt.Println("Source:", req.Source)

	status := "wrong answer"

	if ctx.Err() == context.DeadlineExceeded {
		status = "time limit exceeded"
	}else if err != nil {
		status = "runtime error"
	}else if stdoutBuf.String() == req.ExpectedOutput {
		status = "accepted"
	}

	response := RunResponse{
		Status: status,
		Stdout: stdoutBuf.String(),
		Stderr: stderrBuf.String(),
	}

	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(response)
}
