package models

import "strings"

const (
	// ErrNotFound is returned when a resource cannot be found
	// in the database
	ErrNotFound modelError = "models: resource not found"

	// ErrPasswordIncorrect is returned when an invalid password
	// is used when attempting to authenticate a user.
	ErrPasswordIncorrect modelError = "models: incorrect password provided"

	// ErrEmailRequired is returned when an email address is not
	// provided when creating an user.
	ErrEmailRequired modelError = "models: Email address is required"
	// ErrEmailInvalid is returned when an email address provided
	// does not match any of our requirements
	ErrEmailInvalid modelError = "models: Email address is not valid"

	// ErrEmailTaken is returned when an update or create is attempted
	// with an email address that is already in use.
	ErrEmailTaken modelError = "models: email address is already taken"

	// ErrPasswordTooShort is returned when an update or create is
	// attempted with a user password that is less than 8 characthers
	ErrPasswordTooShort modelError = "models: password must be at least 8 characthers"

	// ErrPasswordRequired is returned when a create is attempted
	// without a user password provided.
	ErrPasswordRequired    modelError = "models: password is required"
	ErrTitleRequired       modelError = "models: title is required"
	ErrDescriptionRequired modelError = "models: description is required"
	ErrApplyAtRequired     modelError = "models: description is required"
	ErrPwResetInvalid      modelError = "models: token provided is not valid"

	// ErrRememberTooShort is returned when a remember token is
	// not at least 32 bytes
	ErrRememberTooShort privateError = "models: Remember token must be at least 32 bytes"

	// ErrRememberRequired is returned when a create or update is attempted
	// without a user remember token hash.
	ErrRememberRequired   privateError = "models: remember is required"
	ErrUserIDRequired     privateError = "models: user ID is required"
	ErrLocationIDRequired privateError = "models: Location ID is required"
	ErrCategoryIDRequired privateError = "models: Category ID is required"
	// ErrIDInvalid is returned when an invalid ID is provided
	// to a method like Delete.
	ErrIDInvalid privateError = "models: ID provided was invalid"

	ErrServiceRequired privateError = "models: service is required"

	ErrCompanyBenefitRequired modelError = "models: cannot update non existent CompanyBenefit"
	ErrBenefitNameRequired    modelError = "models: benefitName is required"
	ErrCompanyProfileRequired modelError = "models: cannot add benefit to non existent profile"
)

type modelError string

func (e modelError) Error() string {
	return string(e)
}
func (e modelError) Public() string {
	s := strings.Replace(string(e), "models: ", "", 1)
	split := strings.Split(s, " ")
	split[0] = strings.Title(split[0])
	return strings.Join(split, " ")
}

type privateError string

func (e privateError) Error() string {
	return string(e)
}
