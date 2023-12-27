package main

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
	"github.com/fedect1/go-api-posts/awsgo"
	"github.com/fedect1/go-api-posts/bd"
	"github.com/fedect1/go-api-posts/handlers"
	"github.com/fedect1/go-api-posts/models"
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

	path := strings.Replace(request.PathParameters["postApp"], os.Getenv("UrlPrefix"), "", -1)
	
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("path"), path)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("method"), request.HTTPMethod)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("user"), SecretModel.Username)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("password"), SecretModel.Password)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("host"), SecretModel.Host)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("database"), SecretModel.Database)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("jwtSign"), SecretModel.JWTSign)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("body"), request.Body)
	awsgo.Ctx = context.WithValue(awsgo.Ctx, models.Key("bucketName"), os.Getenv("BucketName"))

	err = bd.ConectarDB(awsgo.Ctx)
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body: "Error conectando la DB "+err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil	
	}
	restAPI := handlers.Manejadores(awsgo.Ctx, request)
	if restAPI.CustomResp == nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: restAPI.Status,
			Body: restAPI.Message,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil	
	} else {
		return restAPI.CustomResp, nil
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
