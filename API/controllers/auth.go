package controllers

import (
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/samueldaviddelacruz/go-job-board/API/email"
	"github.com/samueldaviddelacruz/go-job-board/API/models"
	"io/ioutil"
	"net/http"
	"time"
)

// NewUsers is used to create a new Users controller.
// This function will panic if the templates are not
// parsed correctly, and should only be used during
// initial setup
func NewAuth(us models.UserService, emailer *email.Client) *Auth {
	return &Auth{
		us:      us,
		emailer: emailer,
	}
}

// Users Represents a Users controller
type Auth struct {
	us      models.UserService
	emailer *email.Client
}

// Create is used to process the signup form when a user
// submits it. This is used to create a new user account.
//
func (u *Auth) Create(w http.ResponseWriter, r *http.Request) {

	credentials := Credentials{
	}

	err := parseJSON(r, &credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	companyUser := models.User{
		RoleID:   1,
		Password: credentials.Password,
		Email:    credentials.Email,
	}

	if err := u.us.Create(&companyUser); err != nil {
		//vd.SetAlert(err)
		respondJSON(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, "resource created successfully")
}

// Login is used to verify the provided email address and
// password and then log the user in if they are correct.
//
// POST /login
func (u *Auth) Login(w http.ResponseWriter, r *http.Request) {
	credentials := Credentials{
	}

	err := parseJSON(r, &credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	companyUser := models.User{
		Password: credentials.Password,
		Email:    credentials.Email,
	}
	_, err = u.us.Authenticate(companyUser.Email, companyUser.Password)
	if err != nil {
		switch err {
		case models.ErrNotFound:
			respondJSON(w, http.StatusUnauthorized, "Invalid email address")
		default:
			//vd.SetAlert(err)
			respondJSON(w, http.StatusInternalServerError, err)
		}
		//u.LoginView.Render(w, r, vd)
		return
	}

	token, err := u.signIn(w, &companyUser)
	if err != nil {
		//u.LoginView.Render(w, r, vd)
		respondJSON(w, http.StatusInternalServerError, err)
		return
	}
	//user2, _ := u.us.ByID(1)
	respondJSON(w, http.StatusOK, string(token))
}

// ResetPwForm is used to process the forgot password form
// and the reset password form
type ResetPwForm struct {
	Email    string `schema:"email"`
	Token    string `schema:"token"`
	Password string `schema:"password"`
}

// POST /forgot
func (u *Auth) InitiateReset(w http.ResponseWriter, r *http.Request) {
	//var vd views.Data
	var form ResetPwForm
	//vd.Yield = &form
	if err := parseForm(r, &form); err != nil {
		//vd.SetAlert(err)
		//u.ForgotPwView.Render(w, r, vd)
		return
	}
	token, err := u.us.InitiateReset(form.Email)
	if err != nil {
		//vd.SetAlert(err)
		//u.ForgotPwView.Render(w, r, vd)
		return
	}

	err = u.emailer.ResetPw(form.Email, token)
	if err != nil {
		//vd.SetAlert(err)
		//u.ForgotPwView.Render(w, r, vd)
		return
	}
	/*
		views.RedirectAlert(w, r, "/reset", http.StatusFound, views.Alert{
			Level:   views.AlertLvlSuccess,
			Message: "Instructions for resetting your password have been emailed to you.",
		})
	*/
}

// CompleteReset processes the reset password form
//
//POST
func (u *Auth) CompleteReset(w http.ResponseWriter, r *http.Request) {
	/*
		//var vd views.Data
		var form ResetPwForm
		//vd.Yield = &form
		if err := parseForm(r, &form); err != nil {
			//vd.SetAlert(err)
			//u.ResetPwView.Render(w, r, vd)
			return
		}
		user, err := u.us.CompleteReset(form.Token, form.Password)
		if err != nil {
			//vd.SetAlert(err)
			//u.ResetPwView.Render(w, r, vd)
			return
		}

		err = u.signIn(w, user)
		if err != nil {
			//vd.SetAlert(err)
			//u.LoginView.Render(w, r, vd)
			return
		}

			views.RedirectAlert(w, r, "/galleries", http.StatusFound, views.Alert{
				Level:   views.AlertLvlSuccess,
				Message: "Your password has been reset and you have been logged in!",
			})
	*/
}

func (u *Auth) signIn(w http.ResponseWriter, user *models.User) ([]byte, error) {
	privateKey, _ := ioutil.ReadFile("./key.priv")
	var hs = jwt.NewHS256(privateKey)
	now := time.Now()
	oneYear := 24 * 30 * 12 * time.Hour
	pl := models.CustomPayload{
		Payload: jwt.Payload{
			Issuer: "Go JobBoard",
			//Subject:        "someone",
			//Audience:       jwt.Audience{"https://golang.org", "https://jwt.io"},
			ExpirationTime: jwt.NumericDate(now.Add(oneYear)),
			NotBefore:      jwt.NumericDate(now.Add(30 * time.Minute)),
			IssuedAt:       jwt.NumericDate(now),
		},
		Email: user.Email,
	}

	token, err := jwt.Sign(pl, hs)
	if err != nil {
		return nil, err
	}
	//jwt := jws.NewJWT()
	return token, nil
}
