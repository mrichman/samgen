package cmd

import (
	"fmt"
	"strings"

	"github.com/awslabs/goformation/v5/cloudformation"
	"github.com/getkin/kin-openapi/openapi3"
	apigw "github.com/mrichman/samgen/apigateway"
	"github.com/mrichman/samgen/util"
	"github.com/spf13/cobra"
)

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import an existing API Gateway REST API",
	Long:  "Import an existing API Gateway REST API",
	Run: func(cmd *cobra.Command, args []string) {

		apiId, _ := cmd.Flags().GetString("rest-api-id")
		stage, _ := cmd.Flags().GetString("stage")
		export, _ := apigw.ExportRestApi(apiId, stage)

		loader := openapi3.NewLoader()
		doc, err := loader.LoadFromData(export)
		if err != nil {
			panic(err)
		}
		err = doc.Validate(loader.Context)
		if err != nil {
			panic(err)
		}

		// Create a new CloudFormation template
		template := cloudformation.NewTemplate()
		transform := "AWS::Serverless-2016-10-31"
		template.Transform = &cloudformation.Transform{String: &transform}

		for _, pathItem := range doc.Paths {
			// fmt.Println("Key:", key)
			for verb, operation := range pathItem.Operations() {
				// fmt.Println("Verb:", verb,"=>", "OperationID:", operation.OperationID)

				function, _ := util.GenerateServerlessFunction(verb, operation)
				template.Resources[strings.Title(operation.OperationID)+"Function"] = function

			}
		}

		// Output the YAML AWS CloudFormation template
		y, err := template.YAML()
		if err != nil {
			fmt.Printf("Failed to generate YAML: %s\n", err)
		} else {
			fmt.Printf("%s\n", string(y))
		}
	},
}

func init() {
	importCmd.Flags().String("rest-api-id", "", "The string identifier of the associated RestApi")
	importCmd.Flags().String("stage", "", "The name of the Stage that will be exported")
	importCmd.MarkFlagRequired("rest-api-id")
	importCmd.MarkFlagRequired("stage")
	rootCmd.AddCommand(importCmd)
}
