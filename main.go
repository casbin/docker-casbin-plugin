// Docker RBAC & ABAC Authorization Plugin based on Casbin.
// Allows only authorized Docker operations based on access control policy.
// AUTHOR: Yang Luo <hsluoyz@qq.com>
// Powered by Casbin: https://github.com/hsluoyz/casbin

package main

import (
	"flag"
	"github.com/docker/go-plugins-helpers/authorization"
	"log"
	"os/user"
	"strconv"
)

const (
	pluginSocket = "/run/docker/plugins/casbin-authz-plugin.sock"
)

var (
	casbinConfig = flag.String("config", "casbin.conf", "Specifies the Casbin configuration file")
)

func main() {
	// Parse command line options.
	flag.Parse()
	log.Println("Casbin config:", *casbinConfig)

	// Create Casbin authorization plugin
	plugin, err := newPlugin(*casbinConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Start service handler on the local sock
	u, _ := user.Lookup("root")
	gid, _ := strconv.Atoi(u.Gid)
	handler := authorization.NewHandler(plugin)
	if err := handler.ServeUnix(pluginSocket, gid); err != nil {
		log.Fatal(err)
	}
}
