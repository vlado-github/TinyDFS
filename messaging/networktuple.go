package messaging

// NetworkTuple contains IP and port
type NetworkTuple interface {
	GetIP() string
	GetPort() string
	GetId() string
	GetQueuePort() string
	GetAvailableStatus() bool
	SetIsAvailable(bool)
}

type networktuple struct {
	IpAddress   string `json:"IpAddress"`
	Port        string `json:"Port"`
	Id          string `json:"Id"`
	QueuePort   string `json:QueuePort`
	IsAvailable bool   `json:IsAvailable`
}

// NewNetworkTuple creates a new instance of network tuple
func NewNetworkTuple(id string, ipAddress string, port string, queuePort string) NetworkTuple {
	return &networktuple{
		IpAddress:   ipAddress,
		Port:        port,
		Id:          id,
		QueuePort:   queuePort,
		IsAvailable: true,
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

func (nt *networktuple) GetQueuePort() string {
	return nt.QueuePort
}

func (nt *networktuple) GetAvailableStatus() bool {
	return nt.IsAvailable
}

func (nt *networktuple) SetIsAvailable(isAvailable bool) {
	nt.IsAvailable = isAvailable
}
