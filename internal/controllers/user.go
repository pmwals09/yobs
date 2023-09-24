package controllers

import (
	"errors"
	"html/template"
	"net/http"
	"os"
	"strings"

	helpers "github.com/pmwals09/yobs/internal"
	"github.com/pmwals09/yobs/internal/models/user"
	"golang.org/x/crypto/bcrypt"
)

type FormData struct {
	Errors map[string]template.HTML
	Values map[string]string
}

func HandleRegisterUser(repo user.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		wd, wdErr := os.Getwd()
		if wdErr != nil {
			helpers.WriteError(w, wdErr)
			return
		}

		r.ParseForm()
		newUserInfo := newUserFromRequest(r)
		if errorData, err := validateUserInfo(newUserInfo, repo); err != nil {
			data := FormData{
				Errors: errorData,
				Values: newUserInfo,
			}
			t, err := template.ParseFiles(
				wd + "/web/template/register-user-form-partial.html",
			)
			if err != nil {
				helpers.WriteError(w, err)
				return
			}
			t.ExecuteTemplate(w, "register-user-form-partial", data)
			return
		}

		// Confirm that this is a unique user (username, email)

		u := user.New(newUserInfo["username"], newUserInfo["email"])
		pwHash, err := bcrypt.GenerateFromPassword(
			[]byte(newUserInfo["password"]),
			14,
		)
		if err != nil {
			data := FormData{
				Errors: map[string]template.HTML{
					"password":       template.HTML("<p  class='text-red-600 col-start-2'>Invalid password - please try another</p>"),
					"passwordRepeat": template.HTML("<p class='text-red-600 col-start-2'>Invalid password - please try another</p>"),
				},
				Values: newUserInfo,
			}
			t, err := template.ParseFiles(
				wd + "/web/template/register-user-form-partial.html",
			)
			if err != nil {
				helpers.WriteError(w, err)
				return
			}
			t.ExecuteTemplate(w, "register-user-form-partial", data)
			return
		}
		u.Password = string(pwHash)

		if err = repo.CreateUser(u); err != nil {
			data := FormData{
				Errors: map[string]template.HTML{
					"overall": template.HTML("<p  class='text-red-600 col-start-2'>Error creating user</p>"),
				},
				Values: newUserInfo,
			}
			t, err := template.ParseFiles(
				wd + "/web/template/register-user-form-partial.html",
			)
			if err != nil {
				helpers.WriteError(w, err)
				return
			}
			t.ExecuteTemplate(w, "register-user-form-partial", data)
			return
		}

		// TODO: Success message with link to login page
		data := FormData{
			Errors: map[string]template.HTML{
				"overall": template.HTML("<p  class='text-green-600 col-start-2'>You're registered! <a href='/login'>Login</a> using these credentials.</p>"),
			},
			Values: map[string]string{},
		}
		t, err := template.ParseFiles(
			wd + "/web/template/register-user-form-partial.html",
		)
		if err != nil {
			helpers.WriteError(w, err)
			return
		}
		t.ExecuteTemplate(w, "register-user-form-partial", data)
		return

	}
}

func newUserFromRequest(r *http.Request) map[string]string {
	r.ParseForm()

	newUserInfo := map[string]string{}
	newUserInfo["username"] = r.PostForm.Get("username")
	newUserInfo["email"] = r.PostForm.Get("email")
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
	if potentialUser != nil {
    errorMessage := "Email or username already in use. You can log in <a href='/login'>here.</a>"
		errorData["username"] = template.HTML(
			openTag + errorMessage + "</p>",
		)
		errorData["email"] = template.HTML(
			openTag + errorMessage + "</p>",
		)
		errorMessages = append(errorMessages, errorMessage)
	}

  if err != nil {
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
