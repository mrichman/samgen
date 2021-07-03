package util

import (
	"embed"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"

	"github.com/awslabs/goformation/v5/cloudformation/serverless"
	"github.com/getkin/kin-openapi/openapi3"
)

// GenerateServerlessFunction creates an AWS::Serverless::Function resource for the given verb and operation
func GenerateServerlessFunction(verb string, operation *openapi3.Operation) (*serverless.Function, error) {

	events := map[string]serverless.Function_EventSource{
		strings.Title(strings.ToLower(verb)) + "Resource": {
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
		Runtime: "python3.8",
		Events:  events,
	}
	return function, nil
}

//go:embed python/*
var f embed.FS

// WriteFunction writes dummy Lambda function content to the specified directory
func WriteFunction(functionPath string) error {
	appPyData, err := f.ReadFile("python/app.py")

	if err != nil {
		log.Error(err)
	}

	log.Infof("Creating directory: %s", functionPath)
	if err := os.MkdirAll(functionPath, os.ModePerm); err != nil {
		log.Fatalf("Could not create directory %s", functionPath)
	}

	err = os.WriteFile(filepath.Join(functionPath, "app.py"), appPyData, 0644)
	if err != nil {
		log.Fatal(err)
	}

	Touch(filepath.Join(functionPath, "requirements.txt"))

	return nil
}
