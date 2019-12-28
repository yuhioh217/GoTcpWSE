package structure

// PQueue to queue the packet's data
type PQueue struct {
	Queue string      // Will use timestamp to define
	Pack  interface{} // Maybe it is struct type
}

// NewPQueue to new a PQueue struct object
func NewPQueue(timestamp string) *PQueue {
	return &PQueue{
		Queue: timestamp,
		Pack:  &Packets{}}
}
