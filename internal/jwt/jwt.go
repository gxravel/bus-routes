package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/gxravel/bus-routes/internal/config"
	ierr "github.com/gxravel/bus-routes/internal/errors"
	log "github.com/gxravel/bus-routes/internal/logger"
	"github.com/gxravel/bus-routes/internal/model"
	"github.com/gxravel/bus-routes/internal/storage"

	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
)

// Manager includes the methods allowed to deal with the token.
type Manager interface {
	save(ctx context.Context, token *Details) error

	Parse(tokenString string) (*Claims, error)
	CheckIfExists(ctx context.Context, tokenUUID string) error
	Delete(ctx context.Context, tokenUUID string) error
	SetNew(ctx context.Context, user *User) (*Details, error)
	Verify(ctx context.Context, tokenString string) (*User, error)
}

// User describes user built into the token
type User struct {
	ID    int64          `json:"id"`
	Email string         `json:"email"`
	Type  model.UserType `json:"type"`
}

// Claims defines JWT token claims.
type Claims struct {
	User *User `json:"user"`
	jwt.StandardClaims
}

// Details defines the structure of a JWT token.
type Details struct {
	String  string
	Expiry  int64
	UUID    string
	Subject int64
}

// JWT contains the fields which interact with the token.
type JWT struct {
	client *storage.Client
	config config.Config
}

func New(client *storage.Client, config config.Config) *JWT {
	return &JWT{client: client, config: config}
}

// create creates the HS512 JWT token with claims.
func create(ctx context.Context, user *User, expiry time.Duration, key string) (*Details, error) {
	now := time.Now()
	token := &Details{}
	token.Expiry = now.Add(expiry).Unix()
	claims := &Claims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			Id:        uuid.NewV4().String(),
			IssuedAt:  now.Unix(),
			ExpiresAt: token.Expiry,
		},
	}
	token.UUID = claims.Id
	token.Subject = claims.User.ID

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	var err error
	token.String, err = jwtToken.SignedString([]byte(key))
	if err != nil {
		log.FromContext(ctx).
			WithErr(err).
			Fatal("failed to sign jwt token with the key")
	}

	return token, err
}

// Parse parses a string token with the key.
func (m *JWT) Parse(tokenString string) (*Claims, error) {
	var key = []byte(m.config.JWT.AccessKey)

	jwtToken, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ierr.
				NewReason(ierr.ErrInvalidJWT).
				WithMessage(fmt.Sprintf("unexpected signing method: %v", t.Header["alg"]))
		}

		return key, nil
	})
	if err != nil || !jwtToken.Valid {
		return nil, ierr.NewReason(ierr.ErrInvalidToken).WithMessage("token validation failed")
	}

	claims, ok := jwtToken.Claims.(*Claims)
	if !ok {
		return nil, ierr.NewReason(ierr.ErrInvalidToken).WithMessage("failed to get claims")
	}

	return claims, nil
}

// save saves the token to the storage database.
func (m *JWT) save(ctx context.Context, token *Details) error {
	expiry := time.Until(time.Unix(token.Expiry, 0))
	return m.client.Set(ctx, token.UUID, token.Subject, expiry).Err()
}

// CheckIfExists checks if token exists in the storage database.
func (m *JWT) CheckIfExists(ctx context.Context, tokenUUID string) error {
	return m.client.Get(ctx, tokenUUID).Err()
}

// Delete deletes token from the storage database.
func (m *JWT) Delete(ctx context.Context, tokenUUID string) error {
	return m.client.Del(ctx, tokenUUID).Err()
}

// SetNew returns the access token.
func (m *JWT) SetNew(ctx context.Context, user *User) (*Details, error) {
	logger := log.FromContext(ctx)

	accessToken, err := create(ctx, user, m.config.JWT.AccessExpiry, m.config.JWT.AccessKey)
	if err != nil {
		return nil, err
	}

	if err = m.save(ctx, accessToken); err != nil {
		logger.WithErr(err).Error("failed to save token to storage")
		return nil, err
	}

	return accessToken, nil
}

// Verify verifies token, and if it presents in storage returns the user.
func (m *JWT) Verify(ctx context.Context, tokenString string) (*User, error) {
	claims, err := m.Parse(tokenString)
	if err != nil {
		return nil, err
	}

	if err := m.CheckIfExists(ctx, claims.Id); err != nil {
		return nil, ierr.NewReason(ierr.ErrTokenExpired)
	}

	return claims.User, nil
}
