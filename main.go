package main

import (
	"context"
	"os"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/fedect1/go-api-posts/awsgo"
	"github.com/fedect1/go-api-posts/secretmanager"
)

func main() {
	lambda.Start(ejecutoLambda)
}

func ejecutoLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res *events.APIGatewayProxyResponse

	awsgo.InicializoAWS()

	if !ValidoParametros() {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body: "Error en las variables de entorno. Deben incluir las variables de entorno",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	}

	SecretModel, err := secretmanager.GetSecret(os.Getenv("SecretName"))
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body: "Error en la lectura de Secret "+err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil	
	}
	
}

func ValidoParametros() bool {
	_, traerParametro := os.LookupEnv("SecretName")
	if !traerParametro {
		return traerParametro
	}
	_, traerParametro = os.LookupEnv("BucketName")
	if !traerParametro {
		return traerParametro
	}
	_, traerParametro = os.LookupEnv("UrlPrefix")
	if !traerParametro {
		return traerParametro
	}
	return traerParametro
}
