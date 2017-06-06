package conf

import (
	"fmt"
	"sort"
	"strings"

	"github.com/coredns/coredns/middleware/vane/models"

	"github.com/astaxie/beego/orm"
)

type Conf struct {
	Agents            map[string]*Agent `json:"agents"`
	DomainPools       SliceMap          `json:"domainpools"`
	Views             SliceMap          `json:"views"`
	ClientSets        SliceMap          `json:"clientsets"`
	ldns              map[string]*DNS
	ldnsMirror        map[string]string
	corednsDomainPool map[string]bool
	corednsHost       []string
	corednsCheckdm    []string
}

func NewConf() *Conf {
	return &Conf{
		Agents:            make(map[string]*Agent, 1),
		DomainPools:       SliceMap{},
		Views:             SliceMap{},
		ClientSets:        SliceMap{},
		ldns:              make(map[string]*DNS, 1),
		ldnsMirror:        make(map[string]string, 1),
		corednsDomainPool: make(map[string]bool, 1),
		corednsHost:       make([]string, 0, 1),
	}
}

func (p *Conf) AddCorednsDomainPool(dp string) {
	if p.corednsDomainPool == nil {
		p.corednsDomainPool = make(map[string]bool, 1)
	}

	p.corednsDomainPool[dp] = true
}

func (p *Conf) AddCorednsHost(host string, checkdm string) {
	if p.corednsHost == nil {
		p.corednsHost = make([]string, 0, 1)
	}

	p.corednsHost = append(p.corednsHost, host)

	if p.corednsCheckdm == nil {
		p.corednsCheckdm = make([]string, 0, 1)
	}

	p.corednsCheckdm = append(p.corednsCheckdm, checkdm)
}

func (p *Conf) AddIpNet(view string, ipnet string) {
	view = strings.Replace(view, "/", "-", -1)
	p.ClientSets.Add(view, ipnet)
}

func (p *Conf) AddDomain(dp string, domain string) {
	if len(domain) <= 1 {
		return
	}

	if domain[0] == '.' {
		domain = domain[1:]
	}

	if p.DomainPools.HasValue(dp, domain) {
		return
	}

	p.DomainPools.Add(dp, domain)
}

func (p *Conf) AddView(view string, agentName string) {
	view = strings.Replace(view, "/", "-", -1)
	if p.Views.HasValue(view, agentName) {
		return
	}

	p.Views.Add(view, agentName)
}

func (p *Conf) CreateAgent(key string, common bool, dp string) (*Agent, error) {
	var agentKey string
	var corednsAgent bool

	if p.Agents == nil {
		p.Agents = make(map[string]*Agent, 1)
	}

	if len(dp) != 0 {
		if _, ok := p.DomainPools[dp]; !ok {
			return nil, fmt.Errorf("DNS agent %s expect domain pool %s but not exists", key, dp)
		}
	}

	if len(dp) != 0 {
		if _, ok := p.corednsDomainPool[dp]; ok {
			agentKey = fmt.Sprintf("%s-coredns", dp)
			corednsAgent = true
		} else {
			agentKey = fmt.Sprintf("%s-%s", dp, key)
		}
	} else {
		agentKey = fmt.Sprintf("default-%s", key)
	}

	if ag, ok := p.Agents[agentKey]; ok {
		return ag, nil
	}

	if corednsAgent {
		p.Agents[agentKey] = &Agent{
			Name:        agentKey,
			Common:      common,
			DomainPool:  dp,
			DNS:         p.corednsHost,
			Ecs:         true,
			CheckDomain: p.corednsCheckdm,
		}

		if mirror, ok := p.ldnsMirror["__coredns__"]; ok {
			p.Agents[agentKey].Mirror = mirror
		} else {
			p.ldnsMirror["__coredns__"] = agentKey
		}

		return p.Agents[agentKey], nil
	}

	ldnsKey := key
	if ldns, ok := p.ldns[ldnsKey]; ok {
		p.Agents[agentKey] = &Agent{
			Name:        agentKey,
			Common:      common,
			DomainPool:  dp,
			DNS:         ldns.Host,
			CheckDomain: ldns.CheckDomain,
			ExForwarder: ldns.ExForwarder,
		}

		if mirror, ok := p.ldnsMirror[ldnsKey]; ok {
			p.Agents[agentKey].Mirror = mirror
		} else {
			p.ldnsMirror[ldnsKey] = agentKey
		}

		return p.Agents[agentKey], nil
	} else {
		return nil, fmt.Errorf("ldns slice cannot be found")
	}
}

func (p *Conf) AddDNS(policy string, priority int, dns string, typ string, checkdm string) error {
	key := fmt.Sprintf("%s-%d", policy, priority)
	if p.ldns == nil {
		p.ldns = make(map[string]*DNS, 1)
	}

	if _, ok := p.ldns[key]; !ok {
		p.ldns[key] = new(DNS)
	}

	switch typ {
	case "ldns":
		p.ldns[key].AddExForwarder(dns)
	case "upstream":
		p.ldns[key].AddHost(dns, checkdm)
	default:
		return fmt.Errorf("unknown ldns type")
	}

	return nil
}

type Agent struct {
	Name        string   `json:"-"`
	DNS         []string `json:"dns"`
	Common      bool     `json:"common"`
	DomainPool  string   `json:"zones,omitempty"`
	Ecs         bool     `json:"ecs,omitempty"`
	ExForwarder []string `json:"exforwarder,omitempty"`
	CheckDomain []string `json:"checkdm,omitempty"`
	Mirror      string   `json:"mirror,omitempty"`
}

type SliceMap map[string][]string

func (v SliceMap) Get(key string) string {
	if v == nil {
		return ""
	}
	vs := v[key]
	if len(vs) == 0 {
		return ""
	}
	return vs[0]
}

func (v SliceMap) Set(key, value string) {
	v[key] = []string{value}
}

func (v SliceMap) HasValue(key, value string) bool {
	if v[key] == nil {
		return false
	}

	for _, c := range v[key] {
		if c == value {
			return true
		}
	}

	return false
}

func (v SliceMap) Add(key, value string) {
	v[key] = append(v[key], value)
}

func (v SliceMap) Del(key string) {
	delete(v, key)
}

func (conf *Conf) BuildConfFromDB(db string) error {
	var query models.Values
	if err := models.InitDB(db); err != nil {
		return err
	}

	o := orm.NewOrm()

	domains, err := models.GetDomainView(o, nil, nil, nil, 0, -1)
	if err != nil {
		return err
	}

	for _, dm := range domains {
		conf.AddDomain(dm.PoolName, dm.Domain)
	}

	cs_wl, err := models.GetClientSetWLView(o, nil, nil, nil, 0, -1)
	if err != nil {
		return err
	}

	query = models.Values{}
	query.Set("typ", "coredns")
	coredns, err := models.GetLDNS(o, query, nil, nil, 0, -1)
	if err != nil {
		return err
	}

	for _, cd := range coredns {
		conf.AddCorednsHost(cd.Addr, cd.Checkdm)
	}

	for _, c := range cs_wl {
		conf.AddIpNet(c.ClientSetName, fmt.Sprintf("%s/%d", c.Ipnet, c.Mask))
	}

	cs, err := models.GetClientSetView(o, nil, nil, nil, 0, -1)
	if err != nil {
		return err
	}

	for _, c := range cs {
		conf.AddIpNet(c.ClientSetName, fmt.Sprintf("%s/%d", c.Ipnet, c.Mask))
	}

	conf.AddIpNet("default", "any")

	policys, err := models.GetPolicyView(o, nil, []string{"policy_id", "priority"}, []string{"asc", "asc"}, 0, -1)
	if err != nil {
		return err
	}

	pm := PolicyMap{}
	for _, policy := range policys {
		if err := conf.AddDNS(policy.PolicyName, policy.Priority, policy.Addr, policy.Typ, policy.Checkdm); err != nil {
			return err
		}

		p := policy
		pm.Add(&p)
	}

	//query := models.Values{}
	//query.Set("domain_pool_id.in", []int{1, 2, 3, 4, 5, 6, 7, 8})
	//srcViews, err := models.GetSrcView(o, query, nil, nil, 0, -1)

	srcViews, err := models.GetSrcView(o, nil, []string{"domain_pool_id"}, []string{"asc"}, 0, -1)
	if err != nil {
		return err
	}

	for _, sv := range srcViews {
		domainPool := ""
		common := true
		if sv.DomainPoolId != 1 {
			common = false
			domainPool = sv.DomainPoolName
		}

		//if sv.ClientSetId == 1 {
		//	continue
		//}

		keys := pm.AgentKeys(sv.PolicyName)

		for _, k := range keys {
			if ag, err := conf.CreateAgent(k, common, domainPool); err != nil {
				return err
			} else {
				conf.AddView(sv.ClientSetName, ag.Name)
			}
		}
	}

	for _, sv := range srcViews {
		domainPool := ""
		common := true
		if sv.DomainPoolId != 1 {
			common = false
			domainPool = sv.DomainPoolName
		}

		if sv.ClientSetId != 1 {
			continue
		}

		keys := pm.AgentKeys(sv.PolicyName)

		for _, k := range keys {
			if ag, err := conf.CreateAgent(k, common, domainPool); err != nil {
				return err
			} else {
				for _, tmp := range srcViews {
					conf.AddView(tmp.ClientSetName, ag.Name)
				}
			}
		}
	}

	return nil
}
