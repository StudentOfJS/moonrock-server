package models

// Response models the data to send back to the handler
type Response struct {
	ServerCode int
	Response   string
}

func getResponse(r string) *Response {
	errorResponseMap := map[string]int{
		"updated":                         200,
		"account successfully confirmed":  200,
		"check your email":                200,
		"ok":                              200,
		"unauthenticated":                 401,
		"server error":                    500,
		"invalid login":                   400,
		"invalid email":                   400,
		"invalid signup":                  400,
		"invalid address":                 400,
		"user doesn't exist":              400,
		"token expired, please try again": 400,
		"update failed":                   400,
		"already signed up":               400,
	}

	response := Response{ServerCode: errorResponseMap[r], Response: r}
	return &response
}
