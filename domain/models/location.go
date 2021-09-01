package models

type Location struct {
	// TODO: do I need separate IDs for database and in-game use?
	//ID:
	Name        string
	Restriction *Restriction
}

type Restriction struct {
	Level int
}
