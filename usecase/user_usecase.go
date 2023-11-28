package usecase

import (
	"errors"
	"fmt"
	"net/http"
	"store/model"
	"store/repository"
	"store/utils"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/mitchellh/mapstructure"
	"gorm.io/gorm"
)

type UserUsecase interface {
	CreateUser(user *model.User) error
	UserLogin(response *gin.H, user *model.User) error
	ReadAllUser(selected *[]model.User, page int, limit int, order string) error
	UpdateUser(updatedUser *model.User) error
	DeleteUser(deletedUser *model.User) error
	UserVerify() gin.HandlerFunc
}

type userUsecase struct {
	dbRepo    repository.DbRepository
	tokenRepo repository.TokenRepository
}

type AuthToken struct {
	UserID   uint
	UserName string
}

type MyClaims struct {
	jwt.StandardClaims
	AuthToken AuthToken `json:"AuthToken,omitempty"`
}

func (u *userUsecase) CreateUser(user *model.User) error {
	HashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	user.Password = HashedPassword
	return u.dbRepo.Create(user)
}

func (u *userUsecase) UserLogin(response *gin.H, user *model.User) error {
	var selected model.User
	err := u.dbRepo.Read(utils.ReadOption{Many: true, Result: &selected, Where: map[string]interface{}{"name": user.Name}})
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return utils.DataNotFoundError()
		} else {
			return err
		}
	}

	isPasswordMatch, err := utils.MatchPassword(selected.Password, user.Password)
	if err != nil {
		return err
	}

	if !isPasswordMatch {
		return utils.WrongPasswordError()
	}

	duration := 12 * time.Hour
	now := time.Now().UTC()
	end := now.Add(duration)

	newToken, err := u.tokenRepo.CreateToken(func(appName string) jwt.Claims {
		return MyClaims{
			StandardClaims: jwt.StandardClaims{
				Issuer:    appName,
				IssuedAt:  now.Unix(),
				ExpiresAt: end.Unix(),
			},
			AuthToken: AuthToken{
				UserID:   selected.ID,
				UserName: selected.Name,
			},
		}
	})
	if err != nil {
		return err
	}
	*response = gin.H{
		"token": newToken,
		"user": gin.H{
			"ID":           selected.ID,
			"userName":     selected.Name,
			"userFullname": selected.FullName,
		},
	}
	return nil

}

func (u *userUsecase) ReadAllUser(selected *[]model.User, page int, limit int, order string) error {
	offset := limit * (page - 1)
	err := u.dbRepo.Read(utils.ReadOption{Limit: &limit, Offset: &offset, Many: true, Result: &selected, Order: &order})
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) UpdateUser(updatedUser *model.User) error {
	err := u.dbRepo.Update(&updatedUser)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) DeleteUser(deletedUser *model.User) error {
	err := u.dbRepo.Delete(&deletedUser)
	if err != nil {
		return err
	}
	return nil
}

func (u *userUsecase) UserVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		var h struct {
			AuthorizationHeader string `header:"Authorization"`
		}
		if err := c.ShouldBindHeader(&h); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthrorized",
			})
			c.Abort()
			return
		}
		tokenString := strings.Replace(h.AuthorizationHeader, "Bearer ", "", -1)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthrorized",
			})
			c.Abort()
			return
		}

		mapClaims, err := u.tokenRepo.VerifyToken(tokenString)
		if err != nil {
			if errors.Is(err, fmt.Errorf("token expired")) {
				c.JSON(http.StatusUnauthorized, gin.H{
					"message": "Token Expired",
				})
				c.Abort()
				return
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthrorized",
			})
			c.Abort()
			return
		}
		fmt.Println(mapClaims["AuthToken"], mapClaims)

		var authToken AuthToken
		mapstructure.Decode(mapClaims["AuthToken"].(map[string]interface{}), &authToken)

		var selected model.User
		err = u.dbRepo.Read(utils.ReadOption{Many: false, Result: &selected, Where: map[string]interface{}{"id": authToken.UserID, "name": authToken.UserName}})
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthrorized",
			})
			c.Abort()
			return
		}

		c.Set("authToken", authToken)
		c.Next()
	}
}

func NewUserUsecase(dbRepo repository.DbRepository, tokenRepo repository.TokenRepository) UserUsecase {
	usecase := new(userUsecase)
	usecase.dbRepo = dbRepo
	usecase.tokenRepo = tokenRepo
	return usecase
}
