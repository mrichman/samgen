package cmd

import (
	"strings"

	"github.com/awslabs/goformation/v5/cloudformation"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/mrichman/samgen/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var specCmd = &cobra.Command{
	Use:   "spec",
	Short: "Import an OpenAPI 3.0 spec file",
	Long:  "Import an OpenAPI 3.0 spec file",
	Run: func(cmd *cobra.Command, args []string) {

		spec, _ := cmd.Flags().GetString("spec")

		loader := openapi3.NewLoader()
		doc, err := loader.LoadFromFile(spec)
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
			log.Errorf("Failed to generate YAML: %s", err)
		} else {
			log.Errorf("%s\n", string(y))
		}
	},
}

func init() {
	specCmd.Flags().String("spec", "", "OpenAPI 3.0 spec file")
	specCmd.MarkFlagRequired("spec")
	rootCmd.AddCommand(specCmd)
}
