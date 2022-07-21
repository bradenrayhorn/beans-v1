package http

import (
	"encoding/json"
	"net/http"

	"github.com/bradenrayhorn/beans/beans"
)

func (s *Server) handleUserRegister() http.HandlerFunc {
	type request struct {
		Username beans.Username `json:"username"`
		Password beans.Password `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			Error(w, err)
			return
		}

		_, err := s.userService.CreateUser(r.Context(), req.Username, req.Password)
		if err != nil {
			Error(w, err)
			return
		}
	}
}

func (s *Server) handleUserLogin() http.HandlerFunc {
	type request struct {
		Username beans.Username `json:"username"`
		Password beans.Password `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			Error(w, err)
			return
		}

		user, err := s.userService.Login(r.Context(), req.Username, req.Password)
		if err != nil {
			Error(w, err)
			return
		}

		session, err := s.sessionRepository.Create(user.ID)
		if err != nil {
			Error(w, beans.WrapError(err, beans.ErrorInternal))
			return
		}

		cookie := http.Cookie{
			Name:     "session_id",
			Value:    string(session.ID),
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
		}

		http.SetCookie(w, &cookie)
	}
}
