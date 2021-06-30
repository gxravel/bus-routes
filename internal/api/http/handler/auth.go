package handler

import (
	"net/http"
	"regexp"
	"strings"

	api "github.com/gxravel/bus-routes/internal/api/http"
	v1 "github.com/gxravel/bus-routes/internal/api/http/handler/v1"
	"github.com/gxravel/bus-routes/internal/dataprovider"
	"github.com/gxravel/bus-routes/internal/model"

	"github.com/pkg/errors"
)

var (
	regPass  = regexp.MustCompile(`^.{4,255}$`)
	regEmail = regexp.MustCompile(`^([a-zA-Z0-9_\-\.]+)@([a-zA-Z0-9_\-\.]+)\.([a-zA-Z]{2,5})$`)
)

func validateUserCredentials(user *v1.User) (err error) {
	if !regPass.MatchString(user.Password) {
		err = errors.New("invalid password: min length - 4")
		return
	}
	if !regEmail.MatchString(user.Email) {
		err = errors.New("invalid email")
		return
	}
	user.Email = strings.ToLower(user.Email)
	return
}

func (s *Server) signup(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user = &v1.User{}
	if err := s.processRequest(r, user); err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if err := validateUserCredentials(user); err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if err := s.busroutes.AddUsers(ctx, user); err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			api.RespondError(ctx, w, http.StatusConflict, errors.New("the email was already taken"))
			return
		}
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	user.Type = model.DefaultUserType

	token, err := s.busroutes.NewJWT(ctx, user)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondDataOK(ctx, w, token)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user = &v1.User{}
	if err := s.processRequest(r, user); err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if err := validateUserCredentials(user); err != nil {
		api.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	filter := dataprovider.NewUserFilter().SelectPassword().ByEmails(user.Email)
	truePassword, err := s.busroutes.CheckPasswordHash(ctx, user.Password, filter)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}
	if !truePassword {
		api.RespondError(ctx, w, http.StatusUnauthorized, errors.New("wrong credentials"))
		return
	}

	filter = dataprovider.NewUserFilter().SelectType().ByEmails(user.Email)
	user.Type, err = s.busroutes.GetUserType(ctx, filter)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	token, err := s.busroutes.NewJWT(ctx, user)
	if err != nil {
		api.RespondError(ctx, w, http.StatusInternalServerError, err)
		return
	}

	api.RespondDataOK(ctx, w, token)
}
