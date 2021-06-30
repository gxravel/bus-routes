package jwt

import (
	"context"
	"time"

	v1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"
	"github.com/gxravel/bus-routes/internal/config"
	"github.com/gxravel/bus-routes/internal/storage"

	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

// Manager includes the methods allowed to deal with the token.
type Manager interface {
	save(ctx context.Context, token *Details) error

	Parse(tokenString string) (*Claims, error)
	CheckIfExists(ctx context.Context, tokenUUID string) error
	Delete(ctx context.Context, tokenUUID string) error
	SetNew(ctx context.Context, user *v1.User) (*v1.Token, error)
	Verify(ctx context.Context, tokenString string) (*v1.User, error)
}

// Claims defines JWT token claims.
type Claims struct {
	User *v1.User `json:"user"`
	jwt.StandardClaims
}

// Details defines the structure of a JWT token.
type Details struct {
	String  string
	Expiry  int64
	UUID    string
	Subject string
}

// JWT contains the fields which interact with the token.
type JWT struct {
	client *storage.Client
	config config.JWT
}

func New(client *storage.Client, config config.JWT) *JWT {
	return &JWT{client: client, config: config}
}

// create creates the HS512 JWT token with claims.
func create(user *v1.User, expiry time.Duration, key string) (*Details, error) {
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

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	var err error
	token.String, err = jwtToken.SignedString([]byte(key))
	return token, err
}

// Parse parses a string token with the key.
func (m *JWT) Parse(tokenString string) (*Claims, error) {
	var key = []byte(m.config.AccessKey)

	jwtToken, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return key, nil
	})

	claims, ok := jwtToken.Claims.(*Claims)
	if !ok || !jwtToken.Valid {
		return nil, errors.New("couldn't handle this token: " + err.Error())
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
func (m *JWT) SetNew(ctx context.Context, user *v1.User) (*v1.Token, error) {
	accessToken, err := create(user, m.config.AccessExpiry, m.config.AccessKey)
	if err != nil {
		return nil, err
	}
	err = m.save(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	token := &v1.Token{
		Token:  accessToken.String,
		Expiry: accessToken.Expiry,
	}

	return token, nil
}

// Verify returns the user.
func (m *JWT) Verify(ctx context.Context, tokenString string) (*v1.User, error) {
	claims, err := m.Parse(tokenString)
	if err != nil {
		return nil, err
	}
	if err := m.CheckIfExists(ctx, claims.Id); err != nil {
		return nil, errors.New("token expired")
	}

	return claims.User, nil
}
