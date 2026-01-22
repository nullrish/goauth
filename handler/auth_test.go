package handler

import (
	"fmt"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nullrish/goauth/internal/generator"
)

func TestJwt(t *testing.T) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "Rishab"
	claims["exp"] = time.Now().Add(72 * time.Hour).Unix()

	data, err := token.SignedString([]byte("SYBAUwmiodj3oi2j902890392==d-2122321dd."))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(data)

	start := time.Now()
	generator.InitializeNode()
	for i := 0; i < 50; i++ {
		fmt.Println(generator.GenerateID())
	}
	fmt.Println("Took:", time.Since(start).Nanoseconds())
}
