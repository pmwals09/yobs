package controllers

import (
	"fmt"
	"net/http"

	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/user"
	"github.com/pmwals09/yobs/web/profile"
)

func HandleGetProfilePage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := r.Context().Value(user.UserCtxKey).(*user.User)
		pa := helpers.ProfileArgs{
			Username: u.Username,
			Email:    u.Email,
		}
		profilepage.ProfilePage(u, pa).Render(r.Context(), w)
	}
}

func HandleGetBasicProfileForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, ok := r.Context().Value(user.UserCtxKey).(*user.User)
		if !ok {
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			return
		}
		var fd helpers.FormData
		fd.AddValue("profile-email", u.Email)
		fd.AddValue("profile-username", u.Username)

		profilepage.ProfilePageForm(fd).Render(r.Context(), w)
	}
}

func HandleUpdateProfile(userRepo user.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUser user.User
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("<p class='text-red-600>Error parsing form</p>"))
			return
		}
		u, ok := r.Context().Value(user.UserCtxKey).(*user.User)
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("<p class='text-red-600>Error getting valid user</p>"))
			return
		}

		newUser.Email = r.PostForm.Get("profile-email")
		newUser.Username = r.PostForm.Get("profile-username")
		newUser.ID = u.ID
		newUser.UUID = u.UUID
		newUser.Password = u.Password

		fmt.Printf("%+v\n", newUser)
		err = userRepo.UpdateUser(&newUser)
		if err != nil {
			var fd helpers.FormData
			fd.AddError("overall", "Error updating user information")
			profilepage.ProfilePageForm(fd).Render(r.Context(), w)
			return
		}
		var pa helpers.ProfileArgs
		pa.Email = newUser.Email
		pa.Username = newUser.Username
		profilepage.BasicProfile(pa).Render(r.Context(), w)
	}
}
