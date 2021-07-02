package apigateway

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/config"

	apigw "github.com/aws/aws-sdk-go-v2/service/apigateway"
)

func ExportRestApi(apiId string, stage string) ([]byte, error) {

	log.Infof("Exporting API %s stage %s\n", apiId, stage)

	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		log.Fatalf("Unable to load AWS SDK config: %v\n", err)
	}
	client := apigw.NewFromConfig(cfg)
	output, err := client.GetExport(context.TODO(), &apigw.GetExportInput{
		ExportType: aws.String("oas30"), // oas30 or swagger
		RestApiId:  aws.String(apiId),   // a1b2c3d4e5
		StageName:  aws.String(stage),   // dev
		Accepts:    aws.String("application/yaml"),
		Parameters: map[string]string{
			"extensions": "integrations",
		},
	})

	if err != nil {
		log.Fatalf("Fatal error getting export from API Gateway: %v\n", err)
	}

	log.Trace(string(output.Body))

	return output.Body, nil
}
