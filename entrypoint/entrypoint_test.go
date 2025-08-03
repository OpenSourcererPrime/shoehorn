package entrypoint

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/OpenSourcererPrime/shoehorn/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEntryPointWithEmptyConfig(t *testing.T) {
	cfg := &config.Config{}

	// Create EntryPoint with empty config
	ep, err := NewEntryPoint(cfg)
	require.NoError(t, err)
	defer ep.Close()

	assert.NotNil(t, ep)
	assert.NotNil(t, ep.watcher)
	assert.Equal(t, *cfg, ep.appConfig)
}

func TestEntryPointWithInvalidTemplate(t *testing.T) {
	testDir := t.TempDir()

	inputFile := filepath.Join(testDir, "input.txt")
	err := os.WriteFile(inputFile, []byte("test content"), 0o644)
	require.NoError(t, err)

	// Create invalid template file with syntax error
	templateFile := filepath.Join(testDir, "invalid.tmpl")
	err = os.WriteFile(templateFile, []byte("{{.input} missing closing brace"), 0o644)
	require.NoError(t, err)

	cfg := &config.Config{
		Generate: []config.GenerateConfig{
			{
				Name:     "output.txt",
				Path:     testDir,
				Strategy: "template",
				Template: templateFile,
				Inputs: []config.InputFile{
					{Name: "input", Path: inputFile},
				},
			},
		},
	}

	ep, err := NewEntryPoint(cfg)
	require.NoError(t, err)
	defer ep.Close()

	assert.NotNil(t, ep)

	// Output file should not exist or be empty
	outputPath := filepath.Join(testDir, "output.txt")
	_, err = os.Stat(outputPath)
	if err == nil {
		// If file exists, it should be empty
		content, err := os.ReadFile(outputPath)
		if err == nil {
			assert.Empty(t, content, "Output file should be empty with invalid template")
		}
	}
}

func TestEntryPointWithInvalidWatchPath(t *testing.T) {
	testDir := t.TempDir()
	inputFile := filepath.Join(testDir, "input.txt")
	err := os.WriteFile(inputFile, []byte("test content"), 0o644)
	require.NoError(t, err)

	nonExistentFile := filepath.Join(testDir, "nonexistent.txt")

	cfg := &config.Config{
		Generate: []config.GenerateConfig{
			{
				Name:     "output.txt",
				Path:     testDir,
				Strategy: "append",
				Inputs: []config.InputFile{
					{Name: "input1", Path: inputFile},
					{Name: "input2", Path: nonExistentFile},
				},
			},
		},
	}

	ep, err := NewEntryPoint(cfg)
	require.NoError(t, err)
	defer ep.Close()

	assert.NotNil(t, ep)

	outputPath := filepath.Join(testDir, "output.txt")
	content, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	assert.Equal(t, "test content\n", string(content))
}

func TestEntryPointWithReadOnlyOutputDir(t *testing.T) {
	testDir := t.TempDir()

	inputFile := filepath.Join(testDir, "input.txt")
	err := os.WriteFile(inputFile, []byte("test content"), 0o644)
	require.NoError(t, err)

	outputDir := filepath.Join(testDir, "readonly")
	err = os.MkdirAll(outputDir, 0o755)
	require.NoError(t, err)
	err = os.Chmod(outputDir, 0o555)
	require.NoError(t, err)
	defer os.Chmod(outputDir, 0o755)

	cfg := &config.Config{
		Generate: []config.GenerateConfig{
			{
				Name:     "output.txt",
				Path:     outputDir,
				Strategy: "append",
				Inputs: []config.InputFile{
					{Name: "input", Path: inputFile},
				},
			},
		},
	}

	ep, err := NewEntryPoint(cfg)
	require.NoError(t, err)
	defer ep.Close()

	assert.NotNil(t, ep)
}

func TestEntryPointWithNoProcess(t *testing.T) {
	testDir := t.TempDir()
	defer os.RemoveAll(testDir)

	inputFile := filepath.Join(testDir, "input.txt")
	err := os.WriteFile(inputFile, []byte("test content"), 0o644)
	require.NoError(t, err)

	cfg := &config.Config{
		Generate: []config.GenerateConfig{
			{
				Name:     "output.txt",
				Path:     testDir,
				Strategy: "append",
				Inputs: []config.InputFile{
					{Name: "input", Path: inputFile},
				},
			},
		},
		// No Process config
	}

	ep, err := NewEntryPoint(cfg)
	require.NoError(t, err)
	defer ep.Close()

	// Try to start a process - should not crash
	ep.StartManagedProcess()

	// Try to reload - should not crash
	ep.reloadManagedProcess()

	// Verify EntryPoint has no managed command
	assert.Nil(t, ep.managedCmd)
}
