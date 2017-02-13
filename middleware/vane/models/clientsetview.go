package models

import (
	"errors"
	"net"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type ClientSetView struct {
	IpnetId       int64  `orm:"pk"`
	IpStart       string `orm:"size(128)"`
	IpEnd         string `orm:"size(128)"`
	Ipnet         string `orm:"size(128)"`
	Mask          uint8
	ClientSetId   int64  `orm:"column(clientset_id)"`
	ClientSetName string `orm:"column(clientset_name);size(128)"`
}

func (m *ClientSetView) TableName() string {
	return "clientset_view"
}

func init() {
	orm.RegisterModel(new(ClientSetView))
}

// GetClientSetViewById retrieves IpsetView by Id. Returns error if
// Id doesn't exist
func GetClientSetView(ip *net.IP) (v []*ClientSetView, err error) {
	o := orm.NewOrm()
	sql := "select * from ipset_view where inet_aton(ip_start) <= , inet_aton(ip_end) >="
	_, err = o.Raw(sql).QueryRows(&v)
	return
}

// GetAllClientSetView retrieves all IpsetView matches certain condition. Returns empty list if
// no records exist
func GetAllClientSetView(query Values, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(ClientSetView))
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

	var l []ClientSetView
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
