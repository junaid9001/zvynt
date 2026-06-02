package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/junaid9001/zvynt/proto/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authClient auth.AuthServiceClient
}

func RegisterAuthRoutes(publicGroup, privateGroup *gin.RouterGroup, srv auth.AuthServiceClient) {

	h := &AuthHandler{authClient: srv}

	publicGroup.GET("auth/health", h.Health)
	publicGroup.POST("auth/signup", h.Signup)
	publicGroup.POST("auth/login", h.Login)

	privateGroup.GET("users/me", h.GetME)

}

func (h *AuthHandler) Health(c *gin.Context) {
	_, err := h.authClient.Health(c.Request.Context(), &auth.Empty{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal sever error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})

}

func (h *AuthHandler) Signup(c *gin.Context) {

	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"email,required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	user, err := h.authClient.Signup(c.Request.Context(), &auth.SignupRequest{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		log.Println("gRPC Signup error: ", err)
		st := status.Convert(err).Code()

		switch st {
		case codes.AlreadyExists:
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		case codes.Unavailable:
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "auth service is currently unreachable"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "an unexpected error occurred"})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user_id": user.UserId})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"email,required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	resp, err := h.authClient.Login(c.Request.Context(), &auth.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		st := status.Convert(err).Code()

		switch st {

		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "email not found",
			})

		case codes.Unauthenticated:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid password",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}

		return
	}

	c.SetCookie(
		"access_token",
		resp.AccessToken,
		60000,
		"/",
		"",
		false, //true in prod
		true,
	)

	c.JSON(http.StatusOK, gin.H{"user_id": resp.UserId, "access_token": resp.AccessToken})
}

func (h *AuthHandler) GetME(c *gin.Context) {

	userID := c.GetString("user_id")

	resp, err := h.authClient.GetUser(c.Request.Context(),
		&auth.UserRequest{
			Identifier: &auth.UserRequest_UserId{UserId: userID},
		})

	if err != nil {

		st := status.Convert(err).Code()

		switch st {

		case codes.NotFound:
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})

		case codes.InvalidArgument:
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "internal server error",
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{"user_id": resp.UserId, "email": resp.Email})
}
