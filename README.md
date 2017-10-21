# Docker RBAC & ABAC Authorization Plugin based on Casbin [![Go Report Card](https://goreportcard.com/badge/github.com/casbin/casbin-authz-plugin)](https://goreportcard.com/report/github.com/casbin/casbin-authz-plugin) [![Build Status](https://travis-ci.org/casbin/casbin.svg?branch=master)](https://travis-ci.org/casbin/casbin) [![GoDoc](https://godoc.org/github.com/casbin/casbin-authz-plugin?status.svg)](https://godoc.org/github.com/casbin/casbin-authz-plugin)

This plugin controls the access to Docker commands based on authorization policy. The functionality of authorization is provided by [Casbin](https://github.com/casbin/casbin). Since Docker doesn't perform authentication by now, there's no user information when executing Docker commands. The access that Casbin plugin can control is actually what HTTP method can be performed on what URL path.

For example, when you run ``docker images`` command, the underlying request is really like:

```
/v1.27/images/json, GET
```

So Casbin plugin helps you decide whether ``GET`` can be performed on ``/v1.27/images/json`` base on the policy rules you write. The policy file is ``basic_policy.csv`` co-located with the plugin binary by default. And its content is:

```
p, /v1.27/images/json, GET
```

The above policy grants anyone to perform ``GET`` on ``/v1.27/images/json``, and deny all other requests. The response should be like below:

```bash
$ docker images
REPOSITORY          TAG                 IMAGE ID            CREATED             SIZE
hello-world         latest              48b5124b2768        3 months ago        1.84 kB

$ docker info
Error response from daemon: authorization denied by plugin casbin-authz-plugin: Access denied by casbin plugin
```

The built-in Casbin model is:

```ini
[request_definition]
r = obj, act

[policy_definition]
p = obj, act

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = r.obj == p.obj && r.act == p.act
```

The built-in Casbin policy is:

```csv
p, /_ping, GET
p, /v1.27/images/json, GET
```

For more information about the Casbin model and policy usage like RBAC, ABAC, please refer to: https://github.com/casbin/casbin

## Build

```bash
$ go get github.com/casbin/casbin-authz-plugin
$ cd $GOPATH/src/github.com/casbin/casbin-authz-plugin
$ make
$ sudo make install
```

## Run

### Run the plugin directly in a shell

```bash
$ cd /usr/lib/docker
$ ./casbin-authz-plugin
```

See whether the plugin starts correctly:

```bash
$ journalctl -xe -u casbin-authz-plugin -f
```

## Enable the authorization plugin on docker engine

### Step-1: Determine where the systemd service of the plugin is located
```bash
$ systemctl status casbin-authz-plugin

Loaded: loaded (/lib/systemd/system/casbin-authz-plugin.service; disabled; vendor preset: enabled)
   Active: active (running) since Sun 2017-10-15 14:21:35 UTC; 5 days ago
 Main PID: 64275 (casbin-authz-pl)
   CGroup: /system.slice/casbin-authz-plugin.service
           └─64275 /usr/lib/docker/casbin-authz-plugin
```
- You can see the directory on the **Loaded** label

### Step-2: Edit the **WorkingDirectory** of th plugin's systemd service
```bash
$ vi /lib/systemd/system/casbin-authz-plugin.service
```
- Make `/usr/local/go/bin/src/github.com/casbin/casbin-authz-plugin` as the value
- If the service directory above is different than the one that returned from the `systemctl status casbin-authz-plugin` please use the latter 

### Step-3: Run the plugin as a systemd service

```bash
$ systemctl daemon-reload
$ systemctl enable casbin-authz-plugin
$ systemctl start casbin-authz-plugin
```

### Step-4: Determine where the docker service is located
```bash
$ systemctl status docker

 docker.service
   Loaded: loaded (/etc/systemd/system/docker.service; enabled; vendor preset: enabled)
   Active: active (running) since Sun 2017-10-15 14:21:35 UTC; 5 days ago
 Main PID: 64276 (dockerd)
   CGroup: /system.slice/docker.service
           ├─  808 docker-containerd-shim a795cd15bbf1e314f45822047b00b43f8848ebfe26091c038f9278eda2f5b306 /var/run/docker/libcontainerd
           ├─64276 dockerd -H tcp://0.0.0.0:2376 -H unix:///var/run/docker.sock --storage-driver aufs --tlsverify --tlscacert /etc/docke
```

### Step-5: Edit the **Execstart** of th plugin's systemd service
```
[Service]
ExecStart=
ExecStart=/usr/bin/docker daemon -H tcp://0.0.0.0:2376 --authorization-plugin=casbin-authz-plugin
```
- Just add `--authorization-plugin=casbin-authz-plugin` if there are more options on the pre-defined `ExecStart` please retain them

### Step-6: Add authorization plugin to the docker engine configuration 

Please add the following cmdline flag to your docker engine (e.g. ExecStart line ``/lib/systemd/system/docker.service``)

```bash
--authorization-plugin casbin-authz-plugin
```

### Step-7: Restart docker engine

```bash
$ systemctl daemon-reload
$ systemctl restart docker
```

## Stop and uninstall the plugin as a systemd service

NOTE: Before doing below, remove the authorization-plugin configuration added above and restart the docker daemon.

Stop the plugin service:

```bash
$ systemctl stop casbin-authz-plugin
$ systemctl disable casbin-authz-plugin
```

Uninstall the plugin service:

```bash
$ make uninstall
```

## Contact

If you have any issues or feature requests, please feel free to contact me at:
- https://github.com/casbin/casbin/issues
- hsluoyz@gmail.com

## License

Apache 2.0
