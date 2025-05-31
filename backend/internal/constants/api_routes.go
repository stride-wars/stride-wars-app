package constants

type ApiRoute string

const (
	// Auth routes
	Signup ApiRoute = "/signup"
	Signin ApiRoute = "/signin"

	// User routes
	GetUserByUsername ApiRoute = "/by-username"
	GetUserByID       ApiRoute = "/by-id"
	UpdateUsername    ApiRoute = "/update"

	// Activity routes
	CreateActivity ApiRoute = "/create"

	// Leaderboard routes
	GetLeaderboardByBBox ApiRoute = "/by-bbox"

	// Test route
	Test ApiRoute = "/test"
)

func (a ApiRoute) String() string {
	return string(a)
}
