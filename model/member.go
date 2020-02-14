package model

import (
	"blockchain.automation/fsm"
)

const (
	Closed   = "Closed"
	Register = "Register"
	Transfer = "Transfer"
	Charge   = "Charge"
	End      = "End"
)

type Member struct {
	MemberId      string
	VsCode        string
	CountryCode   string
	CurrencyCode  string
	MemberRole    string
	WalletAddress string
	CreateDate    string

	Balance float64
	Cookies []string

	events        fsm.Events
	automationFsm *fsm.FSM
}

func NewMember() *Member {
	return new(Member)
}

func (m *Member) Init() {
	m.events = fsm.Events{
		{Name: Register, Src: []string{Closed}, Dst: Register},
		{Name: Transfer, Src: []string{Register}, Dst: Transfer},
		{Name: Charge, Src: []string{Transfer}, Dst: Charge},
		{Name: Transfer, Src: []string{Charge}, Dst: Transfer},
		{Name: Transfer, Src: []string{Transfer}, Dst: Transfer},
		{Name: End, Src: []string{Transfer}, Dst: End},
	}

	m.automationFsm = fsm.NewFSM(
		Closed,
		m.events,
		fsm.Callbacks{},
	)

	logger.Info("Member Init: ", m.MemberId, ": ", m.automationFsm.Current())
	messageQueue.SendMessage(m)

	// go func() {
	//	ch <- true
	// }()
}

func (m *Member) State() string {
	return m.automationFsm.Current()
}

func (m *Member) SetState(state string) error {
	err := m.automationFsm.Event(state)
	return err
}
