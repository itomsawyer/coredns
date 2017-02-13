package vane

type Loader interface {
	Init(nsd string) error
	LoadAll() (*Engine, error)
	//ReloadClientSet(engine Enginer) (int, error)
}

var (
	DBLoaders = map[string]Loader{
		"default": new(MySQLoader),
	}
)

func GetLoader(l string) Loader {
	return DBLoaders[l]
}
