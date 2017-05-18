// Docker RBAC & ABAC Authorization Plugin based on Casbin.
// Allows only authorized Docker operations based on access control policy.
// AUTHOR: Yang Luo <hsluoyz@gmail.com>
// Powered by Casbin: https://github.com/casbin/casbin

package main

import (
	"github.com/docker/go-plugins-helpers/authorization"
	"github.com/casbin/casbin"
	"log"
	"net/url"
)

// Casbin Authorization Plugin struct definition
type CasbinAuthZPlugin struct {
	// Casbin enforcer
	enforcer *casbin.Enforcer
}

// Create a new casbin authorization plugin
func newPlugin(casbinConfig string) (*CasbinAuthZPlugin, error) {
	plugin := &CasbinAuthZPlugin{}

	plugin.enforcer = casbin.NewEnforcer(casbinConfig)

	return plugin, nil
}

// Authorizes the docker client command.
// The command is allowed only if it matches a Casbin policy rule.
// Otherwise, the request is denied!
func (plugin *CasbinAuthZPlugin) AuthZReq(req authorization.Request) authorization.Response {
	// Parse request and the request body
	reqURI, _ := url.QueryUnescape(req.RequestURI)
	reqURL, _ := url.ParseRequestURI(reqURI)

	obj := reqURL.String()
	act := req.RequestMethod

	if plugin.enforcer.Enforce(obj, act) {
		log.Println("obj:", obj, ", act:", act, "res: allowed")
		return authorization.Response{Allow: true}
	} else {
		log.Println("obj:", obj, ", act:", act, "res: denied")
		return authorization.Response{Allow: false, Msg: "Access denied by casbin plugin"}
	}
}

// Authorizes the docker client response.
// All responses are allowed by default.
func (plugin *CasbinAuthZPlugin) AuthZRes(req authorization.Request) authorization.Response {
	// Allowed by default.
	return authorization.Response{Allow: true}
}
