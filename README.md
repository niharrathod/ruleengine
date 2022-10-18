# ruleengine

[![Build](https://github.com/niharrathod/ruleengine/actions/workflows/build.yml/badge.svg?branch=master)](https://github.com/niharrathod/ruleengine/actions/workflows/build.yml)
[![Test](https://github.com/niharrathod/ruleengine/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/niharrathod/ruleengine/actions/workflows/test.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

## Overview

ruleengine is a web interface over [ruleengine-core](https://github.com/niharrathod/ruleengine-core) for control plane and data plane operations, which include following :

- CRUD RuleEngine (Control plane) APIs
- RuleEngine versioning
- RuleEngine evaluate (Data plane) API

### Todo

- [X] YML Configuration file
- [X] [gin](https://github.com/gin-gonic/gin) based web server
- [X] [Zap](https://pkg.go.dev/go.uber.org/zap) based logger
- [X] docker image build
- [X] mongodb client setup
- [ ] request tracing

##### Control plane API

- [X] RuleEngine CRD API
- [X] RuleEngine Update (Default/Enable/Disable) API
- [X] RuleEngine versioning

##### Data plane API

- [ ] RuleEngine evaluate

### Commands

```bash
# Build executable
env GOOS=linux GOARCH=amd64 go build -o ruleengine .         # linux

# Start server
ruleengine -config=config.yml

# Build docker image
docker image build --no-cache --rm -t <appName>:<tag> .
```