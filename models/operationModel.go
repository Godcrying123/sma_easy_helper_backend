package models

import (
	"errors"

	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
	"github.com/mitchellh/mapstructure"
)

// Operation struct is for structuring the operation model with below attribute
type Operation struct {
	OperationID          int
	OperationShortName   string
	OperationDescription string
	NumOfSteps           int
	DetailedSteps        []SubOperation
}

// SubOperation struct is for structuring the operation model with below attribute
type SubOperation struct {
	SubOperationID          int
	SubOperationDescription string
	Machine                 string
	StepType                string
	Commands                []string
	File                    string
	FileChange              string
	URL                     string
	Image                   string
	Notifications           string
	CheckMethod             string
	CheckResults            bool
}

// this is the data structure used in Operation
var (
	OperationMap      map[string]Operation
	OperationList     []Operation
	OperationID       *int
	OperationFileChan chan string
	OperationChan     chan Operation
)

func init() {
	OperationMap = make(map[string]Operation)
	OperationList = make([]Operation, 0)
	OperationFileChan = make(chan string)
	OperationChan = make(chan Operation)
	var WholeOperationFileString string
	go func() {
		FileListDir("./saved_infos/operation/", OperationFileChan)
		close(OperationFileChan)
	}()
	for FileName := range OperationFileChan {
		fileString, err := FileRead(FileName)
		if err != nil {
			beego.Error(err)
		} else {
			WholeOperationFileString += fileString
		}
		go func() {
			OperationJSONRead([]byte(WholeOperationFileString), OperationChan)
			close(OperationChan)
		}()
		for operationEntity := range OperationChan {
			OperationMap[operationEntity.OperationShortName] = operationEntity
			OperationList = append(OperationList, operationEntity)
		}
	}

}

// OperationJSONRead is function is for transfering JSON file to structure
func OperationJSONRead(fileByte []byte, OperationChan chan<- Operation) (chan<- Operation, error) {
	jsons, _ := simplejson.NewJson(fileByte)
	for _, jsonEntity := range jsons.MustArray() {
		operation := Operation{}
		err := mapstructure.WeakDecode(jsonEntity.(map[string]interface{}), &operation)
		if err != nil {
			beego.Error(err)
			return nil, err
		}
		OperationChan <- operation
	}
	return OperationChan, nil
}

// GetOperatoin function is for listing the operation entity by its ID
func GetOperatoin(oname string) (o Operation, err error) {
	if o, ok := OperationMap[oname]; ok {
		return o, nil
	}
	return o, errors.New("Operation not exist")
}

// GetAllOperations function is for listing all operation entity
func GetAllOperations() []Operation {
	return OperationList
}

// DeleteOperation function is for deleting the operation entity by its ID
func DeleteOperation(oname string) {
	delete(OperationMap, oname)
}

// AddOperation function is for adding the operation entity and return its ID
func AddOperation(o Operation) int {
	o.OperationID = *OperationID
	*OperationID++
	OperationMap[o.OperationShortName] = o
	return o.OperationID
}
