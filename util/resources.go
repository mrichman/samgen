package util

import (
	"strings"

	"github.com/awslabs/goformation/v5/cloudformation/serverless"
	"github.com/getkin/kin-openapi/openapi3"
)

func GenerateServerlessFunction(verb string, operation *openapi3.Operation) (*serverless.Function, error) {
	function := &serverless.Function{
		Description: operation.Description, FunctionName: strings.Title(operation.OperationID),
	}
	return function, nil
}
