package board

// Board represent generic board
type Board interface {
	// IsOnline permit to know if board is online
	IsOnline() bool
}
