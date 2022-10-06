# Overview

Defines goals and high-level design for RuleEngine.

## Goal

- Support for RuleEngine control plane operations
  - Create/Delete/Read RuleEngine operation
  - RuleEngine versioning (Tagging), i.e. docker like approach where string tag as version
    - defaultTag for RuleEngine,  a default version of ruleEngine for date plane operation
  - Enable/Disable RuleEngine(identified with rulename and tag)
  - Allow more than one RuleEngine enabled

- Support for RuleEngine data plane operations
  - RuleEngine evaluate operation (with optional tag, i.e. default tagged)

## High-level Design (WIP)

RuleEngine for given tag would follow below state diagram:

[![](https://mermaid.ink/img/pako:eNplkDEOwjAMRa8SeUR0YezAQjkBG4TBEBcqUrtKkyKEuDtp2lIQnvL_s50vP-EshiCH1qOnosKLwzrrVppVrMPiqLJsrTaOIjWDOYoEtownO4FRJFBU7ReZ1P_MDynI0uebUSQQcwwmiydlqfRKyt9FKUCHNsRsShpy6CthhdbKfeohNmnDoDTPk33BEmpyNVYmnuPZMw3-SjVpyOPToLtp0PyKfRi87B58hty7QEsIjZmvB3mJto1ug7wXmfTrDcaAcNg)](https://mermaid.live/edit#pako:eNplkDEOwjAMRa8SeUR0YezAQjkBG4TBEBcqUrtKkyKEuDtp2lIQnvL_s50vP-EshiCH1qOnosKLwzrrVppVrMPiqLJsrTaOIjWDOYoEtownO4FRJFBU7ReZ1P_MDynI0uebUSQQcwwmiydlqfRKyt9FKUCHNsRsShpy6CthhdbKfeohNmnDoDTPk33BEmpyNVYmnuPZMw3-SjVpyOPToLtp0PyKfRi87B58hty7QEsIjZmvB3mJto1ug7wXmfTrDcaAcNg)

#### Control Plane operations

- CreateRuleEngine operation
  - mandatory : [ruleEngineName, Tag, ruleEngineConfig]
  - if RuleEngine exist with ruleEngineName, i.e. introducing new tag(version) to existing ruleEngine
    - validate Tag : valid value and uniqueness
    - validate RuleEngineConfig
    - Persist RuleEngineConfig details
  - else i.e. adding new ruleEngine with tag
    - validate ruleEngineName: valid value and uniqueness
    - validate Tag: valid value
    - validate RuleEngineConfig
    - set defaultTag as empty and enabledTag list as empty
    - Persist RuleEngineConfig details

- DeleteRuleEngine operation
  - mandatory : [ruleEngineName, Tag]
  - validate Tag is not enable. i.e. delete is allowed for disabled tags only
  - mark isDelete Flag true for given ruleEngineName and Tag

- ReadRuleEngine operation
  - mandatory :  [ruleEngineName]
  - optional : [tag]
  - fetch(based on tag as filter) RuleEngine details and return

- SetDefaultTag operation
  - mandatory : [ruleEngineName, tag]
  - validate ruleEngineName and tag exist and given tag must be enable
  - update defaultTag

- RemoveDefaultTag operation
  - mandatory : [ruleEngineName]
  - validate ruleEngineName exist
  - update defaultTag as empty

- Enable operation
  - mandatory : [ruleEngineName, tag]
  - validate ruleEngineName and tag exist, not deleted and given tag must be disabled
  - update enableTags

- Disable operation
  - mandatory : [ruleEngineName, tag]
  - validate ruleEngineName and tag exist and given tag must be enabled
  - update enableTags

- Background worker
  - a background worker for every RuleEngine
  - background worker would observe for following:
    1. enable operation for RuleEngineName:Tag
        - creates new RuleEngine instance and be ready for data plane operation
    2. disable operation for RuleEngineName:Tag  
        - discard existing ruleEngine instance
    3. change in defaultTag
        - update in-memory defaultTag, i.e. route default tag data plane operation to new defaultTag

#### Data Plane operations

- Evaluate
  - mandatory : [ruleEngineName, RuleEngineInput, opType(Complete,AscPriority,DscPriority), limit(only for AscPriority,DscPriority opType)]
  - Optional : [Tag, ruleName]
    - validate existence of ruleengine with ruleEngineName
    - if Tag is not provided, consider defaultTag for operation
    - Evaluate operation based on options

### Data model

RuleEngine uses document store(mongoDB) for persistance layer. Data model is divided in two Entities:

  1. RuleEngine
    - RuleEngineName  string  (unique) (indexed)
    - DefaultTag      string
    - EnableTags      listOfString
  2. Tag
    - RuleEngineName    string (unique) (indexed)
    - Tag               string (unique) (indexed)
    - Config            RuleEngineConfig
    - createdDate       long
    - isDeleted         boolean

### API reference

ToDo
