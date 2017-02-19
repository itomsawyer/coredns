package engine

type SrcView map[ViewKey]View

func (v SrcView) AddView(view View) {
	v[view.ViewKey] = view
}

type View struct {
	ViewKey
	RouteSetID   int
	RouteSetName string
	Upstream     *Upstream
}

type ViewKey struct {
	ClientSetID  int
	DomainPoolID int
}

func (v *View) View() (*Upstream, int) {
	return v.Upstream, v.RouteSetID
}
