package db

type Position struct {
	Status  uint8
	Count   uint8
	Id      uint16
	Cost    float32
	Product string
}

type Receipt struct {
	Status    uint8
	Id        uint16
	Total     uint16
	Pos_len   int
	Positions []Position
	Date      string
}

func NewPosition() Position {
	return Position{0, 0, 0, 0, ""}
}

func NewReceipt() Receipt {
	pos := make([]Position, 0)
	pos = append(pos, NewPosition())
	pos = append(pos, NewPosition())
	pos_len := len(pos) + 1
	return Receipt{0, 0, 0, pos_len, pos, "1999-01-01"}
}
