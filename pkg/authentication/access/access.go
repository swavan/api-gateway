package access

import (
	"github.com/casbin/casbin/util"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/jmoiron/sqlx"
)

type API interface {
	Enforcer() *casbin.Enforcer
}

type AccessControl struct {
	enforcer *casbin.Enforcer
}

func New(authModalConfig string, dep *sqlx.DB, migration bool) (API, error) {
	m := model.NewModel()
	m.AddDef("r", "r", "sub, dom, obj, act")
	m.AddDef("p", "p", "sub, dom, obj, act")
	m.AddDef("g", "g", "_, _, _")
	m.AddDef("e", "e", "some(where (p.eft == allow))")

	m.AddDef("m", "m", `(g(r.sub, p.sub, r.dom)) && (regexMatch(r.dom, p.dom) && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == 'ANY'))  || (r.sub == 'root')`)
	// m =             (g(r.sub, p.sub, r.dom) || p.sub == "Anonymous") && regexMatch(r.dom, p.dom) && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "ANY")

	adapter, err := NewAdapter(dep.DB, "access_rule_store")
	if err != nil {
		return nil, err
	}

	eff, err := casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}
	eff.AddNamedDomainMatchingFunc("g", "", util.KeyMatch)
	return &AccessControl{
		enforcer: eff,
	}, nil
}

func (a *AccessControl) Enforcer() *casbin.Enforcer {
	return a.enforcer
}
