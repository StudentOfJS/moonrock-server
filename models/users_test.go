package models

type testLogin struct {
	username string
	password string
	valid    bool
}

type testUser struct {
	ethereum  string
	firstname string
	lastname  string
	valid     bool
}

type testEmail struct {
	email string
	valid bool
}

var testLogins = []testLogin{
	{username: "test@test.com", password: "12uhuh4rf89J", valid: true},
	{username: "test2@test.com.au", password: "21268634238432dss", valid: true},
	{username: "test@test.co", password: "uhuhsdfs4rf89J", valid: true},
	{username: "test@test.net.au", password: "needtobreakFree", valid: true},
	{username: "@test.com", password: "12uhuh4rf89J", valid: false},
	{username: "test2@testcom", password: "21268634238432dss", valid: false},
	{username: "test@.co", password: "uhuhsdfs4rf89J", valid: false},
	{username: "test@test.net.au", password: "n", valid: false},
}

var testUsers = []testUser{
	{ethereum: "0xe81D72D14B1516e68ac3190a46C93302Cc8eD60f", firstname: "coin", lastname: "lancer", valid: true},
	{ethereum: "0x595832F8FC6BF59c85C527fEC3740A1b7a361269", firstname: "Power", lastname: "Ledger", valid: true},
	{ethereum: "0x6a068E0287e55149a2a8396cbC99578f9Ad16A31", firstname: "dave", lastname: "saville", valid: true},
	{ethereum: "0x08511d6c42Bd247D82746c17a3EEf0Cb235f2c48", firstname: "Ben", lastname: "Georzel", valid: true},
	{ethereum: "08511d6c42Bd247D82746c17a3EE", firstname: "terrence", lastname: "phillip", valid: false},
	{ethereum: "0x08511d6c42Bd247D82746c17a3EEf0Cb235f2c48", firstname: "", lastname: "Morty", valid: false},
	{ethereum: "0x08511d6c42Bd247D82746c17a3EEf0Cb235f2c48", firstname: "Rick", lastname: "", valid: false},
}

var testEmails = []testEmail{
	{email: "test@test.com", valid: true},
	{email: "test2@test.com.au", valid: true},
	{email: "test@test.co", valid: true},
	{email: "test@test.net.au", valid: true},
	{email: "@test.com", valid: false},
	{email: "test2@testcom", valid: false},
	{email: "test@.co", valid: false},
	{email: "test.com", valid: false},
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