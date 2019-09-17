// Copyright 2019 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
