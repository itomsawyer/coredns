package engine

import (
	"github.com/coredns/coredns/middleware/pkg/dnsutil"
	"github.com/coredns/coredns/middleware/vane/models"
	"github.com/itomsawyer/llog"

	"github.com/astaxie/beego/orm"
)

type EngineBuilder struct {
	DBName          string
	Logger          *llog.Logger
	ClientSetView   []models.ClientSetView
	ClientSetWLView []models.ClientSetWLView
	NetLinkView     []models.NetLinkView
	NetLinkWLView   []models.NetLinkWLView
	DomainView      []models.DomainView
	SrcView         []models.SrcView
	DstView         []models.DstView
	PolicyView      []models.PolicyView
	RouteView       []models.RouteView
}

func (b *EngineBuilder) Load() (err error) {
	o := orm.NewOrm()
	if b.DBName != "" {
		o.Using(b.DBName)
	}

	err = o.Begin()
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			o.Rollback()
		} else {
			o.Commit()
		}
	}()

	b.ClientSetView, err = models.GetClientSetView(o, nil, nil, nil, 0, -1)
	if err != nil {
		return err
	}

	b.ClientSetWLView, err = models.GetClientSetWLView(o, nil, nil, nil, 0, -1)
	if err != nil {
		return err
	}

	b.NetLinkView, err = models.GetNetLinkView(o, nil, nil, nil, 0, -1)
	if err != nil {
		return err
	}

	b.NetLinkWLView, err = models.GetNetLinkWLView(o, nil, nil, nil, 0, -1)
	if err != nil {
		return err
	}

	b.DomainView, err = models.GetDomainView(o, nil, nil, nil, 0, -1)
	if err != nil {
		return err
	}

	b.SrcView, err = models.GetSrcView(o, nil, nil, nil, 0, -1)
	if err != nil {
		return err
	}

	b.DstView, err = models.GetDstView(o, nil, nil, nil, 0, -1)
	if err != nil {
		return err
	}

	b.PolicyView, err = models.GetPolicyView(o, nil, nil, nil, 0, -1)
	if err != nil {
		return err
	}

	b.RouteView, err = models.GetRouteView(o, nil, nil, nil, 0, -1)
	if err != nil {
		return err
	}

	return nil
}

func (b *EngineBuilder) Build(e *Engine) (err error) {
	err = b.BuildClientSet(e)
	if err != nil {
		return
	}

	err = b.BuildNetLink(e)
	if err != nil {
		return
	}

	err = b.BuildDomainView(e)
	if err != nil {
		return
	}

	err = b.BuildUpstream(e)
	if err != nil {
		return
	}

	err = b.BuildSrcView(e)
	if err != nil {
		return
	}

	err = b.BuildDstView(e)
	if err != nil {
		return
	}

	err = b.BuildRoute(e)
	if err != nil {
		return
	}

	return
}

func (b *EngineBuilder) BuildClientSet(e *Engine) error {
	for _, v := range b.ClientSetView {
		ipnet, err := dnsutil.ParseIPNet(v.Ipnet, int(v.Mask))
		if err != nil {
			b.Logger.Error("ParseIPNet for clientset %v failed, %s", v, err)
			return err
		}

		err = e.AddClient(ipnet, v.ClientSetId, v.ClientSetName, 0)
		if err != nil {
			b.Logger.Error("AddClientSet %v failed, %s", v, err)
			return err
		}
	}

	for _, v := range b.ClientSetWLView {
		ipnet, err := dnsutil.ParseIPNet(v.Ipnet, int(v.Mask))
		if err != nil {
			b.Logger.Error("ParseIPNet for clientsetwl %v failed, %s", v, err)
			return err
		}

		err = e.AddClientWL(ipnet, v.ClientSetId, v.ClientSetName, 1)
		if err != nil {
			b.Logger.Error("AddClientSetWL %v failed, %s", v, err)
			return err
		}
	}

	return nil
}

func (b *EngineBuilder) BuildNetLink(e *Engine) error {
	for _, v := range b.NetLinkView {
		ipnet, err := dnsutil.ParseIPNet(v.Ipnet, int(v.Mask))
		if err != nil {
			b.Logger.Error("ParseIPNet for netlink %v failed, %s", v, err)
			return err
		}

		err = e.AddNetLink(ipnet, v.NetLinkId, v.Isp, v.Region, 0)
		if err != nil {
			b.Logger.Error("AddNetLink %v failed, %s", v, err)
			return err
		}
	}

	for _, v := range b.NetLinkWLView {
		ipnet, err := dnsutil.ParseIPNet(v.Ipnet, int(v.Mask))
		if err != nil {
			b.Logger.Error("ParseIPNet for netlink %v failed, %s", v, err)
			return err
		}

		err = e.AddNetLinkWL(ipnet, v.NetLinkId, v.Isp, v.Region, 1)
		if err != nil {
			b.Logger.Error("AddNetLinkWL %v failed, %s", v, err)
			return err
		}
	}

	return nil
}

func (b *EngineBuilder) BuildDomainView(e *Engine) error {
	for _, v := range b.DomainView {
		d := Domain{
			ID:       v.DomainId,
			Domain:   v.Domain,
			DmPoolID: v.DomainPoolId,
			DmPool:   v.PoolName,
			Monitor:  v.DomainMonitor,
		}

		err := e.AddDomain(d)
		if err != nil {
			b.Logger.Error("AddDomain %v failed, %s", v, err)
			return err
		}
	}

	return nil
}

func (b *EngineBuilder) BuildUpstream(e *Engine) error {
	for _, v := range b.PolicyView {
		upstream := e.AddUpstream(v.PolicyId, v.PolicyName)

		uh, err := e.AddUpstreamHost(v.Addr, false)
		if err != nil {
			b.Logger.Error("AddUpstreamHost %v failed, %s", v, err)
			return err
		}

		upstream.AddHost(uh, v.Priority)
	}

	return nil
}

func (b *EngineBuilder) BuildSrcView(e *Engine) error {
	for _, v := range b.SrcView {
		up, err := e.GetUpstreamByID(v.PolicyId)
		if err != nil {
			b.Logger.Warn("AddSrcView %v with no policy found", v)
			continue
		}

		view := View{
			ViewKey: ViewKey{
				ClientSetID:  v.ClientSetId,
				DomainPoolID: v.DomainPoolId,
			},
			RouteSetID:   v.RouteSetId,
			RouteSetName: v.RouteSetName,
			Upstream:     up,
		}

		e.AddView(view)
	}

	return nil
}

func (b *EngineBuilder) BuildDstView(e *Engine) error {
	for _, v := range b.DstView {
		dl := DomainLink{
			DomainLinkKey: DomainLinkKey{
				DomainPoolID: v.DomainPoolId,
				NetLinkID:    v.NetLinkId,
			},
			NetLinkSetID: v.DstViewId,
		}

		e.AddDomainLink(dl)
	}

	return nil
}

func (b *EngineBuilder) BuildRoute(e *Engine) error {
	for _, v := range b.RouteView {
		ot := OutLink{v.OutlinkName, v.OutlinkAddr}

		r := Route{
			RouteKey: RouteKey{
				RouteSetID:   v.RoutesetId,
				NetLinkSetID: v.NetlinksetId,
			},

			OutLink:  ot,
			Priority: v.RoutePriority,
			Score:    v.RouteScore,
		}

		e.AddRoute(r)
	}

	return nil
}
