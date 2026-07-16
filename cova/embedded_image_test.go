package cova

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

const embeddedProgramImageScript = `
int initialized = 3;
int zeroed;

int add(int left, int right) {
	return left + right;
}

int script_main() {
	initialized = initialized + 1;
	zeroed = zeroed + 2;
	return add(initialized, zeroed);
}
`

func TestEmbeddedProgramImageFixture(t *testing.T) {
	linked := mustLinkProgram(t, embeddedProgramImageScript, 0, 0)
	image, status := BuildProgramImage(linked)
	if status != VMStatusOK {
		t.Fatalf("BuildProgramImage failed: %s", status)
	}

	fixturePath := embeddedProgramImageFixturePath(t)
	if os.Getenv("CCOVA_UPDATE_EMBEDDED") == "1" {
		writeEmbeddedProgramImageFixture(t, fixturePath, image)
	}

	fixture, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read embedded fixture: %v; refresh with CCOVA_UPDATE_EMBEDDED=1 go test -run TestEmbeddedProgramImageFixture -count=1", err)
	}
	if !bytes.Equal(fixture, image) {
		t.Fatalf("embedded fixture is stale; refresh with CCOVA_UPDATE_EMBEDDED=1 go test -run TestEmbeddedProgramImageFixture -count=1")
	}
}

func embeddedProgramImageFixturePath(t *testing.T) string {
	t.Helper()
	_, sourceFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve embedded fixture source path")
	}
	root := filepath.Dir(filepath.Dir(sourceFile))
	return filepath.Join(root, "embedded", "source", "test", "cpp", "go_program_image.bin")
}

func writeEmbeddedProgramImageFixture(t *testing.T, path string, image []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("create embedded fixture directory: %v", err)
	}
	temporary, err := os.CreateTemp(filepath.Dir(path), ".go_program_image-*.tmp")
	if err != nil {
		t.Fatalf("create temporary embedded fixture: %v", err)
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)

	if _, err := temporary.Write(image); err != nil {
		temporary.Close()
		t.Fatalf("write temporary embedded fixture: %v", err)
	}
	if err := temporary.Chmod(0o644); err != nil {
		temporary.Close()
		t.Fatalf("set embedded fixture permissions: %v", err)
	}
	if err := temporary.Close(); err != nil {
		t.Fatalf("close temporary embedded fixture: %v", err)
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		t.Fatalf("replace embedded fixture: %v", err)
	}
	t.Logf("updated %s (%d bytes)", path, len(image))
}
