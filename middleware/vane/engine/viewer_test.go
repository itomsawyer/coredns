package engine

import (
	"testing"
)

func TestViewAdd(t *testing.T) {
	viewer := SrcView{}
	viewer.AddView(View{
		ViewKey: ViewKey{
			ClientSetID:  1,
			DomainPoolID: 1,
		},
		RouteSetID: 1,
	})
	viewer.AddView(View{
		ViewKey: ViewKey{
			ClientSetID:  2,
			DomainPoolID: 2,
		},
		RouteSetID: 2,
	})
	viewer.AddView(View{
		ViewKey: ViewKey{
			ClientSetID:  2,
			DomainPoolID: 3,
		},
		RouteSetID: 3,
	})
	viewer.AddView(View{
		ViewKey: ViewKey{
			ClientSetID:  1,
			DomainPoolID: 1,
		},
		RouteSetID: 4,
	})

	t.Log(viewer)
}
