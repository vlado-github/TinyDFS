package	consensus

// NetworkTuple contains IP and port
type NetworkTuple interface {
	GetIP() string
	GetPort() string
}

type networktuple struct {
	ipAddress string `json:"ipAddress"`
	port	string `json:"port"`
}

// NewNetworkTuple creates a new instance of network tuple
func NewNetworkTuple(ipAddress string, port string) NetworkTuple {
	return &networktuple {
		ipAddress,
		port,
	}
}

func (nt *networktuple) GetIP() string{
	return nt.ipAddress
}

func (nt *networktuple) GetPort() string{
	return nt.port
}