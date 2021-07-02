package apigateway

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"

	"github.com/aws/aws-sdk-go-v2/config"
	apigw "github.com/aws/aws-sdk-go-v2/service/apigateway"
)

func ExportRestApi(apiId string, stageName string) ([]byte, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	client := apigw.NewFromConfig(cfg)
	output, err := client.GetExport(context.TODO(),&apigw.GetExportInput{
		ExportType: aws.String("oas30"), // oas30 or swagger
		RestApiId: aws.String(apiId), // a1b2c3d4e5
		StageName: aws.String(stageName), // dev
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(output.Body)

	return output.Body, nil
}