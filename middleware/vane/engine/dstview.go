package engine

type DstView map[DomainLinkKey]DomainLink

func (v DstView) AddDomainLink(dl DomainLink) {
	v[dl.DomainLinkKey] = dl
}

type DomainLink struct {
	DomainLinkKey
	NetLinkSetID   int
	NetLinkSetName string
}

type DomainLinkKey struct {
	NetLinkID    int
	DomainPoolID int
}
