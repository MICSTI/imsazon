package user

// Sample UserIds
var (
	U0001 UserId = "U0001"
	U0002 UserId = "U0002"
	U0003 UserId = "U0003"
)

// Sample users
var (
	Rey = &User{U0001, "Rey", "rey@jedi.com", "rey", "rey123", Standard}
	Kylo = &User{U0002, "Kylo", "kylo@firstorder.com", "kylo", "kylo123", Standard}
	Luke = &User{ U0003, "Luke", "luke@jedi.com", "luke", "luke123", Admin}
)