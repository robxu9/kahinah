package kahinah

import "errors"

var (
	// ErrConnectorName - connector name is invalid or already in use
	ErrConnectorName = errors.New("kahinah: connector name is invalid or is already attached")
	// ErrConnectorInit - the Init() call failed
	ErrConnectorInit = errors.New("kahinah: unable to initialise connector")
)

// Connector represents an external hook in to the Kahinah advisory system.
//
// Typically, clients will initialise a Kahinah object, then add appropriate
// connectors via the Attach() function. On attaching, the connector
// will be called with the Init() function, passing in the Kahinah
// object. Then, any updates that pass or fail comment are passed into the
// connector as necessary (may be concurrent - be careful!). The connector
// can also create new updates by calling the NewUpdate() method on the
// Kahinah object.
type Connector interface {
	// Name returns the name of the connector, usually in reverse domain name notation.
	Name() string
	Init(*Kahinah) error
	Pass(*Update)
	Fail(*Update)
}

// Attach attaches a connector to Kahinah, which can be used to
// add updates to the system and pass/fail updates.
func (k *Kahinah) Attach(c Connector) error {
	if _, ok := k.connectors[c.Name()]; ok {
		return ErrConnectorName
	}

	if c.Init(k) != nil {
		return ErrConnectorInit
	}

	k.connectors[c.Name()] = c

	return nil
}
