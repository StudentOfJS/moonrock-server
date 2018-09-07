package models

import (
	"strconv"

	"github.com/asdine/storm"
	"github.com/satori/go.uuid"
	"github.com/studentofjs/moonrock-server/database"
	"github.com/studentofjs/moonrock-server/mailer"
)

// User struct contains all the user data
type User struct {
	Address         string    // this field will not be indexed
	Confirmed       bool      // this field will not be indexed
	CountryCode     string    // this field will not be indexed
	EthereumAddress string    // this field will not be indexed
	FirstName       string    // this field will not be indexed
	Group           string    `storm:"index"`        // this field will be indexed
	ID              int       `storm:"id,increment"` // primary key with auto increment
	LastName        string    // this field will not be indexed
	Password        []byte    // this field will not be indexed
	ResetCode       uuid.UUID // this field will not be indexed
	Username        string    `storm:"unique"` // this field will be indexed with a unique constraint
}

// UpdateContributionAddress uses an ID to find user and updates their contribution address
func UpdateContributionAddress(id int, e string) *Response {
	db, err := database.OpenDB()
	if err != nil {
		return getResponse("server error")
	}
	defer db.Close()
	if err := db.UpdateField(&User{ID: id}, "EthereumAddress", e); err != nil {
		return getResponse("invalid address")
	}
	return getResponse("ok")
}

// ConfirmAccount checks a resetCode against the DB and returns an error string or
func ConfirmAccount(c string) *Response {
	rc, _ := uuid.FromString(c)
	db, err := database.OpenDB()
	if err != nil {
		return getResponse("server error")
	}
	defer db.Close()
	var user User
	if err := db.One("ResetCode", rc, &user); err != nil {
		return getResponse("user doesn't exist")
	}
	if err := db.UpdateField(&User{ID: user.ID}, "Confirmed", true); err != nil {
		return getResponse("token expired, please try again")
	}
	return getResponse("account successfully confirmed")
}

// ForgotPassword sends a reset email with unique password reset link
func ForgotPassword(u string) *Response {
	resetcode := uuid.Must(uuid.NewV4())
	rc := resetcode.String()

	db, err := database.OpenDB()
	if err != nil {
		return getResponse("server error")
	}
	defer db.Close()

	var user User
	if err := db.One("Username", u, &user); err != nil {
		return getResponse("invalid login")
	}
	if err := db.UpdateField(&User{ID: user.ID}, "ResetCode", resetcode); err != nil {
		return getResponse("server error")
	}

	r := mailer.NewRequest([]string{u}, "Moonrock password reset")
	r.Send("templates/reset_template.html", map[string]string{
		"reset":    rc,
		"username": u,
	})
	return getResponse("check your email")
}

// GetContributionAddress returns the saved address of the user
func GetContributionAddress(i string) (string, *Response) {
	var user User
	db, err := database.OpenDB()
	if err != nil {
		return "", getResponse("server error")
	}
	defer db.Close()
	id, e := strconv.Atoi(i)
	if e != nil {
		return "", getResponse("unauthenticated")
	}
	err = db.One("ID", id, &user)
	if err != nil {
		return "", getResponse("user doesn't exist")
	}

	return user.EthereumAddress, getResponse("ok")
}

// Register validates the user signup form and saves to db
func Register(a, c, e, f, l, p, u string) *Response {
	resetcode := uuid.Must(uuid.NewV4())

	if err := LoginValid(u, p); err != nil {
		return getResponse("invalid signup")
	}

	if err := UserValid(e, f, l); err != nil {
		return getResponse("invalid signup")
	}
	// Generate "hash" to store from username password
	hash, err := HashPassword(p)
	if err != nil {
		return getResponse("server error")
	}

	user := User{
		Address:         a,
		Confirmed:       false,
		CountryCode:     c,
		EthereumAddress: e,
		FirstName:       f,
		Group:           "public_investor",
		LastName:        l,
		Password:        hash,
		ResetCode:       resetcode,
		Username:        u,
	}
	// Start boltDB
	db, err := database.OpenDB()
	if err != nil {
		return getResponse("server error")
	}
	defer db.Close()
	if err := db.Save(&user); err == storm.ErrAlreadyExists {
		return getResponse("already signed up")
	}

	r := mailer.NewRequest([]string{u}, "Moonrock Account Confirmation")
	r.Send("templates/register_template.html", map[string]string{
		"country":  c,
		"ethereum": e,
		"name":     f,
	})
	return getResponse("ok")
}

// ResetPassword handles the reset code checking and password change
func ResetPassword(p, r, u string) *Response {
	// Generate "hash" from password
	hash, err := HashPassword(p)
	if err != nil {
		return getResponse("server error")
	}

	reset, e := uuid.FromString(r)
	if e != nil {
		return getResponse("server error")
	}

	db, err := database.OpenDB()
	if err != nil {
		return getResponse("server error")
	}
	defer db.Close()

	var user User
	err = db.One("Username", u, &user)
	if err != nil {
		return getResponse("invalid")
	}
	if uuid.Equal(user.ResetCode, reset) {
		if err := db.UpdateField(&User{Username: u}, "Password", hash); err != nil {
			return getResponse("token expired, please try again")
		}
		newResetCode := uuid.Must(uuid.NewV4())
		db.UpdateField(&User{ResetCode: reset}, "ResetCode", newResetCode)
		return getResponse("ok")
	}
	return getResponse("token expired, please try again")
}

// UpdateUserDetails updates user details supplied to API
func UpdateUserDetails(a, c, f, i, l string) *Response {
	id, _ := strconv.Atoi(i)

	db, err := database.OpenDB()
	if err != nil {
		return getResponse("server error")
	}
	defer db.Close()

	if err := db.Update(&User{
		ID:          id,
		Address:     a,
		CountryCode: c,
		FirstName:   f,
		LastName:    l,
	}); err != nil {
		return getResponse("server error")
	}
	return getResponse("ok")
}
