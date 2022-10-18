# Overview

Define goals and high-level design for RuleEngine.

## Goal (WIP)

- Support for control plane operations
  - Create/Update/Delete/Read RuleEngine
  - RuleEngine versioning (Tagging), i.e. docker like approach where string tag as version
    - Allow multiple RuleEngine versions
    - Default version for data plane operation
  - Enable specific RuleEngine version
  - Disable specific RuleEngine version

- Support for data plane operations
  - RuleEngine evaluate
    - proxy to RuleEngineCore Evaluate operation
    - Optional version as an input, if version is not provided consider a default version.

## High-level Design (WIP)

A RuleEngine is defined as name, default tag(default version) and tag list(list of versions). Every Tag defines as tag name, enable/disable flag and RuleEngineConfig(config defined as RuleEngine-Core, used to instantiate RuleEngine). RuleEngine tags follows docker style versioning. A RuleEngine can have multiple tags but at a time only one tag would act as default for data plane operation. Once RuleEngine with a Tag is created, enabling a tag prepares that specific versioned RuleEngine ready for data plane operations, internally it creates RuleEngine instance. Disabled Tag i.e. removes RuleEngine instance but maintains RuleEngineConfig with tag as label for future.

Evaluate as Data plane operation is proxy to RuleEngineCore Evaluate operation, ruleEngineName is required to identify specific instance, tag is optional. if tag is not provided default tagged ruleEngine is picked for evaluation.

**RuleEngine with tag would follow below state diagram:**

[![](https://mermaid.ink/img/pako:eNplkDEOwjAMRa8SeUR0YezAQjkBG4TBEBcqUrtKkyKEuDtp2lIQnvL_s50vP-EshiCH1qOnosKLwzrrVppVrMPiqLJsrTaOIjWDOYoEtownO4FRJFBU7ReZ1P_MDynI0uebUSQQcwwmiydlqfRKyt9FKUCHNsRsShpy6CthhdbKfeohNmnDoDTPk33BEmpyNVYmnuPZMw3-SjVpyOPToLtp0PyKfRi87B58hty7QEsIjZmvB3mJto1ug7wXmfTrDcaAcNg)](https://mermaid.live/edit#pako:eNplkDEOwjAMRa8SeUR0YezAQjkBG4TBEBcqUrtKkyKEuDtp2lIQnvL_s50vP-EshiCH1qOnosKLwzrrVppVrMPiqLJsrTaOIjWDOYoEtownO4FRJFBU7ReZ1P_MDynI0uebUSQQcwwmiydlqfRKyt9FKUCHNsRsShpy6CthhdbKfeohNmnDoDTPk33BEmpyNVYmnuPZMw3-SjVpyOPToLtp0PyKfRi87B58hty7QEsIjZmvB3mJto1ug7wXmfTrDcaAcNg)

#### Control Plane operations

- Create RuleEngine operation
  - mandatory : [ruleEngineName, Tag, ruleEngineConfig]
  - if RuleEngine exist with ruleEngineName, i.e. introducing new tag(version) to existing ruleEngine
    - validate Tag : valid value and uniqueness
    - validate RuleEngineConfig
    - Persist Tag with RuleEngineConfig
  - else i.e. adding new ruleEngine with tag
    - validate ruleEngineName: valid value and uniqueness
    - validate Tag: valid value and uniqueness
    - validate RuleEngineConfig
    - Persist RuleEngine, Tag with RuleEngineConfig

- Delete RuleEngine operation
  - mandatory : [ruleEngineName]
  - Remove RuleEngine and all versions from persistance layer

- Delete RuleEngine having specific tag
  - mandatory : [ruleEngineName, tag]
  - validate: tag is not enabled and not set as default
  - Remove RuleEngine and all versions from persistance layer

- Read RuleEngine operation
  - mandatory :  [ruleEngineName]
  - fetch details from persistance layer

- SetDefaultTag operation
  - mandatory : [ruleEngineName, tag]
  - validate : tag is enabled
  - Update RuleEngine entity with default tag

- RemoveDefaultTag operation
  - mandatory : [ruleEngineName]
  - Update RuleEngine entity with default tag as none

- Enable operation
  - mandatory : [ruleEngineName, tag]
  - Update RuleEngine entity with tag enable flag as true

- Disable operation
  - mandatory : [ruleEngineName, tag]
  - Update RuleEngine entity with tag enable flag as true

- Background worker
  - Background worker would observe changes in persistence layer and acts as following:
    1. enable operation for RuleEngineName:Tag
        - creates new RuleEngine instance and be ready for data plane operation
    2. disable operation for RuleEngineName:Tag  
        - discard existing ruleEngine instance
    3. change in defaultTag
        - update in-memory defaultTag, i.e. route default tag data plane operation to new defaultTag
    4. delete RuleEngine
        - Remove all RuleEngine Instances.

#### Data Plane operations

- Evaluate
  - mandatory : [ruleEngineName, RuleEngineInput, opType(Complete,AscPriority,DscPriority), limit(only for AscPriority,DscPriority opType)]
  - Optional : [Tag, ruleName]
  - validate existence of ruleengine with ruleEngineName
  - if Tag is not provided, consider defaultTag for operation
  - Evaluate operation based on options

### Data model

RuleEngine uses document store(mongoDB) for persistance layer. Data model is divided in two Entities:

**RuleEngine**
    {}
**RuleEngineConfig**
    {}

### API reference

ToDo
