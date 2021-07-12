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

var (
	MetaDetailedRoutesAccept = &Meta{
		XName: "x_detailed-routes_accept",
		XType: "direct",
		QName: "",
		Key:   "key_detailed-routes_accept",
	}

	MetaDetailedRoutesTransmit = &Meta{
		XName: "x_detailed-routes_transmit",
		XType: "direct",
		QName: "",
		Key:   "key_detailed-routes_transmit",
	}

	MetaDetailedRoutesRPC = &Meta{
		QName: "q_detailed-routes_rpc",
	}
)
