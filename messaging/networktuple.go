package messaging

// NetworkTuple contains IP and port
type NetworkTuple interface {
	GetIP() string
	GetPort() string
	GetId() string
}

type networktuple struct {
	IpAddress string `json:"IpAddress"`
	Port      string `json:"Port"`
	Id        string `json:"Id"`
}

// NewNetworkTuple creates a new instance of network tuple
func NewNetworkTuple(id string, ipAddress string, port string) NetworkTuple {
	return &networktuple{
		IpAddress: ipAddress,
		Port:      port,
		Id:        id,
	}
}

func (nt *networktuple) GetIP() string {
	return nt.IpAddress
}

func (nt *networktuple) GetPort() string {
	return nt.Port
}

func (nt *networktuple) GetId() string {
	return nt.Id
}
