package storage

// storage.go provides minimal JSON load/save helpers with atomic writes
// used by the app to persist active tasks and archived tasks.
import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func Load(path string, dst any) error {
	b, err := os.ReadFile(path)
	if err != nil { // Checks if file doesn't exist, which is not an error for Load.
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read %s: %w", path, err)
	}
	if len(bytes.TrimSpace(b)) == 0 { // Treat empty file as empty JSON, which is not an error for Load.
		return nil
	}
	if err := json.Unmarshal(b, dst); err != nil { // Checks if JSON is invalid, which is an error for Load.
		return fmt.Errorf("unmarshal %s: %w", path, err)
	}
	return nil
}

func Save(path string, src any) error {
	// Ensure the directory exists before writing the file.
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("mkdir %s: %w", filepath.Dir(path), err)
	}

	b, err := json.MarshalIndent(src, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}
	b = append(b, '\n')

	// Write to a temporary file first, then rename it to the target path for atomicity.
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, b, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", tmpPath, err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("rename %s -> %s: %w", tmpPath, path, err)
	}
	return nil
}
