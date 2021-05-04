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

  * [List Secrets](#list-secrets)

  * [List Images](#list-images)

* [Get Commands](#get-commands)

  * [Get Container Summary](#get-container-summary)

  * [Get Container Details](#get-container-details)

  * [Get Container Tasks](#get-container-tasks)

  * [Get Container Events](#get-container-events)

  * [Get Server Summary](#get-server-summary)

  * [Get Server Details](#get-server-details)

  * [Get S3 Storage Summary](#get-s3-storage-summary)

  * [Get S3 Storage Details](#get-s3-storage-details)

  * [Get Website Summary](#get-website-summary)

  * [Get Website Details](#get-website-details)

* [Update Commands](#update-commands)

  * [Redeploy a Container](#redeploy-a-container)

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
  update      Update a resource

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

## List Commands

The `list` subcommands allows listing of spaces and resources within one or more spaces.

### List Spaces

Example:

```bash
spinup list spaces
```

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

```bash
spinup list resources funSpace sensitive
```

```json
[
  {
    "admin": "s_service_user",
    "created_at": "2020-01-09 15:54:24",
    "flavor": "s3-website",
    "id": 1940,
    "is_a": "storage",
    "name": "spintst000794-mytest.somesite.org",
    "space_name": "funSpace",
    "status": "created",
    "type_name": "S3 Bucket",
    "updated_at": "2020-03-13 12:56:33"
  },
  {
    "admin": "s_service_user",
    "created_at": "2020-04-17 15:06:49",
    "flavor": "linux",
    "id": 2072,
    "ip": "10.12.34.56",
    "name": "spintst-000818.spinup.yale.edu",
    "server_id": "i-0ac04fb882ad31xxxx",
    "space_name": "funSpace",
    "status": "created",
    "type_name": "Centos 8",
    "task": "c50b1fec-db75-40da-9ab6-13dbc280de99",
    "updated_at": "2020-04-17 15:07:54"
  },
  {
    "admin": "s_service_user",
    "created_at": "2020-04-21 13:51:42",
    "flavor": "fargate",
    "id": 2073,
    "name": "spintst-000819-testContainer",
    "space_name": "sensitive",
    "status": "created",
    "type_name": "Container Service",
    "updated_at": "2020-04-21 13:51:44"
  }
]
```

### List Secrets

```bash
spinup list secrets funSpace sensitive
```

```json
[
  {
    "name": "MySecretString",
    "description": "beer",
    "space_name": "funSpace"
  }
]
```

### List Images

```bash
spinup list images funSpace sensitive
```

```json
[
  {
    "architecture": "x86_64",
    "created_at": "2020/05/05 15:25:06",
    "created_by": "s_service_user",
    "description": "My Perfect Image. So nice.",
    "id": "ami-008f837464770f445",
    "name": "Perfect",
    "server_name": "spintst-123456.spinup.yale.edu",
    "state": "available",
    "status": "created",
    "volumes": {
      "/dev/sda1": {
        "delete_on_termination": true,
        "encrypted": false,
        "snapshot_id": "snap-00001111122223334444",
        "volume_size": 30,
        "volume_type": "gp2"
      }
    },
    "offering_name": "CentOS 7",
    "space_name": "funSpace"
  }
]
```

## Get Commands

The `get` subcommands allow you to get detailed information about a resource in a space.

### Get Container Summary

```bash
spinup get funSpace spintst-000848-testService
```

```json
{
  "id": "2120",
  "name": "spintst-000848-testService",
  "status": "created",
  "type": "Container Service",
  "flavor": "fargate",
  "security": "low",
  "space_id": "128",
  "beta": false,
  "size": "nano.512MB",
  "tryit": false,
  "state": "ACTIVE"
}
```

### Get Container Details

```bash
spinup get funSpace spintst-000848-testService -d
```

```json
{
  "id": "2120",
  "name": "spintst-000848-testService",
  "status": "created",
  "type": "Container Service",
  "flavor": "fargate",
  "security": "low",
  "beta": false,
  "size": "nano.512MB",
  "tryit": false,
  "state": "ACTIVE",
  "details": {
    "containers": [
      {
        "auth": false,
        "image": "yalespinup/nginxproxy",
        "name": "nginxproxy",
        "env": {
          "BACKEND_URL": "http://127.0.0.1:8080"
        },
        "portMappings": [
          "8443/tcp"
        ]
      },
      {
        "auth": true,
        "image": "yalespinup/testapi",
        "name": "api",
        "env": {
          "FOOFOOFOO": "kJBDGKLBEGLKWBGLsndlkFNFGLKEN",
          "dKJGBLSGB": "LDNGLWK"
        },
        "healthcheck": {
          "command": [
            "CMD-SHELL",
            "wget -qO- localhost:8080/v1/test/ping || exit 1"
          ],
          "interval": 30,
          "retries": 3,
          "startperiod": 0,
          "timeout": 5
        },
        "mountpoints": [
          {
            "containerpath": "/foobar",
            "readonly": false,
            "sourcevolume": "spintst-000a16-TestFS"
          }
        ],
        "portMappings": [
          "8080/tcp"
        ],
        "secrets": {
          "DERPDERP": "MySecretString"
        }
      }
    ],
    "desiredCount": 1,
    "endpoint": "spintst-000848-testService.svc.spinup.yale.edu",
    "pendingCount": 1,
    "runningCount": 0,
    "spot": false,
    "volumes": [
      {
        "name": "spintst-000a16-TestFS",
        "type": "persistent",
        "nfs_volume": "fs-6ec3009b"
      }
    ]
  }
}
```

### Get Container Tasks

```bash
spinup get funSpace spintst-000848-testService --tasks
```

```json
{
  "tasks": [
    {
      "availabilityZone": "us-east-1a",
      "capacityProvider": "FARGATE_SPOT",
      "cpu": "256",
      "createdAt": "2020-05-21T20:01:58Z",
      "id": "40caf80571634d8db8116c9ee070e5a0",
      "ipAddress": "10.1.2.3",
      "lastStatus": "RUNNING",
      "launchType": "FARGATE",
      "memory": "512",
      "platformVersion": "1.3.0",
      "pullStartedAt": "2020-05-21T20:02:18Z",
      "pullStoppedAt": "2020-05-21T20:02:27Z",
      "stopCode": "",
      "stoppedAt": "",
      "stoppedReason": "",
      "stoppingAt": "",
      "containers": [
        {
          "exitCode": "",
          "healthStatus": "UNKNOWN",
          "image": "yalespinup/testapi",
          "lastStatus": "RUNNING",
          "name": "api",
          "reason": ""
        },
        {
          "exitCode": "",
          "healthStatus": "UNKNOWN",
          "image": "yalespinup/nginxproxy",
          "lastStatus": "RUNNING",
          "name": "nginxproxy",
          "reason": ""
        }
      ],
      "version": 4
    }
  ]
}
```

### Get Container Events

```bash
spinup get funSpace spintst-000848-testService --events
```

```json
[
  {
    "createdAt": "2020-05-21T20:03:14Z",
    "id": "0755e819-5f2b-4ded-aaf8-90de7d215a63",
    "message": "(service spintst-000848-testService) has reached a steady state."
  },
  {
    "createdAt": "2020-05-21T20:01:58Z",
    "id": "75d25a93-9cb3-431b-a248-ca306e37d0f4",
    "message": "(service spintst-000848-testService) has started 1 tasks: (task 40caf80571634d8db8116c9ee070e5a0)."
  },
  {
    "createdAt": "2020-05-21T20:01:57Z",
    "id": "eac34bd9-7b6d-432c-8a49-edaaba41df0b",
    "message": "(service spintst-000848-testService) stopped 1 pending tasks."
  },
  {
    "createdAt": "2020-05-21T20:01:29Z",
    "id": "7aa71f43-7427-4252-a811-79a574098302",
    "message": "(service spintst-000848-testService) has started 1 tasks: (task b1de7eec4b6845b79ac16c7e8786a8c1)."
  },
  {
    "createdAt": "2020-05-21T19:53:36Z",
    "id": "97ba8254-85ca-40d0-8898-a7860a07b190",
    "message": "(service spintst-000848-testService) has reached a steady state."
  },
  {
    "createdAt": "2020-05-21T19:53:25Z",
    "id": "e1e405b7-604c-40d2-bfca-447db95bd2a4",
    "message": "(service spintst-000848-testService) has stopped 2 running tasks: (task cc962c70fa74424fa711c7e88d78c21a) (task 72c70c1a70d647909df56f197aebe3e4)."
  }
]
```

### Get Server Summary

```bash
spinup get funSpace spintst-000818.spinup.yale.edu
```

```json
TODO
```

### Get Server Details

```bash
spinup get funSpace spintst-000818.spinup.yale.edu -d
```

```json
TODO
```

### Get S3 Storage Summary

```bash
spinup get funSpace spintst000794-mytest.somesite.org
```

```json
TODO
```

### Get S3 Storage Details

```bash
spinup get funSpace spintst000794-mytest.somesite.org -d
```

```json
TODO
```

## Update Commands

The `update` subcommands allow you to make changes to an existing resource.

### Redeploy a Container

To redeploy an existing container service, using the existing configuration and tag.  This will force the latest image with the defined tag to be pulled and redeployed.  Container (re)deployments are rolling.  This is useful if you have a tag that gets updated with the latest release and you want deploy that via an automated pipeline.

```bash
spinup update funSpace spintst-000848-testService -r
```

```json
OK
```

## Author

E Camden Fisher <camden.fisher@yale.edu>

## License

GNU Affero General Public License v3.0 (GNU AGPLv3)
Copyright (c) 2021 Yale University
