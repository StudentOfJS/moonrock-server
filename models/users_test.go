package models

import (
	"testing"

	"github.com/asdine/storm"
	uuid "github.com/satori/go.uuid"
	"github.com/studentofjs/moonrock-server/database"
)

type testUser struct {
	address   string
	confirmed bool
	country   string
	eth       string
	firstname string
	group     string
	id        int
	lastname  string
	password  string
	reset     uuid.UUID
	user      string
}

var testUsers = []testUser{
	{
		address:   "1 Chester Field Green, Baltimore Fields, Baltimore, MA",
		confirmed: false,
		country:   "US",
		eth:       "0xe81D72D14B1516e68ac3190a46C93302Cc8eD60f",
		firstname: "Teddy",
		group:     "public_investor",
		lastname:  "Weinstein",
		password:  "TotalMayhem",
		reset:     uuid.Must(uuid.NewV4()),
		user:      "teddy.w@test.com",
	},
	{
		address:   "12 Bacon Court, Saltash, Cornwall, UK",
		confirmed: false,
		country:   "UK",
		eth:       "0x6a068E0287e55149a2a8396cbC99578f9Ad16A31",
		firstname: "Dave",
		group:     "public_investor",
		lastname:  "Saville",
		password:  "Loser",
		reset:     uuid.Must(uuid.NewV4()),
		user:      "dave@test.com",
	},
	{
		address:   "82 Avalon Plains, Esperance, WA, Australia",
		confirmed: false,
		country:   "AU",
		eth:       "0xe81D72D14B1516e68ac3190a46C93302Cc8eE60c",
		firstname: "Brad",
		group:     "public_investor",
		lastname:  "Tad",
		password:  "SurfOrDie2",
		reset:     uuid.Must(uuid.NewV4()),
		user:      "tad@test.com",
	},
	{
		address:   "8 Tornado Alley, Aliceville, Wisconsin",
		confirmed: false,
		country:   "US",
		eth:       "0xe81D72D14B1516e68ac3190a46C93302Cc8eD60f",
		firstname: "Avril",
		group:     "public_investor",
		lastname:  "Smith",
		password:  "fhweuhwriwe34",
		reset:     uuid.Must(uuid.NewV4()),
		user:      "a.s@test.com",
	},
	{
		address:   "Fisherman's Cottage, Smugglers Cove, Turks and Cacos",
		confirmed: false,
		country:   "TC",
		eth:       "0x595832F8FC6BF59c85C527fEC3740A1b7a361269",
		firstname: "Peter",
		group:     "public_investor",
		lastname:  "Marston",
		password:  "r4j3ok4j50f",
		reset:     uuid.Must(uuid.NewV4()),
		user:      "peter@test.com",
	},
}

// Register validates the user signup form and saves to db
func TestValidRegister(t *testing.T) {

	for _, r := range testUsers {
		if err := LoginValid(r.user, r.password); err != nil {
			t.Error("invalid username or password")
		}
		if err := UserValid(r.eth, r.firstname, r.lastname); err != nil {
			t.Error("invalid signup details")
		}

		hash, err := HashPassword(r.password)
		if err != nil {
			t.Error("server error")
		}

		user := User{
			Address:         r.address,
			Confirmed:       r.confirmed,
			CountryCode:     r.country,
			EthereumAddress: r.eth,
			FirstName:       r.firstname,
			Group:           r.group,
			LastName:        r.lastname,
			Password:        hash,
			ResetCode:       r.reset,
			Username:        r.user,
		}

		db, err := database.OpenTestDB()
		if err != nil {
			t.Error("server error")
		}
		defer db.Close()
		if err := db.Save(&user); err == storm.ErrAlreadyExists {
			t.Error("user already signed up")
		}

	}
}
