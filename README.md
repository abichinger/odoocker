# Odoocker

Odoocker is a CLI app to run odoo tests. It creates a docker container for PostgreSQL and Odoo. Then it installs the modules and runs the specified tests.

## Requirements

Ensure that docker is installed.

## Install

### Go install

Use `go install` to build and install odoocker.

```sh
go install github.com/abichinger/odoocker
```

### Binary

You can download the latest binary from the [Github releases](https://github.com/abichinger/odoocker/releases) page.

## Usage

```sh
$> odoocker -help
odoocker v1.0.0 - Setup and run odoo inside a docker container

Available commands:

   test   Run odoo tests 
   run    Runs a one-time command inside the odoo container

odoocker test - Run odoo tests
Flags:

  -addons string
        path to odoo addons folder or a single addon (default ".")
  -help
        Get help on the 'odoocker test' command.
  -m value
        List of modules to install
  -pg string
         (default "15")
  -t value
        odoo test-tags
  -tempdir string
        temporary directory (default ".odoocker")
  -term
        is terminal (default true)
  -tour
        installs websocket-client and chrome
  -v string
        Odoo version (default "latest")
  -verbose
        verbose output
```