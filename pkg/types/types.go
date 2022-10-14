package types

type CardID int

type LatestKnownCardsSource interface {
	getLatestKnownCards() []CardID
}
