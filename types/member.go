package types

import "github.com/hashicorp/serf/serf"

type Member struct {
	serf.Member

	Id         string `json:"id"`
	StatusText string `json:"statusText"`
}
