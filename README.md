# ruleengine

[![Build](https://github.com/niharrathod/ruleengine/actions/workflows/build.yml/badge.svg?branch=master)](https://github.com/niharrathod/ruleengine/actions/workflows/build.yml)
[![Test](https://github.com/niharrathod/ruleengine/actions/workflows/test.yml/badge.svg?branch=master)](https://github.com/niharrathod/ruleengine/actions/workflows/test.yml)
[![Coverage Status](https://coveralls.io/repos/github/niharrathod/ruleengine/badge.svg?branch=master)](https://coveralls.io/github/niharrathod/ruleengine?branch=master)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](https://opensource.org/licenses/MIT)

## Overview

ruleengine is a web interface over [ruleengine-core](https://github.com/niharrathod/ruleengine-core). It provides management layer for RuleEngine, which include following operations :

- CRUD RuleEngine APIs
- RuleEngine evaluate API
- Enable/Disable RuleEngine
- RuleEngine versioning

### Todo

- [X] YML Configuration file
- [X] [gin](https://github.com/gin-gonic/gin) based web server
- [X] [Zap](https://pkg.go.dev/go.uber.org/zap) based logger
- [X] docker image build
- [ ] mongodb client setup
- [ ] request tracing
- [ ] RuleEngine CRUD API
- [ ] RuleEngine Enable/Disable API
- [ ] RuleEngine versioning

#### Commands

##### build executable

```bash
    $ env GOOS=linux GOARCH=amd64 go build -o ruleengine .         # for linux
```

##### Start server

```bash
    $ ruleengine -config=config.yml
```

##### To build docker image

```bash
    $ docker image build --no-cache --rm -t <appName>:<tag> .
```
