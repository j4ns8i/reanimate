// package horde provides methods to reanimate and manage hordes of workers.
package horde

type Reanimator interface {
	// Reanimate a new horde (i.e. create and attach)
	Reanimate() (Horde, error)
	// List existing hordes
	List() ([]Horde, error)
}

type Horde interface {
	// Summon forth the horde (i.e. attach)
	Summon() error
	// Return the name of the horde
	Name() string
	// Destroy the horde (i.e. detach and clean up)
	Destroy() error
}
