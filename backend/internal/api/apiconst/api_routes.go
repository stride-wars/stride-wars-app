package apiroute

type ApiRoute string

const (
	// Auth routes
	Signup ApiRoute = "/signup"
	Signin ApiRoute = "/signin"

	// User routes
	UpdateUsername ApiRoute = "/update"

	// Activity routes
	CreateActivity ApiRoute = "/create"

	// Leaderboard routes
	GetLeaderboardByBBox ApiRoute = "/bbox"

	// Test route
	Test ApiRoute = "/test"
)

func (a ApiRoute) String() string {
	return string(a)
}
