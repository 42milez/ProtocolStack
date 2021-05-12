package e

type Error int

const(
	OK Error = iota
	AlreadyOpened
	CantOpen
)
