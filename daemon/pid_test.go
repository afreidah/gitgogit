package daemon

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWritePID_RoundTrip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.pid")
	if err := WritePID(path); err != nil {
		t.Fatalf("WritePID() error: %v", err)
	}
	got, err := ReadPID(path)
	if err != nil {
		t.Fatalf("ReadPID() error: %v", err)
	}
	if got != os.Getpid() {
		t.Errorf("ReadPID() = %d, want %d", got, os.Getpid())
	}
}

func TestWritePID_CreatesParentDirs(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "c.pid")
	if err := WritePID(path); err != nil {
		t.Fatalf("WritePID() error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Errorf("pid file not created: %v", err)
	}
}

func TestWritePID_NoTmpLeftover(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.pid")
	if err := WritePID(path); err != nil {
		t.Fatalf("WritePID() error: %v", err)
	}
	if _, err := os.Stat(path + ".tmp"); !os.IsNotExist(err) {
		t.Error("expected .tmp file to be gone after WritePID")
	}
}

func TestReadPID_FileNotExist(t *testing.T) {
	_, err := ReadPID(filepath.Join(t.TempDir(), "nonexistent.pid"))
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestReadPID_InvalidContents(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.pid")
	if err := os.WriteFile(path, []byte("not-a-number\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := ReadPID(path)
	if err == nil {
		t.Error("expected error for non-integer pid file contents")
	}
}

func TestRemovePID_Idempotent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.pid")
	if err := RemovePID(path); err != nil {
		t.Errorf("first RemovePID() on missing file error: %v", err)
	}
	if err := RemovePID(path); err != nil {
		t.Errorf("second RemovePID() on missing file error: %v", err)
	}
}

func TestRemovePID_RemovesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.pid")
	if err := WritePID(path); err != nil {
		t.Fatal(err)
	}
	if err := RemovePID(path); err != nil {
		t.Fatalf("RemovePID() error: %v", err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("pid file still exists after RemovePID")
	}
}

func TestIsRunning_CurrentProcess(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.pid")
	if err := WritePID(path); err != nil {
		t.Fatal(err)
	}
	pid, running, err := IsRunning(path)
	if err != nil {
		t.Fatalf("IsRunning() error: %v", err)
	}
	if !running {
		t.Error("IsRunning() = false for current process, want true")
	}
	if pid != os.Getpid() {
		t.Errorf("IsRunning() pid = %d, want %d", pid, os.Getpid())
	}
}

func TestIsRunning_StalePID(t *testing.T) {
	path := filepath.Join(t.TempDir(), "stale.pid")
	// 99999999 exceeds Darwin's max PID and should never be a live process.
	if err := os.WriteFile(path, []byte("99999999\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, running, err := IsRunning(path)
	if err != nil {
		t.Fatalf("IsRunning() unexpected error: %v", err)
	}
	if running {
		t.Error("IsRunning() = true for impossible PID, want false")
	}
}

func TestIsRunning_NoFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.pid")
	pid, running, err := IsRunning(path)
	if err != nil {
		t.Errorf("IsRunning() error for missing file: %v", err)
	}
	if running {
		t.Error("IsRunning() = true for missing file, want false")
	}
	if pid != 0 {
		t.Errorf("IsRunning() pid = %d for missing file, want 0", pid)
	}
}
