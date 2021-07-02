package cmd

import (
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

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
		outfile, _ := cmd.Flags().GetString("output")
		verbose, _ := cmd.Flags().GetBool("verbose")

		if verbose {
			log.SetLevel(log.DebugLevel)
			log.Debug("Verbose logs enabled")
		}

		export, _ := apigw.ExportRestApi(apiId, stage)

		loader := openapi3.NewLoader()
		doc, err := loader.LoadFromData(export)
		if err != nil {
			log.Fatalf("Could not load Open API spec: %v", err)
		}
		err = doc.Validate(loader.Context)
		if err != nil {
			log.Fatalf("Could not validate Open API spec: %v", err)
		}

		// Create a new CloudFormation template
		template := cloudformation.NewTemplate()
		transform := "AWS::Serverless-2016-10-31"
		template.Transform = &cloudformation.Transform{String: &transform}

		for _, pathItem := range doc.Paths {
			for verb, operation := range pathItem.Operations() {
				log.Debugf("Generating resource for: %s /%s", verb, operation.OperationID)
				function, _ := util.GenerateServerlessFunction(verb, operation)
				template.Resources[strings.Title(operation.OperationID)+"Function"] = function
			}
		}

		// Output the YAML AWS CloudFormation template
		y, err := template.YAML()
		if err != nil {
			log.Fatalf("Failed to generate YAML: %s", err)
		} else {
			log.Tracef("%s", string(y))
			f, err := os.Create(outfile)
			if err != nil {
				log.Fatalf("Could not create file %s: %v", outfile, err)
			}
			_, err = f.Write(y)
			if err != nil {
				log.Fatalf("Could not write file %s: %v", outfile, err)
				f.Close()
				return
			}
		}
	},
}

func init() {
	importCmd.Flags().String("rest-api-id", "", "The string identifier of the associated RestApi (e.g. a1b2c3d4e5)")
	importCmd.Flags().String("stage", "", "The name of the Stage that will be exported (e.g. prod)")
	importCmd.Flags().String("output", "", "The filename to save the SAM template as (e.g. template.yaml)")
	importCmd.PersistentFlags().Bool("verbose", false, "Verbose logging")
	importCmd.MarkFlagRequired("rest-api-id")
	importCmd.MarkFlagRequired("stage")
	importCmd.MarkFlagRequired("output")
	rootCmd.AddCommand(importCmd)
}
