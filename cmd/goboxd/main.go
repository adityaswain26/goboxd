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

type LanguageConfig struct{
	FileName string
	Compiled bool
}

var languages = map[string]LanguageConfig{
	"py3":{
		FileName: "main.py",
		Compiled: false,
	},
	"cpp":{
		FileName: "main.cpp",
		Compiled: true,
	},
}

func buildCommand(ctx context.Context, lang LanguageConfig, sourcePath string, tempDir string) (*exec.Cmd, error) {

	if !lang.Compiled {
		return exec.CommandContext(ctx, "python3", sourcePath), nil
	}

	binaryPath := filepath.Join(tempDir, "main")

	compileCmd := exec.Command("g++", sourcePath, "-o", binaryPath)

	output, err := compileCmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("compile error: %s", string(output))
	}

	return exec.CommandContext(ctx, binaryPath), nil
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
	lang, ok := languages[req.Language]
	if !ok {
		http.Error(w, "unsupported language", http.StatusBadRequest)
		return
	}
	sourcePath := filepath.Join(tempDir, lang.FileName)
	err = os.WriteFile(sourcePath, []byte(req.Source), 0644)
	if err != nil {
		http.Error(w, "failed to write source file", http.StatusInternalServerError)
		return
	}
	fmt.Println("Source file written:", sourcePath)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	cmd, err := buildCommand(ctx, lang, sourcePath, tempDir)
	if err != nil {
		response := RunResponse{
			Status: "compile error",
			Stdout: "",
			Stderr: err.Error(),
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
		return
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
