package nettree

import (
	"errors"
	"net"
)

var (
	errBitsRangeInvalid = errors.New("bits range is invalid")
	errDuplicated       = errors.New("Duplicated key")
)

type NetTree struct {
	node  [2]*NetTree
	value interface{}
}

func (t *NetTree) Insert(key uint32, bits int, v interface{}) error {
	return t.insert(key, bits, v, false)
}

func (t *NetTree) InsertOrUpdate(key uint32, bits int, v interface{}) error {
	return t.insert(key, bits, v, true)
}

func (t *NetTree) insert(key uint32, bits int, v interface{}, force bool) error {
	if bits < 0 || bits > 32 {
		return errBitsRangeInvalid
	}

	if bits == 0 {
		t.value = v
		return nil
	}

	cur := t
	for i := 0; i < bits; i++ {
		bit := getBit(key, i)
		if bit == 0 {
			if cur.node[0] == nil {
				cur.node[0] = new(NetTree)
			}

			cur = cur.node[0]
		} else {
			if cur.node[1] == nil {
				cur.node[1] = new(NetTree)
			}

			cur = cur.node[1]
		}
	}

	if cur.value != nil && !force {
		return errDuplicated
	}

	cur.value = v
	return nil
}

func (t *NetTree) Find(key uint32, bits int) interface{} {
	var found interface{}

	if bits < 0 || bits > 32 {
		return nil
	}

	cur := t
	if cur.value != nil {
		found = cur.value
	}

	for i := 0; i < bits; i++ {
		bit := getBit(key, i)
		if bit == 0 {
			if cur.node[0] == nil {
				break
			}
			cur = cur.node[0]
		} else {
			if cur.node[1] == nil {
				break
			}
			cur = cur.node[1]
		}

		if cur.value != nil {
			found = cur.value
		}
	}

	return found
}

func (t *NetTree) InsertByIPNet(cidr *net.IPNet, v interface{}) error {
	size, _ := cidr.Mask.Size()
	return t.Insert(ipToUint(cidr.IP.To4()), size, v)
}

func (t *NetTree) InsertOrUpdateByIPNet(cidr *net.IPNet, v interface{}) error {
	if cidr == nil {
		return errors.New("invalid cidr which is nil")
	}

	size, _ := cidr.Mask.Size()
	return t.InsertOrUpdate(ipToUint(cidr.IP.To4()), size, v)
}

func (t *NetTree) FindByIP(ip net.IP) interface{} {
	if ip == nil {
		return nil
	}

	return t.Find(ipToUint(ip.To4()), 32)
}

// From: http://stackoverflow.com/questions/2249731/how-to-get-bit-by-bit-data-from-a-integer-value-in-c

// Return bit k from n. We count from the left.
// So k = 0 is the first bit on the left and k = 31 is the last bit on the right.
func getBit(n uint32, k int) byte {
	return byte((n & (uint32(0x80000000) >> uint(k))) >> uint(32-k-1))
}

func ipToUint(ip net.IP) uint32 {
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}
