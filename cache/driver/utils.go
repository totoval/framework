package driver

import "github.com/totoval/framework/helpers/zone"

func durationFromNow(future zone.Time) zone.Duration {
	return future.Sub(zone.Now())
}

type key struct {
	raw    string
	prefix string
}

func newKey(raw string, prefix string) *key {
	k := key{}
	k.prefix = prefix
	k.raw = raw
	return &k
}
func (k *key) Raw() string {
	return k.raw
}
func (k *key) Prefixed() string {
	return k.prefix + k.raw
}
