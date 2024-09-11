# spinup

The spinup user interface in CLI form.  The interface takes on a git-like style and outputs JSON for easy consumption by computers and humans with computer brains.  To make the output more friendly to the masses, consider using [jq](https://stedolan.github.io/jq/).

---
[![goreleaser](https://github.com/YaleSpinup/spinup-cli/actions/workflows/releaser.yml/badge.svg)](https://github.com/YaleSpinup/spinup-cli/actions/workflows/releaser.yml)
[![tests](https://github.com/YaleSpinup/spinup-cli/actions/workflows/tests.yaml/badge.svg)](https://github.com/YaleSpinup/spinup-cli/actions/workflows/tests.yaml)
## Table of Contents

- [spinup](#spinup)
  - [Table of Contents](#table-of-contents)
  - [Getting Started](#getting-started)
    - [Download](#download)
    - [Running the command](#running-the-command)
  - [Configuration](#configuration)
    - [Configure with the configuration utility](#configure-with-the-configuration-utility)
  - [Get Commands](#get-commands)
  - [Update Commands](#update-commands)
    - [Containers](#containers)
      - [Redeploy](#redeploy)
      - [Scale](#scale)
  - [Author](#author)
  - [License](#license)

## Getting Started

Spinup is a cross-compiled static binary with support for many platforms.  Download an install the relevant binary for your system.

### Installation

#### MacOS
`spinup-cli` is available through Homebrew.

```sh
brew install yalespinup/tools/spinup
```

See https://github.com/YaleSpinup/homebrew-tools for more Spinup related CLI tools.

#### Linux
Download the precomplied and compressed binary from the [releases page](https://github.com/YaleSpinup/spinup-cli/releases). Decompress the file and move it to your preferred installation directory.

```sh
wget https://github.com/YaleSpinup/spinup-cli/releases/download/v0.4.10/spinup-cli_0.4.10_linux_amd64.tar.gz
tar -xz -f spinup-cli_0.4.10_linux_amd64.tar.gz
mv spinup /usr/local/bin/
sudo chown root:root /usr/local/bin/spinup
```

#### Windows
- TODO

### Running the command

```bash
# spinup help
A small CLI for interacting with Yale's Spinup service

Usage:
  spinup [command]

Available Commands:
  completion  generate the autocompletion script for the specified shell
  configure   Configure Spinup CLI
  get         Get information about a resource in a space
  help        Help about any command
  new         Create new resources
  update      Update a resource in a space
  version     Display version information about the spinup-cli.

Flags:
      --config string    config file (default is $HOME/.spinup.yaml)
      --debug            Enable debug logging
  -h, --help             help for spinup
  -s, --spaces strings   Default Space(s)
  -t, --token string     Spinup API Token
      --url string       The base url for Spinup
  -v, --verbose          Enable verbose logging

Use "spinup [command] --help" for more information about a command.
```

## Configuration

By default the configuration lives in ~/.spinup.{yml|json}.

All fields in the configuration can be overridden on the command line.  All fields in the configuration file are optional
and will act as defaults.  Most users will probably want the `url`, and `token`.

Supported configuration items:

| property | type         | description                 |
|:---------|:------------:|:---------------------------:|
| url      | string       | spinup url                  |
| token    | string       | spinup token                |
| spaces   | string array | default list of space names |

Example `~/.spinup.json`:

```json
{
  "url": "https://spinup.example.edu",
  "token": "xxxxxyyyyyyy",
  "spaces": ["my_space_1", "my_space_2"]
}
```

### Configure with the configuration utility

```bash
spinup configure
```

## Get Commands

The `get` subcommands allow you to get detailed information about spinup resources.

```bash
# spinup get --help
Get information about a resource in a space

Usage:
  spinup get [command]

Available Commands:
  container   Get a container service
  database    Get a container service
  images      Get a list of images for a space
  secrets     Get a list of secrets for a space
  server      Get a server service
  space       Get details about your space(s)
  spaces      Get a list of your space(s)
  storage     Get a storage service

Flags:
  -d, --details   Get detailed output about the resource
  -h, --help      help for get

Global Flags:
      --config string    config file (default is $HOME/.spinup.yaml)
      --debug            Enable debug logging
  -s, --spaces strings   Default Space(s)
  -t, --token string     Spinup API Token
      --url string       The base url for Spinup
  -v, --verbose          Enable verbose logging

Use "spinup get [command] --help" for more information about a command.
```
## Update Commands

The `update` subcommands allow you to make changes to an existing resource.  Currently only container updates are supported.

```bash
# spinup update --help
Update a resource in a space

Usage:
  spinup update [command]

Available Commands:
  container   Update a container service

Flags:
  -h, --help   help for update

Global Flags:
      --config string    config file (default is $HOME/.spinup.yaml)
      --debug            Enable debug logging
  -s, --spaces strings   Default Space(s)
  -t, --token string     Spinup API Token
      --url string       The base url for Spinup
  -v, --verbose          Enable verbose logging

Use "spinup update [command] --help" for more information about a command.
```

### Containers

#### Redeploy

Redeploy an existing container service, using the existing configuration and tag.  This will force the latest image with the defined tag to be pulled and redeployed.  Container (re)deployments are rolling.  This is useful if you have a tag that gets updated with the latest release and you want deploy that via an automated pipeline.

```bash
spinup update funSpace/spintst-000848-testService -r
```

```json
OK
```

#### Scale

Scale an existing container service, using the existing configuration and tag.  This will set the desired count to the passed value.  Values between 0 and 10 are supported.

```bash
spinup update container funSpace spintst-000848-testService --scale 2
```

```json
OK
```

## Author

* E Camden Fisher <camden.fisher@yale.edu>
* Brandon Tassone <brandon.tassone@yale.edu>

## License

GNU Affero General Public License v3.0 (GNU AGPLv3)
Copyright (c) 2021 Yale University
