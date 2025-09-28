package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v4"
)


func ValidateJWTMiddleware(next func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) ) func (request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func (request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		token := extractTokenFromHeader(request.Headers)
		if token == "" {
			return events.APIGatewayProxyResponse{
				Body:  "Unauthorized Request",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}
		claims, err := pasreToken(token)
		if err != nil {
			return events.APIGatewayProxyResponse{
				Body:  "Unauthorized Request",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		expires := claims["exp"].(float64)
		if expires < float64(time.Now().Unix()) {
			return events.APIGatewayProxyResponse{
				Body:  "Unauthorized Request, Token Expired",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		request.Headers["username"] = claims["user"].(string)
		return next(request)
	}
}


func Protectedhandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "Protected Resource",
		StatusCode: http.StatusOK,
	}, nil
}


func extractTokenFromHeader(headers map[string]string) string {
	authToken , ok := headers["Authorization"]
	if !ok {
		return ""
	}
	splitToken := strings.Split(authToken, "Bearer ")

	if len(splitToken) != 2 {
		return ""
	}
	return splitToken[1]
}

func pasreToken(tokenString string) (jwt.MapClaims, error)  {
	secret := "myPersonalAmitSecret"
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	claims , ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}