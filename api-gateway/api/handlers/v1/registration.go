package v1

import (
	"context"
	"encoding/json"
	"fmt"

	// "fmt"
	"log"
	"net/http"
	"net/smtp"

	// "strconv"
	"strings"
	"time"

	"api-gateway/api/handlers/models"
	token "api-gateway/api/handlers/tokens"
	"api-gateway/pkg/etc"
	l "api-gateway/pkg/logger"
	pbu "api-gateway/protos/template-service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"google.golang.org/protobuf/encoding/protojson"
)

// jwt

// Register ...
// @Summary Register
// @Description Api for registration
// @Tags register
// @Accept json
// @Produce json
// @Param User body models.User true "createUserModel"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/register/ [post]
func (h *handlerV1) Register(c *gin.Context) {
	var (
		body        models.RegisterUser
		jsonMarshal protojson.MarshalOptions
	)

	jsonMarshal.UseProtoNames = true
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to bind json", l.Error(err))
		return
	}

	body.Email = strings.TrimSpace(body.Email)
	body.Email = strings.ToLower(body.Email)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	if err = body.Validate(); err != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "This password is already in use or email error, please choose another",
		})
		h.log.Error("failed to check email uniques", l.Error(err))
		return
	}

	exists, err := h.serviceManager.UserService().CheckUniques(ctx, &pbu.CheckUniquesRequest{
		Field: "email",
		Value: body.Email,
	})

	if err != nil && exists.IsExist {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to check email uniques", l.Error(err))
		return
	}

	if exists.IsExist {
		c.JSON(http.StatusConflict, gin.H{
			"error": "This email is already in use, please choose another",
		})
		h.log.Error("failed to check email uniques", l.Error(err))
		return
	}

	exists, err = h.serviceManager.UserService().CheckUniques(ctx, &pbu.CheckUniquesRequest{
		Field: "user_name",
		Value: body.UserName,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		h.log.Error("failed to check username uniques", l.Error(err))
	}

	if exists.IsExist {
		c.JSON(http.StatusConflict, gin.H{
			"error": "This username is already in use, please choose another",
		})
		h.log.Error("failed to check username uniques", l.Error(err))
		return

	}

	byteData, err := json.Marshal(body)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Not marshaled code",
		})
	}

	client := redis.NewClient(&redis.Options{
		Addr:     "redisdb:6379",
	})

	code := etc.GenerateCode(6)

	err = client.Set(context.Background(), code, byteData, time.Hour*2).Err()
	if err != nil {
		log.Fatal(err)
	}

	auth := smtp.PlainAuth("", "asadfaxriddinov611@gmail.com", "drkeagdlwrfanrdp", "smtp.gmail.com")
	err = smtp.SendMail("smtp.gmail.com:587", auth, "asadfaxriddinov611@gmail.com", []string{body.Email}, []byte(code))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	c.JSON(http.StatusOK, true)
}

// Register ...
// @Summary Login
// @Description Api for Login
// @Tags register
// @Accept json
// @Produce json
// @Param email query string true "EMAIL"
// @Param password query string true "PASSWORD"
// @Success 200 {object} models.User
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/login/ [get]
func (h *handlerV1) LogIn(c *gin.Context) {

	// fmt.Println("kirdi")
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	email := c.Query("email")
	password := c.Query("password")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	respouser, err := h.serviceManager.UserService().GetUserByEmail(ctx, &pbu.EmailRequest{
		Email:    email,
		Password: password,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect email or password. Please try again",
		})
		h.log.Error(err.Error())
		return
	}

	if !etc.CheckPasswordHash(password, respouser.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect password. Please try again",
		})
	}

	// Create access and refresh tokens JWT
	h.jwthandler = token.JWTHandler{
		Sub:       respouser.Id,
		Iss:       time.Now().String(),
		Exp:       time.Now().Add(time.Hour * 6).String(),
		Role:      "user",
		SigninKey: h.cfg.SigningKey,
		Timeot:    h.cfg.AccessTokenTimout,
	}

	// aksestoken bn refreshtokeni generatsa qiliah
	access, _, err := h.jwthandler.GenerateAuthJWT()

	if err != nil {
		c.JSON(http.StatusInternalServerError, "error generating token")
		return
	}

	fmt.Println(access)

	c.JSON(http.StatusOK, models.UserByAccess{
		Id:          respouser.Id,
		Name:        respouser.Name,
		LastName:    respouser.LastName,
		Email:       respouser.Email,
		Password:    respouser.Password,
		UserName:    respouser.UserName,
		AccessToken: access,
	})
}

// Verification
// @Summary Verification User
// @Description LogIn - Api for verification users
// @Tags register
// @Accept json
// @Produce json
// @Param email query string true "Email"
// @Param code query string true "Code"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/verification [get]
func (h *handlerV1) Verification(c *gin.Context) {

	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	email := c.Query("email")
	code := c.Query("code")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	// time
	// time.Sleep(time.Duration(time.Minute))

	rdb := redis.NewClient(&redis.Options{
		Addr:     "redisdb:6379",
	})
	defer rdb.Close()

	// fmt.Println("keldi")

	val, err := rdb.Get(ctx, code).Result()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Verification code is expired",
		})
		h.log.Error("Verification code is expired", l.Error(err))
		return
	}

	// fmt.Println("ulgurdi")
	var userdetail models.User
	if err := json.Unmarshal([]byte(val), &userdetail); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unmarshiling error",
		})
		h.log.Error("error unmarshalling userdetail", l.Error(err))
		return
	}

	fmt.Println("email", email)
	fmt.Println("userdetail", userdetail.Email)

	if email != userdetail.Email {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Incorrect email. Try again",
		})
		return
	}

	id := uuid.New().String()

	hashPassword, err := etc.HashPassword(userdetail.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error message",
		})
		h.log.Error("Error hashing possword", l.Error(err))
		return
	}

	// Create access and refresh tokens JWT
	h.jwthandler = token.JWTHandler{
		Sub:       id,
		Iss:       time.Now().String(),
		Exp:       time.Now().Add(time.Hour * 6).String(),
		Role:      "admin",
		SigninKey: h.cfg.SigningKey,
		Timeot:    h.cfg.AccessTokenTimout,
	}
	// aksestoken bn refreshtokeni generatsa qiliah
	access, refresh, err := h.jwthandler.GenerateAuthJWT()

	if err != nil {
		c.JSON(http.StatusInternalServerError, "error generating token")
		return
	}

	createdUser, err := h.serviceManager.UserService().CreateUser(ctx, &pbu.User{
		Id:           id,
		Name:         userdetail.Name,
		LastName:     userdetail.LastName,
		Email:        userdetail.Email,
		Password:     hashPassword,
		UserName:     userdetail.UserName,
		RefreshToken: refresh,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error creating user",
		})
		h.log.Error("failed to create user", l.Error(err))
		return
	}

	response := &models.UserResponse{
		Id:           id,
		Name:         createdUser.Name,
		LastName:     createdUser.LastName,
		Username:     createdUser.UserName,
		Email:        createdUser.Email,
		Password:     hashPassword,
		AccessToken:  access,
		RefreshToken: refresh,
		// CreatedAt:    time.Now().String(),
		// UpdatedAt:    time.Now().String(),
	}

	c.JSON(http.StatusOK, response)
}

// frontdan reshresh token keladi
// token bn userni malumotini olaman
// userni id bn acc ref token generatsiya qilaman
// userni ref tokenini update qilaman
// userni hamma malumoti bilan accses tokeni qaytaraman

// @Summary Verification User
// @Description refresh token user
// @Tags register
// @Accept json
// @Produce json
// @Param refresh_token query string true "REFRESHTOKEN"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.StandardErrorModel
// @Failure 500 {object} models.StandardErrorModel
// @Router /v1/refreshusertoken [get]
func (h *handlerV1) RefreshUserToken(c *gin.Context) {

	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	refresh_token := c.Query("refresh_token")
	// code := c.Query("code")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.CtxTimeout))
	defer cancel()

	user, err := h.serviceManager.UserService().GetUserByRefreshToken(ctx, &pbu.RefreshToken{
		RefreshToken: refresh_token,
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error getRefreshToken user",
		})
		h.log.Error("failed to get user", l.Error(err))
		return
	}

	// Create access and refresh tokens JWT
	h.jwthandler = token.JWTHandler{
		Sub:       user.Id,
		Iss:       time.Now().String(),
		Exp:       time.Now().Add(time.Hour * 6).String(),
		Role:      "user",
		SigninKey: h.cfg.SigningKey,
		Timeot:    h.cfg.AccessTokenTimout,
	}

	// aksestoken bn refreshtokeni generatsa qiliah
	access, refresh, err := h.jwthandler.GenerateAuthJWT()
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error generating token")
		return
	}

	user.RefreshToken = refresh

	updateuser, err := h.serviceManager.UserService().UpdateUser(ctx, user)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Error updating user",
		})
		h.log.Error("failed to updating user", l.Error(err))
		return
	}

	c.JSON(http.StatusOK, models.UserByAccess{
		Id:          updateuser.Id,
		Name:        updateuser.Name,
		LastName:    updateuser.LastName,
		Email:       updateuser.Email,
		Password:    updateuser.Password,
		UserName:    updateuser.UserName,
		AccessToken: access,
	})
}
