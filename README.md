# Synchronize Enivornments Simply, Stupid (SESS)

SESS is a command line tool for synchronizing development environments.

## Building

`make build` will build the project and put the binary in `target/sess`.

## Usage

```
sess capture
```
Captures the current relevant environment state and stores it in a configuration YAML file.

```
sess install [file]
```
Install the tools and extensions specified in the configuration file.

```
sess validate [file]
```
Validate that the current environment matches the provided configuration file.

## Configuration

SESS uses YAML files to define environment configurations. Here's an example:

```yaml
tools:
  - name: node
    version: ">= 21.x"
    install_command: "nvm install 21"
    version_check_command: "node --version"

extensions:
  - id: "denoland.vscode-deno"
    ide: "vscode"
    version: "latest"

env_variables:
  JAVA_HOME: "/usr/lib/jvm/java-11"
  GOPATH: "/home/user/go"

dotfiles:
  - ".gitconfig"
  - ".bashrc"
```

### Configuration Sections

- **tools**: Define development tools with their version requirements and installation commands
  - `name`: Tool identifier
  - `version`: Version constraint (supports semver ranges)
  - `install_command`: Command to install the tool
  - `version_check_command`: Command to verify the installed version

- **extensions**: IDE extensions to install
  - `id`: Extension identifier
  - `ide`: IDE type (currently supports "vscode")
  - `version`: Desired version ("latest" or specific version)

- **env_variables**: Environment variables to set
  - Key-value pairs of environment variable names and their values

- **dotfiles**: List of dotfiles to track and synchronize
