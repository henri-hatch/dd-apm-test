package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	ddlambda "github.com/DataDog/datadog-lambda-go"

	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"

	"github.com/gin-gonic/gin"

	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
)

var ginLambda *ginadapter.GinLambdaV2

func init() {
	r := gin.Default()
	r.Use(gintrace.Middleware("apm-test"))
	r.GET("/ping", PingPong)

	ginLambda = ginadapter.NewV2(r)
}

func PingPong(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func HandleRequest(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// Submit a custom metric
	ddlambda.Metric(
		"test.ppm_metric",                  // Metric name
		12.45,                              // Metric value
		"product:ppm_bot", "hotel:trivago", // Associated tags
	)

	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	// Start the tracing
	tracer.Start()
	defer tracer.Stop()
	// Wrap the handler function with Datadog Lambda Wrapper
	lambda.Start(ddlambda.WrapFunction(HandleRequest, nil))
}
