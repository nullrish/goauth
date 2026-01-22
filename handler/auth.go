// Package handler is used to handle the request on the register routes
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/mail"
	"regexp"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nullrish/goauth/database"
	"github.com/nullrish/goauth/internal/auth"
	"github.com/nullrish/goauth/internal/generator"
	"github.com/nullrish/goauth/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func existsByEmail(e string) (bool, error) {
	var count int64
	db := database.DB
	err := db.Model(&model.User{}).Where("email = ?", e).Count(&count).Error
	return count > 0, err
}

func existsByUsername(u string) (bool, error) {
	var count int64
	db := database.DB
	err := db.Model(&model.User{}).Where("username = ?", u).Count(&count).Error
	return count > 0, err
}

func getUserByEmail(e string) (*model.User, error) {
	db := database.DB
	var user model.User
	if err := db.Where(&model.User{Email: e}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func getUserByUsername(u string) (*model.User, error) {
	db := database.DB
	var user model.User
	if err := db.Where(&model.User{Username: u}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func isEmail(e string) bool {
	_, err := mail.ParseAddress(e)
	return err == nil
}

func isUser(u string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._-]{3,30}$`)
	return re.Match([]byte(u))
}

func Register(c fiber.Ctx) error {
	type RegisterInput struct {
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		Password    string `json:"password"`
	}

	type NewUser struct {
		Username    string `json:"username"`
		DisplayName string `json:"display_name"`
		Email       string `json:"email"`
	}

	input := new(RegisterInput)
	if err := json.Unmarshal(c.Body(), &input); err != nil {
		fmt.Println(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on register request", "data": err})
	}

	if !isEmail(input.Email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Enter a valid email", "data": nil})
	}

	if !isUser(input.Username) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Username must be 3-30 characters & can only contain -, _, ., alphabets, numbers"})
	}

	exists, err := existsByEmail(input.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Couldn't check for existing emails", "data": err})
	}
	if exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Email already taken", "data": err})
	}
	exists, err = existsByUsername(input.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Couldn't check for existing usernames", "data": err})
	}
	if exists {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Username already taken", "data": err})
	}

	hash, err := hashPassword(input.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Couldn't hash password", "data": err})
	}
	db := database.DB
	user := &model.User{
		ID:            generator.GenerateID(),
		Username:      input.Username,
		DisplayName:   input.DisplayName,
		Email:         input.Email,
		PhoneNumber:   input.PhoneNumber,
		Password:      hash,
		EmailVerified: false,
		PhoneVerified: false,
	}
	if err := db.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Couldn't create user", "data": err})
	}

	newUser := NewUser{
		Username:    user.Username,
		DisplayName: user.DisplayName,
		Email:       user.Email,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": "success", "message": "Created User", "data": newUser})
}

func Login(c fiber.Ctx) error {
	type LoginInput struct {
		Identity string `json:"identity"`
		Password string `json:"password"`
	}
	type UserData struct {
		ID       int64  `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	input := new(LoginInput)
	var userData UserData

	if err := json.Unmarshal(c.Body(), &input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on login request", "data": err})
	}

	identity := input.Identity
	pass := input.Password
	userModel, err := new(model.User), *new(error)

	if isEmail(identity) {
		userModel, err = getUserByEmail(identity)
	} else {
		userModel, err = getUserByUsername(identity)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Internal Server Error", "data": err})
	} else if userModel == nil {
		CheckPasswordHash(pass, "")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password", "data": err})
	} else {
		userData = UserData{
			ID:       userModel.ID,
			Username: userModel.Username,
			Email:    userModel.Email,
			Password: userModel.Password,
		}
	}

	if !CheckPasswordHash(pass, userData.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Invalid identity or password", "data": err})
	}

	/* 	token := jwt.New(jwt.SigningMethodHS256)

	   	claims := token.Claims.(jwt.MapClaims)
	   	claims["username"] = userData.Username
	   	claims["user_id"] = userData.ID
	   	claims["user_email"] = userData.Email
	   	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	   	t, err := token.SignedString([]byte(os.Getenv("SECRET")))
	   	if err != nil {
	   		return c.SendStatus(fiber.StatusInternalServerError)
	   	}
	*/
	t, err := auth.SignJWT(userModel)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Successful login", "data": t})
}

func VerifyAuth(c fiber.Ctx) error {
	type Input struct {
		Token string `json:"token"`
	}
	input := new(Input)
	if err := json.Unmarshal(c.Body(), &input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid json request.", "data": err})
	}

	t, err := auth.VerifyJWT(input.Token)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Can't verify session token", "data": err})
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Can't verify claims", "data": err})
	}

	type JWTClaims struct {
		ID         float64 `json:"id"`
		Username   string  `json:"username"`
		IssuedTime float64 `json:"iat"`
		ExpiryTime float64 `json:"exp"`
	}

	jwtClaims := &JWTClaims{
		ID:         claims["id"].(float64),
		Username:   claims["username"].(string),
		IssuedTime: claims["iat"].(float64),
		ExpiryTime: claims["exp"].(float64),
	}

	return c.JSON(fiber.Map{"status": "success", "message": "Verified The Token", "data": jwtClaims})
}
