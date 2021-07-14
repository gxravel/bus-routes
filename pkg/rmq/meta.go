package rmq

// Meta contains of meta information for message delivering.
type Meta struct {
	XName         string
	XType         string
	QName         string
	Key           string
	CorrID        string
	Mode          uint8
	PrefetchCount int
}

func GetMetaDetailedRoutesAccept() *Meta {
	return &Meta{
		XName: "x_detailed-routes_accept",
		XType: "direct",
		QName: "",
		Key:   "key_detailed-routes_accept",
	}
}

func GetMetaDetailedRoutesTransmit() *Meta {
	return &Meta{
		XName: "x_detailed-routes_transmit",
		XType: "direct",
		QName: "",
		Key:   "key_detailed-routes_transmit",
	}
}

func GetMetaDetailedRoutesRPC() *Meta {
	return &Meta{
		Key:           "key_detailed-routes_rpc",
		PrefetchCount: 0,
	}
}
