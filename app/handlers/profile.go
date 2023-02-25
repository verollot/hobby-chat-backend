package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/frisk038/livechat/business/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type business interface {
	CreateUser(ctx context.Context, user models.User) error
	SetHobbies(ctx context.Context, user string, hobby string) error
	GetHobbies(ctx context.Context, userID string) ([]models.Hobby, error)
	DelHobbies(ctx context.Context, userID string, hobbyID uuid.UUID) error
}

type HandlerProfile struct {
	business business
}

func NewHandlerProfile(b business) HandlerProfile {
	return HandlerProfile{business: b}
}

func (hp *HandlerProfile) PostUsers(c *gin.Context) {
	ctx := c.Request.Context()
	var user models.User
	if err := c.BindJSON(&user); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	//TODO better fine tuning err
	if err := hp.business.CreateUser(ctx, user); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (hp *HandlerProfile) PostUsersHobbies(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("user_id")
	hobby := c.Param("hobby")
	if len(userID) == 0 || len(hobby) == 0 {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("id and hobby are mandatory"))
		return
	}

	//TODO better fine tuning err
	if err := hp.business.SetHobbies(ctx, userID, hobby); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}

func (hp *HandlerProfile) GetUsersHobbies(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("user_id")
	if len(userID) == 0 {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("id and hobby are mandatory"))
		return
	}

	//TODO better fine tuning err
	hobbies, err := hp.business.GetHobbies(ctx, userID)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"hobbies": hobbies})
}

func (hp *HandlerProfile) DelUsersHobbies(c *gin.Context) {
	ctx := c.Request.Context()
	userID := c.Param("user_id")
	hobby := c.Param("hobby_id")
	if len(userID) == 0 || len(hobby) == 0 {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("id and hobby are mandatory"))
		return
	}
	hobbyID, err := uuid.Parse(hobby)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, fmt.Errorf("id is not an uuid"))
		return
	}

	//TODO better fine tuning err
	if err := hp.business.DelHobbies(ctx, userID, hobbyID); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
