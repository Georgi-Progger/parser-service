package proxy

type Proxy struct {
	Id     uint
	Body   string `json:"proxy_body"`
	Active bool   `json:"active"`
}
