package product

// sample ProductIds
var (
	P0001 ProductId = "P0001"
	P0002 ProductId = "P0002"
	P0003 ProductId = "P0003"
	P0004 ProductId = "P0004"
	P0005 ProductId = "P0005"
)

// sample products
var (
	Lightsaber = &Product{
		P0001,
		"Lightsaber",
		"The perfect lightsaber for an aspiring Jedi",
		"Weapons",
		"http://images.buystarwarstoys.com/products/9288/1-1/ahsoka-tano-toy-lightsaber.jpg",
		999.99,
		10,
	}

	MilleniumFalcon = &Product{
		P0002,
		"The Millenium Falcon",
		"The fastest ship in the entire gallaxy - finished the Kessel Run in less than 12 parsecs",
		"Mobility",
		"http://ksassets.timeincuk.net/wp/uploads/sites/54/2017/11/Millenium-Falcon.jpg",
		30000.00,
		1,
	}
)