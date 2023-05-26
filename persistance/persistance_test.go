package persistance

import (
	"fmt"
	"testing"
)

func TestCreate(t *testing.T) {
	bp := BadgerPersistance{inMemory: false}

	badger := bp.Create("db", "kacpietrzak@gmail.com")
	cards := badger.ReadAll()

	pp := PostgresPersistance{}
	postgres := pp.Create("db", "kacpietrzak@gmail.com")

	for _, card := range cards {
		fmt.Printf("%v\n", card)
		postgres.Add(&card)
	}
}
