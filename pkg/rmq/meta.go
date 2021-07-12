package rmq

// Meta contains of meta information for message delivering.
type Meta struct {
	XName  string
	XType  string
	QName  string
	Key    string
	CorrID string
	Mode   uint8
}
