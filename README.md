# Shoehorn

> [!WARNING]
> This project is still under development and some features may not be fully reliable yet.

Shoehorn is a generic container entrypoint capable of configuration file generation and process lifecycle management within containers.
For example, this entrypoint can be used to generate a composite configuration file for an application that requires both sensitive credentials and non-sensitive configuration details, allowing you to manage as much as possible in a gitops-driven way while still keeping sensitive information out of your git repository.
The manner of handling sensitive information is left up to the user, but this entrypoint can be used to generate a composite configuration file that includes both sensitive and non-sensitive information.

It supports two strategies for generating the output file:

1. **Append**: Concatenates multiple input files into a single output file, preserving the order of the input files.
2. **Template**: Uses a template file to generate the output file, allowing for more complex configurations and variable substitution.

## Features

- **File Watching**: Monitors source files for changes
- **Configuration Generation**:
  - Combines multiple input files into a single output
  - Supports simple concatenation (append) or templating
- **Process Management**:
  - Starts another process within the container
  - Forwards stdin and CLI arguments to the managed process
  - Manages process lifecycle (start, stop, reload)
- **Configuration via YAML**: Simple, declarative configuration
- **Minimal Dependencies**: Uses only two external libraries (fsnotify and yaml)

## Configuration

The entrypoint is configured via a YAML file with the following structure:

```yaml
generate:
  - name: my-composite-config.yaml # Name of the output file
    path: /my/output/directory/ # Path for the output file
    strategy: append # Either 'append' or 'template'
    template: /my/template/config.tpl # Used when strategy=template
    inputs: # Input files to watch
      - name: my-config-1 # Template variable name when using strategy=template
        path: /some/config.yml # Path to the input file
      - name: my-credentials-secret
        path: /secrets/credentials/my-credentials
process:
  path: /my/binary/process # Process to manage
  reload:
    enabled: false # Whether to reload on config changes
    method: restart # 'restart' or 'signal'
    signal: SIGHUP # Signal to send when method=signal
  args: [] # Default args for the process
```

## Usage

The entrypoint is designed to replace the original entrypoint of a container.

```dockerfile
FROM your-base-image
COPY entrypoint /entrypoint/bin/entrypoint
COPY config.yaml /entrypoint/config/config.yaml
ENTRYPOINT ["/entrypoint"]
```

Any arguments after the config file will be forwarded to the managed process in addition to the args specified in the config.

## Strategies for File Generation

### Append Strategy

The `append` strategy simply concatenates all input files, preserving their order as specified in the configuration. This is useful for combining configuration files where order matters.

### Template Strategy

The `template` strategy uses Go's text/template package to render a template file. Each input file's content is made available as a variable in the template, using the name specified in the configuration.

Example template:

```
# Server configuration
{{.`server-config`}}

# SSL certificates
{{.`ssl-certs`}}
```

## Process Reload Methods

### Restart Method

The `restart` method completely stops and restarts the managed process when configuration files change.

### Signal Method

The `signal` method sends a specified signal (e.g., SIGHUP) to the managed process, allowing it to reload its configuration without restarting.

## Building

This project is built with `make`. See either `make help` or check the `Makefile` for additional info.

To compile the binary run `make build`. The compiled binary will be at `build/shoehorn`

```bash
make build
```

## Example Use Cases

- Generate Nginx configuration by combining multiple site configs
- Combine Kubernetes manifests from different sources
- Template together database configuration with credentials from a secrets file
- Watch for changes in ConfigMaps and Secrets in Kubernetes deployments
