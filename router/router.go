// Package router implements how to treat a connection which should be proxied, directed or denied.
package router

type policy int

const (
	Proxy policy = iota
	Direct
	Deny
)

type Rule struct {
	Address string
	Policy  policy
}

type Router struct {
	Rules   []*Rule
	Default policy
}

func (r *Router) Get(address string) policy {
	for _, r := range r.Rules {
		if r.Address == address {
			return r.Policy
		}
	}

	return r.Default
}
