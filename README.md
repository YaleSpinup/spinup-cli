# spinup

The spinup user interface in CLI form.  The interface takes on a git-like style and outputs JSON for easy consumption by computers and humans with computer brains.  To make the output more friendly to the masses, consider using [jq](https://stedolan.github.io/jq/).

## Table of Contents

* [Getting Started](#getting-started)

  * [Download](#download)

  * [Running Spinup Cli](#running-the-command)

* [Configuration](#configuration)

* [List Commands](#list-commands)

  * [List Spaces](#list-spaces)

  * [List Resources](#list-resources)

* [Get Commands](#get-commands)

* [Author](#author)

* [License](#license)

## Getting Started

Spinup is a cross-compiled static binary with support for many platforms.  Download an install the relevant binary for your system.

### Download

* TODO
* TODO

### Running the command

```
# spinup help
A small CLI for interacting with Yale's Spinup service

Usage:
  spinup [command]

Available Commands:
  get         Get details about a resource
  help        Help about any command
  list        List spinup objects
  new         Create new resources

Flags:
      --config string     config file (default is $HOME/.spinup.yaml)
      --debug             Enable debug logging
  -h, --help              help for spinup
  -p, --password string   Spinup password
  -s, --spaces strings    Space ID
      --url string        The base url for Spinup
  -u, --username string   Spinup username
  -v, --verbose           Enable verbose logging

Use "spinup [command] --help" for more information about a command.
```

## Configuration

By default the configuration lives in ~/.spinup.{yml|json}.

All fields in the configuration can be overridden on the command line.  All fields in the configuration file are optional
and will act as defaults.  Most users will probably want the `username`, `password` and `url` configured.

Supported configuration items:

| property | type         | description               |
|:---------|:------------:|:-------------------------:|
| username | string       | spinup user               |
| password | string       | spinup password           |
| url      | string       | spinup url                |
| spaces   | number array | default list of space ids |

Example `~/.spinup.json`:

```json
{
  "username": "s_service_user",
  "password": "secretPassword",
  "url": "https://spinup.example.edu",
  "spaces": ["11111", "22222"]
}
```

## List Commands

The `list` subcommands allows listing of spaces and resources within one or more spaces.

### List Spaces

Example:

`spinup list spaces`

```json
[
  {
    "id": 128,
    "name": "funSpace",
    "owner": "someUser",
    "security": "low",
    "created_at": "2018-11-15 18:47:50",
    "mine": true
  },
  {
    "id": 136,
    "name": "sensitive",
    "owner": "otherUser",
    "security": "moderate",
    "created_at": "2019-04-10 08:43:58",
    "mine": true
  }
]
```

### List Resources

Resources can be listed for one or more spaces by space name or space ID.

Example (the following are equivalent):

`spinup list resources funSpace sensitive`

`spinup list resources -s 128 -s 136`

```json
[
  {
    "admin": "s_service_user",
    "created_at": "2020-01-09 15:54:24",
    "flavor": "s3-website",
    "id": 1940,
    "is_a": "storage",
    "name": "spintst000794-mytest.somesite.org",
    "size_id": 117,
    "space_id": 128,
    "status": "created",
    "type_id": 38,
    "updated_at": "2020-03-13 12:56:33"
  },
  {
    "admin": "s_service_user",
    "created_at": "2020-04-17 15:06:49",
    "flavor": "linux",
    "id": 2072,
    "ip": "10.12.34.56",
    "is_a": "server",
    "name": "spintst-000818.spinup.yale.edu",
    "server_id": "i-0ac04fb882ad31xxxx",
    "size_id": 136,
    "space_id": 128,
    "status": "created",
    "type_id": 18,
    "task": "c50b1fec-db75-40da-9ab6-13dbc280de99",
    "updated_at": "2020-04-17 15:07:54"
  },
  {
    "admin": "s_service_user",
    "created_at": "2020-04-21 13:51:42",
    "flavor": "fargate",
    "id": 2073,
    "is_a": "container",
    "name": "spintst-000819-testContainer",
    "size_id": 97,
    "space_id": 136,
    "status": "created",
    "type_id": 24,
    "updated_at": "2020-04-21 13:51:44"
  }
]
```

## Get Commands

The `get` subcommands allow you to get detailed information about a resource.

### TODO

## Author

E Camden Fisher <camden.fisher@yale.edu>

## License

GNU Affero General Public License v3.0 (GNU AGPLv3)
Copyright (c) 2020 Yale University
