package controllers

import (
	"database/sql"
	"errors"
	"html/template"
	"net/http"
	"strings"

	"github.com/google/uuid"
	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/session"
	"github.com/pmwals09/yobs/internal/models/user"
	"golang.org/x/crypto/bcrypt"
  templates "github.com/pmwals09/yobs/web/template"
)

func HandleRegisterUser(repo user.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		newUserInfo := newUserFromRequest(r)
		if errorData, err := validateUserInfo(newUserInfo, repo); err != nil {
			data := helpers.FormData{
				Errors: errorData,
				Values: newUserInfo,
			}
      templates.RegisterUserForm(data).Render(r.Context(), w)
      return
		}

		// Confirm that this is a unique user (username, email)

		u := user.New(newUserInfo["username"], newUserInfo["email"])
		pwHash, err := bcrypt.GenerateFromPassword(
			[]byte(newUserInfo["password"]),
			14,
		)
		if err != nil {
			data := helpers.FormData{
				Errors: map[string]template.HTML{
					"password":       template.HTML("<p  class='text-red-600 col-start-2'>Invalid password - please try another</p>"),
					"passwordRepeat": template.HTML("<p class='text-red-600 col-start-2'>Invalid password - please try another</p>"),
				},
				Values: newUserInfo,
			}
      templates.RegisterUserForm(data).Render(r.Context(), w)
      return
		}
		u.Password = string(pwHash)

		if err = repo.CreateUser(u); err != nil {
			data := helpers.FormData{
				Errors: map[string]template.HTML{
					"overall": template.HTML("<p  class='text-red-600 col-start-2'>Error creating user</p>"),
				},
				Values: newUserInfo,
			}
      templates.RegisterUserForm(data).Render(r.Context(), w)
      return
		}

		// TODO: Success message with link to login page
		data := helpers.FormData{
			Errors: map[string]template.HTML{
				"overall": template.HTML("<p  class='text-green-600 col-start-2'>You're registered! <a href='/login'>Login</a> using these credentials.</p>"),
			},
			Values: map[string]string{},
		}

    templates.RegisterUserForm(data).Render(r.Context(), w)
    return
	}
}

func HandleLogInUser(userRepo user.Repository, sessionRepo session.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		newUserInfo := newUserFromRequest(r)

		user, err := userRepo.GetUserByEmailOrUsername(
			newUserInfo["usernameOrEmail"],
			newUserInfo["usernameOrEmail"],
		)

		if err != nil {
			// TODO: return form w/ messages
			helpers.WriteError(w, err)
			return
		}

		if err = bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(newUserInfo["password"]),
		); err != nil {
			// TODO: return form w/ messages
			helpers.WriteError(w, err)
			return
		}

    // TODO: I don't like this arrangement. It would be really easy to create
    // a session object, but not insert it into the database before sending a
    // cookie to the browser
    s := session.New()
    s.WithUser(user)
    sessionRepo.CreateSession(s)
    cName := "yobs"
    c := http.Cookie{
      Name: cName,
      Value: s.UUID.String(),
      Path: "/",
      Expires: s.Expiration,
      HttpOnly: true,
      Secure: true,
      SameSite: http.SameSiteDefaultMode,
    }
    
    // NOTE: I was running into an issue where certain requests from HTMX would
    // have 2 "yobs" cookies - the current one and an old one. We should only
    // have 1 session running at a time, so we will remove the old cookie and
    // remove the old session from the database before setting the new cookie
    currentCookies := r.Cookies()
    for _, cookie := range currentCookies {
      if cookie.Name == cName {
        cookie.MaxAge = -1
        http.SetCookie(w, cookie)
        oldUuid, err := uuid.Parse(cookie.Value)
        if err != nil {
          helpers.WriteError(w, err)
          return
        }
        err = sessionRepo.DeleteSessionByUUID(oldUuid)
        if err != nil {
          helpers.WriteError(w, err)
          return
        }
        break
      }
    }
    http.SetCookie(w, &c)

    // NOTE: This seems unintuitive, but it works with the HTMX model. One may
    // expect to do http.Redirect(w, r, "/home", http.StatusFound), but that
    // will replace the hx-boost area with the contents of the page to which
    // we redirect. Done enough times, this is like the "infinite mirror"
    // effect. In retrospect, a 302 forwards a GET request to the new path, and
    // the server returns the full HTML of the page to replace the boosted area,
    // so it makes sense. This is the alternative proposed by HTMX
    w.Header().Add("HX-Redirect", "/home")
    return
	}
}

func newUserFromRequest(r *http.Request) map[string]string {
	r.ParseForm()

	newUserInfo := map[string]string{}
	newUserInfo["username"] = r.PostForm.Get("username")
	newUserInfo["email"] = r.PostForm.Get("email")
	newUserInfo["usernameOrEmail"] = r.PostForm.Get("username-or-email")
	newUserInfo["password"] = r.PostForm.Get("password")
	newUserInfo["passwordRepeat"] = r.PostForm.Get("password-repeat")
	return newUserInfo
}

func validateUserInfo(userInfo map[string]string, repo user.Repository) (map[string]template.HTML, error) {
	errorData := map[string]template.HTML{}
	errorMessages := []string{}
	openTag := "<p class='text-red-600 col-start-2'>"

	// Confirm that the password fields match
	if userInfo["password"] != userInfo["passwordRepeat"] {
		errorMessage := "Password fields must match"
		errorData["password"] = template.HTML(
			openTag + errorMessage + "</p>",
		)
		errorData["passwordRepeat"] = template.HTML(
			openTag + errorMessage + "</p>",
		)
		errorMessages = append(errorMessages, errorMessage)
	}

	// Confirm that the email and username are unique
	potentialUser, err := repo.GetUserByEmailOrUsername(
		userInfo["email"],
		userInfo["username"],
	)
	if potentialUser.ID != 0 {
		errorMessage := "Email or username already in use. You can log in <a href='/login'>here.</a>"
		errorData["username"] = template.HTML(
			openTag + errorMessage + "</p>",
		)
		errorData["email"] = template.HTML(
			openTag + errorMessage + "</p>",
		)
		errorMessages = append(errorMessages, errorMessage)
	}

	if err != nil && err != sql.ErrNoRows {
		errorMessage := "Error validating username and email"
		errorData["username"] = template.HTML(
			openTag + errorMessage + "</p>",
		)
		errorData["email"] = template.HTML(
			openTag + errorMessage + "</p>",
		)
		errorMessages = append(errorMessages, errorMessage)
	}

	if len(errorMessages) == 0 {
		return nil, nil
	}
	return errorData, errors.New(strings.Join(errorMessages, ", "))

}
