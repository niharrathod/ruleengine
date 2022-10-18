package entities

import (
	"fmt"

	ruleenginecore "github.com/niharrathod/ruleengine-core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RuleEngine struct {
	Name           string          `bson:"name"`
	DefaultTag     string          `bson:"defaultTag"`
	Tags           map[string]*Tag `bson:"tags"`
	LastUpdateTime int64           `bson:"lastUpdateTime"`
}

type Tag struct {
	Name           string             `bson:"name"`
	EngineConfigID primitive.ObjectID `bson:"engineConfigId"`
	IsEnable       bool               `bson:"isEnable"`
}

type EngineConfig struct {
	ID               primitive.ObjectID               `bson:"_id"`
	EngineCoreConfig *ruleenginecore.RuleEngineConfig `bson:"config"`
}

type CompleteRuleEngine struct {
	Name       string                  `json:"name"`
	DefaultTag string                  `json:"defaultTag"`
	Tags       map[string]*TagResponse `json:"tags"`
}

type TagResponse struct {
	IsEnable bool                             `json:"isEnable"`
	Config   *ruleenginecore.RuleEngineConfig `json:"config"`
}

type Error struct {
	ErrCode  uint   `json:"errCode"`
	ErrMsg   string `json:"errMsg"`
	OtherMsg string `json:"otherMsg"`
}

func NewError(errCode uint) *Error {
	errMsg, ok := errCodeToMessage[errCode]
	if !ok {
		errMsg = "UnknownFailure"
	}

	return &Error{
		ErrCode: errCode,
		ErrMsg:  errMsg,
	}
}

func NewErrorWithMsg(errCode uint, otherMsg string) *Error {
	errMsg, ok := errCodeToMessage[errCode]
	if !ok {
		errMsg = "UnknownFailure"
	}

	return &Error{
		ErrCode:  errCode,
		ErrMsg:   errMsg,
		OtherMsg: otherMsg,
	}
}

func (err *Error) Error() string {
	return fmt.Sprintf("errCode:%v errMsg:%v. %v", err.ErrCode, err.ErrMsg, err.OtherMsg)
}

const (
	ErrCodeParsingFailed                   = 1
	ErrCodeInvalidRuleEngineName           = 2
	ErrCodeInvalidTagName                  = 3
	ErrCodeInvalidRuleEngineConfig         = 4
	ErrCodeDatastoreFailed                 = 5
	ErrCodeRuleEngineNotFound              = 6
	ErrCodeTagNotFound                     = 7
	ErrCodeTagDeleteNotAllowed             = 8
	ErrCodeTagDisableNotAllowed            = 9
	ErrCodeDefaultTagExistAndMustBeEnabled = 10
	ErrCodeTagAlreadyExist                 = 11
	ErrCodeEvaluationFailed                = 12
)

var errCodeToMessage = map[uint]string{
	ErrCodeParsingFailed:                   "Parsing failed",
	ErrCodeInvalidRuleEngineName:           "Invalid ruleEngineName. alphabetic([a-z][A-Z]) and maximum 30 characters allowed",
	ErrCodeInvalidTagName:                  "Invalid tag. alphabetic([a-z][A-Z]) and maximum 30 characters allowed",
	ErrCodeInvalidRuleEngineConfig:         "RuleEngineConfig is invalid",
	ErrCodeDatastoreFailed:                 "Internal datastore failure",
	ErrCodeRuleEngineNotFound:              "RuleEngine not found",
	ErrCodeTagNotFound:                     "Tag not found",
	ErrCodeTagDeleteNotAllowed:             "Could not delete tag, either set as default or enabled",
	ErrCodeTagDisableNotAllowed:            "Could not disable default tag",
	ErrCodeDefaultTagExistAndMustBeEnabled: "Could not set defaultTag, either not found or not enabled",
	ErrCodeTagAlreadyExist:                 "Tag already exist",
}
