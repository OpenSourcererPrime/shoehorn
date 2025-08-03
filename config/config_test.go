package config

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var configTests = []struct {
	name           string
	content        string
	expectedConfig *Config
	expectedError  error
}{
	{
		name: "valid config",
		content: `
generate:
  - name: test_file.yml
    path: /etc/
    strategy: append
    inputs:
      - name: test.yml
        path: /sources/

process:
  path: test_process
  reload:
    enabled: true
    method: restart
  args:
    - arg1
    - arg2
`,
		expectedConfig: &Config{
			Generate: []GenerateConfig{
				{
					Name:     "test_file.yml",
					Path:     "/etc/",
					Strategy: "append",
					Inputs: []InputFile{
						{
							Path: "/sources/",
							Name: "test.yml",
						},
					},
				},
			},
			Process: ProcessConfig{
				Path: "test_process",
				Reload: ReloadConfig{
					Enabled: true,
					Method:  "restart",
				},
				Args: []string{"arg1", "arg2"},
			},
		},
		expectedError: nil,
	},
	{
		name: "config with reload via signal method",
		content: `
generate:
  - name: test_file.yml
    path: /etc/
    strategy: append
    inputs:
      - name: test.yml
        path: /sources/

process:
  path: test_process
  reload:
    enabled: true
    method: signal
    signal: SIGHUP
  args:
    - arg1
    - arg2
`,
		expectedConfig: &Config{
			Generate: []GenerateConfig{
				{
					Name:     "test_file.yml",
					Path:     "/etc/",
					Strategy: "append",
					Inputs: []InputFile{
						{
							Path: "/sources/",
							Name: "test.yml",
						},
					},
				},
			},
			Process: ProcessConfig{
				Path: "test_process",
				Reload: ReloadConfig{
					Enabled: true,
					Method:  "signal",
					Signal:  "SIGHUP",
				},
				Args: []string{"arg1", "arg2"},
			},
		},
		expectedError: nil,
	},
	{
		name: "invalid strategy",
		content: `
generate:
  - name: test_file.yml
    path: /etc/
    strategy: non_existent_strategy
    inputs:
      - name: test.yml
        path: /sources/

process:
  path: test_process
  reload:
    enabled: true
    method: restart
  args:
    - arg1
    - arg2
`,
		expectedConfig: nil,
		expectedError:  &ErrorInvalidStrategy{Strategy: "non_existent_strategy", Name: "test_file.yml"},
	},
	{
		name: "missing template path",
		content: `
generate:
  - name: test_file.yml
    path: /etc/
    strategy: template
    inputs:
      - name: test.yml

process:
  path: test_process
  reload:
    enabled: true
    method: restart
  args:
    - arg1
    - arg2
`,
		expectedConfig: nil,
		expectedError:  &ErrorInvalidStrategy{Strategy: "non_existent_strategy", Name: "test_file.yml"},
	},
	{
		name: "invalid reload method",
		content: `
generate:
  - name: test_file.yml
    path: /etc/
    strategy: template
    inputs:
      - name: test.yml
      	path: /sources/

process:
  path: test_process
  reload:
    enabled: true
    method: non_existent_method
  args:
    - arg1
    - arg2
`,
		expectedConfig: nil,
		expectedError:  &ErrorInvalidReloadMethod{Method: "non_existent_method"},
	},
	{
		name: "missing signal for reload method",
		content: `
generate:
  - name: test_file.yml
    path: /etc/
    strategy: template
    inputs:
      - name: test.yml
      	path: /sources/

process:
  path: test_process
  reload:
    enabled: true
    method: signal
  args:
    - arg1
    - arg2
`,
		expectedConfig: nil,
		expectedError:  &ErrorMissingSignal{},
	},
}

func TestLoadConfig(t *testing.T) {
	for _, tc := range configTests {
		t.Run(tc.name, func(t *testing.T) {
			r := strings.NewReader(tc.content)
			appConfig, err := LoadConfig(r)
			if tc.expectedError != nil {
				fmt.Println("Expected error:", tc.expectedError)
				assert.ErrorAs(t, err, &tc.expectedError, "Expected error: %v, got: %v", tc.expectedError, err)
			} else {
				assert.NoError(t, err, "Expected no error, got: %v", err)
			}
			assert.Equal(t, tc.expectedConfig, appConfig, "Expected config: %v, got: %v", tc.expectedConfig, appConfig)
		})
	}
}
