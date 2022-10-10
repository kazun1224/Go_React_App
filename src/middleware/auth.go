package middleware

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

const SecretKey = "secret"

type ClaimsWithScope struct {
	jwt.StandardClaims
	Scope string
}

func IsAuthenticate(ctx *fiber.Ctx) error {
	// cookieから情報を取得
	cookie := ctx.Cookies("jwt")

	// token取得
	token, err := jwt.ParseWithClaims(
		cookie,
		&ClaimsWithScope{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		},
	)

	if err != nil || !token.Valid {
		ctx.Status(fiber.StatusUnauthorized) // 401
		return ctx.JSON(fiber.Map{
			"message": "認証がされていません",
		})
	}
	payload := token.Claims.(*ClaimsWithScope)
	IsAmbassador := strings.Contains(ctx.Path(), "/api/ambassador")

	if (payload.Scope == "admin" && IsAmbassador) || (payload.Scope == "ambassador" && !IsAmbassador) {
		ctx.Status(fiber.StatusUnauthorized) // 401
		return ctx.JSON(fiber.Map{
			"message": "認証が許可されません",
		})
	}

	return ctx.Next()
}

func GenerateJWT(id uint, scope string) (string, error) {
	// トークンの発行
	payload := ClaimsWithScope{}
	payload.Subject = strconv.Itoa(int(id))
	payload.ExpiresAt = time.Now().Add(time.Hour * 24).Unix()
	payload.Scope = scope

	return jwt.NewWithClaims(jwt.SigningMethodHS256, payload).SignedString([]byte(SecretKey))
}

func GetUserID(ctx *fiber.Ctx) (uint, error) {
	// cookieから情報を取得
	cookie := ctx.Cookies("jwt")

	// token取得
	token, err := jwt.ParseWithClaims(
		cookie,
		&ClaimsWithScope{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		},
	)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	payload := token.Claims.(*ClaimsWithScope)
	id, _ := strconv.Atoi(payload.Subject)
	return uint(id), nil
}
