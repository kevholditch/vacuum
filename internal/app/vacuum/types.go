package vacuum

type Resource interface {
	ID() *string
}

type Region string

type Vacuumer interface {
	Identify(region Region) (Resources, error)
	Clean(resources Resources, cleaned func(amount int)) error
	Type() string
}

type Resources interface {
	Region() Region
	Resources() []Resource
}

type resources struct {
	resources []Resource
	region    Region
}

func (r *resources) Region() Region {
	return r.region
}

func (r *resources) Resources() []Resource {
	return r.resources
}
