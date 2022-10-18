package types

import (
	"fmt"
)

type Species uint8

const (
	Dog Species = iota + 1
	Cat
)

func (s Species) String() string {
	species := []string{"dog", "cat"}
	if s < Dog || s > Cat {
		panic(fmt.Sprintf("Unexpected species: %d", s))
	}
	return species[s-1]
}

type CardType int

const (
	Found CardType = iota + 1
	Lost
)

func (c CardType) String() string {
	cardTypes := []string{"found", "lost"}
	if c < Found || c > Lost {
		panic(fmt.Sprintf("Unexpected card type: %d", c))
	}
	return cardTypes[c-1]
}

type Sex int

const (
	Male Sex = iota + 1
	Female
)

func (s Sex) String() string {
	sexes := []string{"unknown", "male", "female"}
	if s < Male || s > Female {
		panic(fmt.Sprintf("Unexpected animal sex: %d", s))
	}
	return sexes[s-1]
}
