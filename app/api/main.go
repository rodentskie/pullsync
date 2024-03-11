package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"slack-pr-lambda/api/routes"
	"slack-pr-lambda/constants"
	"slack-pr-lambda/env"
	"slack-pr-lambda/logger"
	"syscall"

	"go.uber.org/zap"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/httpadapter"
)

var httpLambda *httpadapter.HandlerAdapter

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return httpLambda.ProxyWithContext(ctx, req)
}

func main() {
	l := logger.LoggerConfig()
	zapLog, _ := l.Build()
	ports := constants.Port()

	defer func() {
		if err := zapLog.Sync(); err != nil && !errors.Is(err, syscall.EINVAL) {
			log.Fatalf("error closing the logger. %v\n", err)
		}
	}()

	portString := fmt.Sprintf(":%d", ports.MainApi)

	port := env.GetEnv("PORT", portString)
	host := env.GetEnv("HOST", "localhost")
	env := env.GetEnv("ENV", "local")

	mux := http.NewServeMux()
	routes.MainRoutes(mux)

	if env == "local" {
		sandboxLink := fmt.Sprintf("http://%s%s", host, port)
		zapLog.Info("running at 🚀⚙️",
			zap.String("link", sandboxLink),
		)

		if err := http.ListenAndServe(port, mux); err != nil && err != http.ErrServerClosed {
			zapLog.Fatal("error serve api",
				zap.String("port", port),
				zap.Error(err),
			)
		}
	}

	httpLambda = httpadapter.New(mux)
	lambda.Start(Handler)
}
