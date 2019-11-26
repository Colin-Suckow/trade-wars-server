package tradewars

//TODO: Finish elemental stuff later. Currently out of scope

type Element struct {
	id            int
	description   string
	strongAgaisnt int
	weakAgainst   int
}

/*
	Balistic    Element = Element{0, "Basic physical", 4, 1}
	Plasma      Element = Element{1, "Green magic stuff", 0, 2}
	Electricity Element = Element{2, "Yellow shocking stuff", 1, 3}
	SpaceWater  Element = Element{3, "Its water, but in space", 2, 4}
	SpaceRock   Element = Element{4, "Wow we need better names", 3, 0}
*/
