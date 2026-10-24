package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"notix/internal/models"
	"notix/internal/repository"
)

var ErrInvalidCredentials = errors.New("invalid credentials")

type AccountStore interface {
	FindByEmail(email string) (models.Cuenta, error)
}

type AuthService struct {
	accounts AccountStore
	secret   []byte
	expiry   time.Duration
}

type TokenClaims struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"rol"`
	jwt.RegisteredClaims
}

type LoginResult struct {
	Token     string
	ExpiresAt time.Time
	User      models.Usuario
	Role      string
}

func NewAuthService(accounts AccountStore, secret []byte, expiry time.Duration) *AuthService {
	return &AuthService{accounts: accounts, secret: secret, expiry: expiry}
}

func (s *AuthService) Login(email, password string) (LoginResult, error) {
	account, err := s.accounts.FindByEmail(email)
	if errors.Is(err, repository.ErrNotFound) {
		return LoginResult{}, ErrInvalidCredentials
	}
	if err != nil {
		return LoginResult{}, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(account.ClaveHash), []byte(password)); err != nil {
		return LoginResult{}, ErrInvalidCredentials
	}

	expiresAt := time.Now().Add(s.expiry)
	claims := TokenClaims{
		UserID: account.UsuarioID,
		Role:   account.Rol,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "notix-api",
			Subject:   fmt.Sprint(account.UsuarioID),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
	if err != nil {
		return LoginResult{}, fmt.Errorf("sign access token: %w", err)
	}
	return LoginResult{Token: token, ExpiresAt: expiresAt, User: account.Usuario, Role: account.Rol}, nil
}

func (s *AuthService) Authenticate(token string) (TokenClaims, error) {
	claims := &TokenClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing algorithm %s", token.Method.Alg())
		}
		return s.secret, nil
	})
	if err != nil || !parsed.Valid || claims.UserID == 0 {
		return TokenClaims{}, ErrInvalidCredentials
	}
	return *claims, nil
}
