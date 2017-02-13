package vane

type DBer interface {
	Open(nsd string) error
	Load() error
	GetClientSetID(ip string) (int, error)
	GetDomainID(domain string) (int, error)
	GetPolicy(client int, domain int) (policy int, router int, err error)
	GetOutLink(router int, domain int, netlink int) error
}
