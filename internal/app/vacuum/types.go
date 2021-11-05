package vacuum

type Resource interface {
	ID() *string
}

type Region string

type Vacuumer interface {
	Identify(region Region) (Resources, error)
	Clean(resources Resources) error
}

type Resources interface {
	Region() Region
	Resources() []Resource
}
