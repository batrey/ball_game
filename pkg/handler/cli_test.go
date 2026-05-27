package handler

import (
	"bytes"
	"errors"
	"testing"
)

func TestRunCLIPrintsMaxTouches(t *testing.T) {
	// given
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	service := &fakeTouchCounter{result: 4}

	// when
	exitCode := RunCLI([]string{"-file", "players.txt"}, stdout, stderr, service)

	// then
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if stdout.String() != "4\n" {
		t.Fatalf("expected stdout %q, got %q", "4\n", stdout.String())
	}
	if stderr.String() != "" {
		t.Fatalf("expected empty stderr, got %q", stderr.String())
	}
	if service.gotFileName != "players.txt" {
		t.Fatalf("expected service to receive file name %q, got %q", "players.txt", service.gotFileName)
	}
}

func TestRunCLIRequiresFileFlag(t *testing.T) {
	// given
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	service := &fakeTouchCounter{}

	// when
	exitCode := RunCLI(nil, stdout, stderr, service)

	// then
	if exitCode != 2 {
		t.Fatalf("expected exit code 2, got %d", exitCode)
	}
	if stdout.String() != "" {
		t.Fatalf("expected empty stdout, got %q", stdout.String())
	}
	if stderr.String() != "missing required -file flag\n" {
		t.Fatalf("expected missing file error, got %q", stderr.String())
	}
	if service.calls != 0 {
		t.Fatalf("expected service not to be called, got %d calls", service.calls)
	}
}

func TestRunCLIPrintsHelp(t *testing.T) {
	// given
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	service := &fakeTouchCounter{}

	// when
	exitCode := RunCLI([]string{"-help"}, stdout, stderr, service)

	// then
	if exitCode != 0 {
		t.Fatalf("expected exit code 0, got %d", exitCode)
	}
	if stdout.String() != "" {
		t.Fatalf("expected empty stdout, got %q", stdout.String())
	}
	expectedStderr := "Usage of ballgame:\n  -file string\n    \tPath to the input file\n"
	if stderr.String() != expectedStderr {
		t.Fatalf("expected help output %q, got %q", expectedStderr, stderr.String())
	}
	if service.calls != 0 {
		t.Fatalf("expected service not to be called, got %d calls", service.calls)
	}
}

func TestRunCLIReturnsServiceError(t *testing.T) {
	// given
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	service := &fakeTouchCounter{err: errors.New("bad input")}

	// when
	exitCode := RunCLI([]string{"-file", "players.txt"}, stdout, stderr, service)

	// then
	if exitCode != 1 {
		t.Fatalf("expected exit code 1, got %d", exitCode)
	}
	if stdout.String() != "" {
		t.Fatalf("expected empty stdout, got %q", stdout.String())
	}
	if stderr.String() != "error: bad input\n" {
		t.Fatalf("expected service error, got %q", stderr.String())
	}
}

type fakeTouchCounter struct {
	result      int
	err         error
	gotFileName string
	calls       int
}

func (f *fakeTouchCounter) MaxTouchesFromFile(fileName string) (int, error) {
	f.calls++
	f.gotFileName = fileName
	if f.err != nil {
		return 0, f.err
	}
	return f.result, nil
}
