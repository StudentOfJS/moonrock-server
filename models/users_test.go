package models

import (
	"fmt"
	"testing"

	"github.com/asdine/storm"
	uuid "github.com/google/uuid"
	"github.com/studentofjs/moonrock-server/database"
	"golang.org/x/crypto/bcrypt"
)

type TestResetCodes struct {
	ID         string `storm:"index"`
	ResetCodes map[int]string
}

type testCompleteUser struct {
	address   string
	confirmed bool
	country   string
	eth       string
	firstname string
	group     string
	id        int
	lastname  string
	password  string
	user      string
}

var testCompleteUsers = []testCompleteUser{
	{
		id:        1,
		address:   "1 Chester Field Green, Baltimore Fields, Baltimore, MA",
		confirmed: false,
		country:   "US",
		eth:       "0xe81D72D14B1516e68ac3190a46C93302Cc8eD60f",
		firstname: "Teddy",
		group:     "public_investor",
		lastname:  "Weinstein",
		password:  "TotalMayhemw13",
		user:      "teddy@test.com",
	},
	{
		id:        2,
		address:   "12 Bacon Court, Saltash, Cornwall, UK",
		confirmed: false,
		country:   "UK",
		eth:       "0x6a068E0287e55149a2a8396cbC99578f9Ad16A31",
		firstname: "Dave",
		group:     "public_investor",
		lastname:  "Saville",
		password:  "Loser322452",
		user:      "dave@test.com",
	},
	{
		id:        3,
		address:   "82 Avalon Plains, Esperance, WA, Australia",
		confirmed: false,
		country:   "AU",
		eth:       "0xe81D72D14B1516e68ac3190a46C93302Cc8eE60c",
		firstname: "Brad",
		group:     "public_investor",
		lastname:  "Tad",
		password:  "SurfOrDie2",
		user:      "tad@test.com",
	},
	{
		id:        4,
		address:   "8 Tornado Alley, Aliceville, Wisconsin",
		confirmed: false,
		country:   "US",
		eth:       "0xe81D72D14B1516e68ac3190a46C93302Cc8eD60f",
		firstname: "Avril",
		group:     "public_investor",
		lastname:  "Smith",
		password:  "fhweuhJwriwe34",
		user:      "al@test.com",
	},
	{
		id:        5,
		address:   "Fisherman's Cottage, Smugglers Cove, Turks and Cacos",
		confirmed: false,
		country:   "TC",
		eth:       "0x595832F8FC6BF59c85C527fEC3740A1b7a361269",
		firstname: "Peter",
		group:     "public_investor",
		lastname:  "Marston",
		password:  "r4j3ok4j50f",
		user:      "peter@test.com",
	},
}

var resetCodes = make(map[int]string)

// Register validates the user signup form and saves to db
func TestValidRegister(t *testing.T) {
	db, err := database.OpenTestDB("../database/")
	if err != nil {
		t.Error("server error")
		return
	}
	defer db.Close()
	for i, r := range testCompleteUsers {
		if err := LoginValid(r.user, r.password); err != nil {
			t.Errorf("invalid username or password %d %v", r.id, err)
			return
		}
		if err := UserValid(r.eth, r.firstname, r.lastname); err != nil {
			t.Error("invalid signup details")
			return
		}

		hash, err := HashPassword(r.password)
		if err != nil {
			t.Error("server error")
			return
		}
		reset := uuid.New().String()
		resetCodes[i] = reset

		user := User{
			Address:         r.address,
			Confirmed:       r.confirmed,
			CountryCode:     r.country,
			EthereumAddress: r.eth,
			FirstName:       r.firstname,
			Group:           r.group,
			LastName:        r.lastname,
			Password:        hash,
			ResetCode:       reset,
			Username:        r.user,
		}
		if err := db.Save(&user); err == storm.ErrAlreadyExists {
			t.Error("user already signed up")
			return
		}
	}
	testResetCodes := TestResetCodes{
		ID:         "test",
		ResetCodes: resetCodes,
	}
	if err := db.Save(&testResetCodes); err == storm.ErrAlreadyExists {
		t.Error("user already signed up")
		return
	}

}

func TestResestCodesCreated(t *testing.T) {
	for k, v := range resetCodes {
		fmt.Printf("key: %d, value: %s", k, v)
		t.Errorf("key: %d, value: %s", k, v)
	}
}

func TestConfirmAccount(t *testing.T) {
	db, err := database.OpenTestDB("../database/")
	if err != nil {
		t.Error("server error")
		return
	}
	defer db.Close()
	for _, v := range resetCodes {
		var user User
		if err := db.One("ResetCode", v, &user); err != nil {
			t.Errorf("failed searching user by reset code: %v", err)
			return
		}
		if err := db.UpdateField(&User{ID: user.ID}, "Confirmed", true); err != nil {
			t.Error("failed trying to update user to confirmed true")
			return
		}
	}
}

func TestUpdateContributionAddress(t *testing.T) {
	db, err := database.OpenTestDB("../database/")
	if err != nil {
		t.Error("server error")
		return
	}
	defer db.Close()

	for _, u := range testCompleteUsers {
		if err := db.UpdateField(&User{ID: u.id}, "EthereumAddress", "0xCaE9eFE97895EF43e72791a10254d6abDdb17Ae9"); err != nil {
			t.Error("failed to update eth address")
			return
		}
	}
}

func TestResetPassword(t *testing.T) {
	db, err := database.OpenTestDB("../database/")
	if err != nil {
		t.Error("opening test db failed")
		return
	}
	defer db.Close()
	hash, err := HashPassword("this_is_a_test")
	if err != nil {
		t.Error("password hashing failed")
		return
	}

	var testResetCodes TestResetCodes
	if err := db.One("ID", "test", &testResetCodes); err != nil {
		t.Errorf("failed searching user by reset code: %v", err)
		return
	}

	for i, u := range testCompleteUsers {
		var user User
		err = db.One("Username", u.user, &user)
		if err != nil {
			t.Error("can't locate user in db")
			return
		}

		if user.ResetCode == testResetCodes.ResetCodes[i] {
			if err := db.UpdateField(&User{ID: u.id}, "Password", hash); err != nil {
				t.Errorf("updating password field failed, with error: %v", err)
				return
			}
			reset := uuid.New()
			newResetCode := reset.String()
			db.UpdateField(&User{ID: u.id}, "ResetCode", newResetCode)
		} else {
			t.Errorf("reset codes not equal: %s vs %s", user.ResetCode, testResetCodes.ResetCodes[i])
			return
		}
	}
}

// UpdateUserDetails updates user details supplied to API
func TestUpdateUserDetails(t *testing.T) {
	db, err := database.OpenTestDB("../database/")
	if err != nil {
		t.Error("opening test db failed")
		return
	}
	defer db.Close()

	for _, u := range testCompleteUsers {
		var user User
		if err = db.One("ID", u.id, &user); err != nil {
			t.Error("cannot find user")
			return
		}
		user.Address = "test_address"
		user.CountryCode = "AU"
		user.FirstName = "test"
		user.LastName = "change"

		if err := db.Update(&user); err != nil {
			t.Errorf("Updating user details failed with error: %v", err)
			return
		}

		if err := db.One("ID", u.id, &user); err != nil {
			t.Error("User not searchable after update")
			return
		}
		if user.Username != u.user || user.Group != u.group {
			t.Errorf("non updated fields mutated during update, : %s vs %s, %s vs %s", user.Username, u.user, user.Group, u.group)
			return
		}

		if user.Address != "test_address" || user.FirstName != "test" {
			t.Error("update occured without changing requsted fields")
			return
		}
	}
}

func TestDeleteUsers(t *testing.T) {
	db, err := database.OpenTestDB("../database/")
	if err != nil {
		t.Error("opening test db failed")
		return
	}
	defer db.Close()

	for _, u := range testCompleteUsers {
		var user User
		if err := db.One("ID", u.id, &user); err != nil {
			t.Error("failed trying to delete, user not found")
			return
		}
		if err := bcrypt.CompareHashAndPassword(user.Password, []byte("this_is_a_test")); err != nil {
			t.Error("failed trying to delete, passwords don't match")
			return
		}
		if err := db.DeleteStruct(&user); err != nil {
			t.Errorf("failed trying to delete, delete struct failed with error: %v", err)
			return
		}
	}
	var users []User
	if err := db.Range("ID", 1, 5, &users); err != nil {
		fmt.Println(err.Error())
	}
	leftOver := len(users)
	if leftOver > 0 {
		t.Errorf("not all users deleted, %d remain", leftOver)
		return
	}

}
