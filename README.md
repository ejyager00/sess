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
