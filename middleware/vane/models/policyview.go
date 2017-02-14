package models

import (
	"errors"
	"reflect"
	"sort"
	"strings"

	"github.com/astaxie/beego/orm"
)

type PolicyView struct {
	PolicyId       int `orm:"pk"`
	PolicyName     string
	PolicySequence int
	Priority       int
	Op             string `orm:"size(128)"`
	OpTyp          string `orm:"size(128)"`
	LdnsId         int
	Name           string `orm:"size(128)"`
	Addr           string `orm:"size(128)"`
	Typ            string `orm:"size(128)"`
	RrsetId        int
}

func init() {
	orm.RegisterModel(new(PolicyView))
}

type PolicySlice []*PolicyView

func (p PolicySlice) Len() int {
	return len(p)
}

func (p PolicySlice) Less(i, j int) bool {
	return p[i].Priority < p[j].Priority
}

func (p PolicySlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p PolicySlice) Hosts() [][]string {
	if len(p) == 0 {
		return nil
	}

	sort.Sort(p)

	prio := -1
	hosts := make([][]string, 0, 4)
	hs := make([]string, 0, 1)

	for _, p := range p {
		if p.Priority > prio {
			if hs != nil {
				hosts = append(hosts, hs)
			}

			hs = make([]string, 0, 1)
			prio = p.Priority
		}

		hs = append(hs, p.Addr)
	}

	if len(hs) > 0 {
		hosts = append(hosts, hs)
	}

	return hosts
}

type PolicySet map[int]PolicySlice

func (p PolicySet) Add(id int, pv *PolicyView) {
	if p[id] == nil {
		p[id] = PolicySlice{pv}
		return
	}

	p[id] = append(p[id], pv)
}

// GetAllPolicyView retrieves all PolicyView matches certain condition. Returns empty list if
// no records exist
func GetAllPolicyView(o orm.Ormer, query Values, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	if o == nil {
		o = orm.NewOrm()
	}
	qs := o.QueryTable(new(PolicyView))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v...)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []PolicyView
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}
