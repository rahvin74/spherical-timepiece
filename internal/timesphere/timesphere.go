package timesphere

// MinuteBall represents a steel ball as it moves through the clock mechanism.
// It keeps track of it's original position.
type MinuteBall struct {
	originalPosition int
}

// GetOriginalPosition returns the original position property of the MinuteBall
func (m *MinuteBall) GetOriginalPosition() int {
	return m.originalPosition
}

// SetOriginalPosition sets the original position property on the MinuteBall
func (m *MinuteBall) SetOriginalPosition(position int) {
	m.originalPosition = position
}
