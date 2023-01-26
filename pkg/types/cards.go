package types

import (
	"fmt"
)

type Species uint8

const (
	Dog Species = iota + 1
	Cat
	Bird
)

func (s Species) String() string {
	species := []string{"dog", "cat", "bird"}
	if s < Dog || s > Bird {
		panic(fmt.Sprintf("Unexpected species: %d", s))
	}
	return species[s-1]
}

type EventType int

const (
	Found EventType = iota + 1
	Lost
)

func (c EventType) String() string {
	cardTypes := []string{"found", "lost"}
	if c < Found || c > Lost {
		panic(fmt.Sprintf("Unexpected card type: %d", c))
	}
	return cardTypes[c-1]
}

type Sex int

const (
	UndefinedSex Sex = iota
	Male
	Female
)

func (s Sex) String() string {
	sexes := []string{"unknown", "male", "female"}
	if s < UndefinedSex || s > Female {
		panic(fmt.Sprintf("Unexpected animal sex: %d", s))
	}
	return sexes[s]
}
