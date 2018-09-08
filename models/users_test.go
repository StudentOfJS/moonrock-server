package models


var testUsers = []User{
  {
    Address:         "1 Chester Field Green, Baltimore Fields, Baltimore, MA",
    Confirmed:       false,
    CountryCode:     "US",
    EthereumAddress: "0xe81D72D14B1516e68ac3190a46C93302Cc8eD60f",
    FirstName:       "Teddy",
    Group:           "public_investor",
    LastName:        "Weinstein",
    Password:        "TotalMayhem",
    ResetCode:       resetcode,
    Username:        "teddy.w@test.com",
  },
  {
    Address:         "12 Bacon Court, Saltash, Cornwall, UK",
    Confirmed:       false,
    CountryCode:     "UK",
    EthereumAddress: "0x6a068E0287e55149a2a8396cbC99578f9Ad16A31",
    FirstName:       "Dave",
    Group:           "public_investor",
    LastName:        "Saville",
    Password:        "Loser",
    ResetCode:       resetcode,
    Username:        "dave@test.com",
  },
  {
    Address:         "82 Avalon Plains, Esperance, WA, Australia",
    Confirmed:       false,
    CountryCode:     "AU",
    EthereumAddress: "0xe81D72D14B1516e68ac3190a46C93302Cc8eE60c",
    FirstName:       "Brad",
    Group:           "public_investor",
    LastName:        "Tad",
    Password:        "SurfOrDie2",
    ResetCode:       resetcode,
    Username:        "tad@test.com",
  },
  {
    Address:         "8 Tornado Alley, Aliceville, Wisconsin",
    Confirmed:       false,
    CountryCode:     "US",
    EthereumAddress: "0xe81D72D14B1516e68ac3190a46C93302Cc8eD60f",
    FirstName:       "Avril",
    Group:           "public_investor",
    LastName:        "Smith",
    Password:        "fhweuhwriwe34",
    ResetCode:       resetcode,
    Username:        "a.s@test.com",
  },
  {
    Address:         "Fisherman's Cottage, Smugglers Cove, Turks and Cacos",
    Confirmed:       false,
    CountryCode:     "TC",
    EthereumAddress: "0x595832F8FC6BF59c85C527fEC3740A1b7a361269",
    FirstName:       "Peter",
    Group:           "public_investor",
    LastName:        "Marston",
    Password:        "r4j3ok4j50f",
    ResetCode:       resetcode,
    Username:        "peter@test.com",
  },
}



// Register validates the user signup form and saves to db
func Register(t *testing.T){

  a, c, e, f, l, p, u string
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