package consensus

import "errors"

// Command is a message used in Raft algorithm to
// communicate between nodes.
type Command interface {
	SetAck()
	SetNAck()
	IsAck() bool
	GetType() CommandType
	Validate() error
}

type command struct {
	termId  int
	cmdType CommandType
	values  map[string]string
	isAck   bool
}

// NewCommand creates new instance of command
func NewCommand(termId int, cmdType CommandType, values map[string]string) Command {
	return &command{
		termId:  termId,
		cmdType: cmdType,
		values:  values,
		isAck:   false,
	}
}

func (cmd *command) SetAck() {
	cmd.isAck = true
}

func (cmd *command) SetNAck() {
	cmd.isAck = false
}

func (cmd *command) IsAck() bool {
	return cmd.isAck
}

func (cmd *command) GetType() CommandType {
	return cmd.cmdType
}

func (cmd *command) Validate() error {
	switch cmd.cmdType {
	case REQUESTVOTE:
		{
			if cmd.termId <= 0 {
				return errors.New("command's termId must be larger than 0")
			}
			return nil
		}
	case APPENDENTRY:
		{
			if cmd.termId <= 0 {
				return errors.New("command's termId must be larger than 0")
			}
			if len(cmd.values) <= 0 {
				return errors.New("command values can not be empty")
			}
			return nil
		}
	case HEARTBEAT:
		{
			if cmd.termId <= 0 {
				return errors.New("command's termId must be larger than 0")
			}
			return nil
		}
	default:
		{
			if cmd.termId <= 0 {
				return errors.New("command's termId must be larger than 0")
			}
			return nil
		}
	}
}
