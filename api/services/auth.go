package services

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/tiago123456789/monitor_background_job/models"
	"github.com/tiago123456789/monitor_background_job/repositories"
	"golang.org/x/crypto/bcrypt"
)

type AuthInterface interface {
	Login(credential models.Credential) (models.ResponseAuthenticateion, error)
	IsAuthenticated(token string) error
}

type Auth struct {
	CompanyRepository repositories.CompanyRepositoryInterface
}

func NewAuth(companyRepository repositories.CompanyRepositoryInterface) *Auth {
	return &Auth{
		CompanyRepository: companyRepository,
	}
}

func (a *Auth) IsAuthenticated(token string) error {
	_, err := models.NewToken().IsValid(token)
	if err != nil {
		return err
	}
	return nil
}

func (a *Auth) Login(
	credential models.Credential) (models.ResponseAuthenticateion, error) {
	result, err := a.CompanyRepository.FindByEmail(credential.Email)
	if err != nil {
		return models.ResponseAuthenticateion{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(credential.Password))
	if err != nil {
		return models.ResponseAuthenticateion{}, err
	}

	tokenModel := models.NewToken()
	atClaims := jwt.MapClaims{}
	atClaims["companyId"] = result.ID
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	token, _ := tokenModel.Get(atClaims)

	return models.ResponseAuthenticateion{
		AccessToken: token,
	}, nil
}
