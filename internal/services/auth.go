package service

import (
	"context"
	"erp/config"
	"erp/internal/api_errors"
	"erp/internal/constants"
	"erp/internal/domain"
	"erp/internal/infrastructure"
	"erp/internal/models"
	"erp/internal/repository"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	dbTransaction infrastructure.DatabaseTransaction
	userService   domain.UserService
	ufw           *repository.UnitOfWork
	jwtService    domain.JwtService
	config        *config.Config
}

func NewAuthService(
	dbTransaction infrastructure.DatabaseTransaction,
	userService domain.UserService,
	ufw *repository.UnitOfWork,
	jwtService domain.JwtService,
	config *config.Config,
) domain.AuthService {
	return &authService{
		dbTransaction: dbTransaction,
		userService:   userService,
		ufw:           ufw,
		jwtService:    jwtService,
		config:        config,
	}
}

func (a *authService) Register(ctx context.Context, req domain.RegisterRequest) (user *models.User, err error) {
	role := constants.RoleCustomer

	if req.RequestFrom != string(constants.Web) {
		role = constants.RoleSeller
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt password")
	}

	req.Password = string(encryptedPassword)
	user, err = a.userService.Create(ctx, &models.User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      role,
	})

	return user, err
}

func (a *authService) Login(ctx context.Context, req domain.LoginRequest) (res domain.LoginResponse, err error) {
	user, err := a.userService.GetByEmail(ctx, req.Email)
	if err != nil {
		return res, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return res, errors.New(api_errors.ErrInvalidPassword)
	}

	accessToken, refreshToken, err := a.jwtService.GenerateAuthTokens(user.ID.String())
	if err != nil {
		return res, err
	}

	err = a.dbTransaction.WithTransaction(func(tx *infrastructure.Database) error {
		if _, err := a.ufw.TokenRepository.Upsert(tx, ctx, &models.Token{
			UserID:    user.ID.String(),
			ExpiredAt: a.config.Jwt.AccessTokenExpiresIn,
		}); err != nil {
			return err
		}

		if err := a.ufw.UserRepository.UpdateLastLogin(tx, ctx, user.ID.String()); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return res, err
	}

	res = domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    a.config.Jwt.AccessTokenExpiresIn,
		UserInfo: domain.UserInfoLoginResponse{
			ID:        user.ID.String(),
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
		},
	}

	return res, nil
}
