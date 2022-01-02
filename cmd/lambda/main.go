package main

import (
	"context"

	"github.com/nns7/userapp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

func init() {
	router := gin.Default()
	// ユーザー作成
	router.POST("/users", userapp.PostUsers)
	// ユーザー一覧の取得
	router.GET("/users", userapp.GetUsers)
	// ユーザーの更新
	router.PUT("/users/:user_id", userapp.PutUser)
	// ユーザーの削除
	router.DELETE("/users/:user_id", userapp.DeleteUser)
	// ユーザーの検索
	router.GET("/users/search", userapp.SearchUser)
	ginLambda = ginadapter.New(router)
}

func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func main() {
	lambda.Start(Handler)
}
