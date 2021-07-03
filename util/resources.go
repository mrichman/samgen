package util

import (
	"fmt"
	"strings"

	"github.com/awslabs/goformation/v5/cloudformation/serverless"
	"github.com/getkin/kin-openapi/openapi3"
)

func GenerateServerlessFunction(verb string, operation *openapi3.Operation) (*serverless.Function, error) {

	events := map[string]serverless.Function_EventSource{
		strings.ToTitle(verb) + "Resource": {
			Type: "Api",
			Properties: &serverless.Function_Properties{
				ApiEvent: &serverless.Function_ApiEvent{
					Path:   "/" + operation.OperationID,
					Method: strings.ToLower(verb),
				},
			},
		},
	}

	functionPath := fmt.Sprintf("functions/%s", operation.OperationID)
	function := &serverless.Function{
		Description:  operation.Description,
		FunctionName: strings.Title(operation.OperationID),
		Handler:      "app.lambda_handler",
		CodeUri: &serverless.Function_CodeUri{
			String: &functionPath,
		},
		Runtime: "python3.9",
		Events:  events,
	}
	return function, nil
}
