// Docker RBAC & ABAC Authorization Plugin based on Casbin.
// Allows only authorized Docker operations based on access control policy.
// AUTHOR: Yang Luo <hsluoyz@gmail.com>
// Powered by Casbin: https://github.com/casbin/casbin

package main

import (
	"flag"
	"log"
	"os"
	"os/user"
	"strconv"

	"github.com/docker/go-plugins-helpers/authorization"
)

const (
	pluginSocket = "/run/docker/plugins/casbin-authz-plugin.sock"
)

var (
	casbinModel = flag.String("model", "examples/basic_model.conf", "Specifies the Casbin model file")
	casbinPolicy = flag.String("policy", "examples/basic_policy.csv", "Specifies the Casbin policy file")
)

func main() {
	// Parse command line options.
	flag.Parse()
	pwd, _ := os.Getwd()
	log.Println("Current directory:", pwd)
	log.Println("Casbin model:", *casbinModel)
	log.Println("Casbin policy:", *casbinPolicy)

	// Create Casbin authorization plugin
	plugin, err := newPlugin(*casbinModel, *casbinPolicy)
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
