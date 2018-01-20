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
		"The perfect lightsaber for every aspiring Jedi",
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

	BB8 = &Product{
		P0003,
		"BB 8",
		"Extraordinarily helpful droid",
		"Droids",
		"https://images.fun.com/products/34909/2-1-63328/star-wars-episode-7-rey-jakku-and-bb8-black-series-set.jpg",
		12499,
		3,
	}

	Podracer = &Product{
		P0004,
		"Podracer",
		"Lightning-fast podracer - nobody will be able to beat you",
		"Mobility",
		"https://images-na.ssl-images-amazon.com/images/I/41j3vMHSX0L._AA300_.jpg",
		3499.00,
		6,
	}

	CarboniteFreezer = &Product{
		P0005,
		"Carbonite Freezer",
		"Very useful in case you need to freeze someone in carbonite",
		"Utilities",
		"https://s-i.huffpost.com/gen/1359887/images/o-HAN-SOLO-CARBONITE-facebook.jpg",
		39999.99,
		2,
	}
)