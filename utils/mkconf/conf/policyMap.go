package conf

import (
	"fmt"
	"sort"

	"github.com/coredns/coredns/middleware/vane/models"
)

type PolicyMap map[string]models.PolicySlice

func (p PolicyMap) Add(v *models.PolicyView) {
	if p[v.PolicyName] == nil {
		p[v.PolicyName] = models.PolicySlice{v}
		return
	}

	p[v.PolicyName] = append(p[v.PolicyName], v)
	sort.Sort(p[v.PolicyName])
}

func (p PolicyMap) AgentKeys(policyName string) (keys []string) {
	if p[policyName] == nil {
		return nil
	}

	set := make(map[string]bool, 1)
	for _, v := range p[policyName] {
		key := fmt.Sprintf("%s-%d", v.PolicyName, v.Priority)
		if _, ok := set[key]; ok {
			continue
		}

		set[key] = true
		keys = append(keys, key)
	}

	return keys
}

func (p PolicyMap) AllAgentsKeys() (keys []string) {
	set := make(map[string]bool, 1)
	for _, policy := range p {
		for _, v := range policy {
			key := fmt.Sprintf("%s-%d", v.PolicyName, v.Priority)
			if _, ok := set[key]; ok {
				continue
			}
			set[key] = true
			keys = append(keys, key)
		}
	}

	return keys
}
