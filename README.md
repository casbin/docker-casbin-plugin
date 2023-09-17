# Docker Authorization Plugin Based on Casbin

[![Go Report Card](https://goreportcard.com/badge/github.com/casbin/casbin-authz-plugin)](https://goreportcard.com/report/github.com/casbin/casbin-authz-plugin) [![Build Status](https://travis-ci.org/casbin/casbin.svg?branch=master)](https://travis-ci.org/casbin/casbin) [![GoDoc](https://godoc.org/github.com/casbin/casbin-authz-plugin?status.svg)](https://godoc.org/github.com/casbin/casbin-authz-plugin)

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

## For "non-golang developer" users
```bash
$ apt install golang-go  # install go language
$ mkdir /usr/local/go
$ export GOPATH=/usr/local/go
```
- The installation command above is for Ubuntu, other distros may have different commands for installing go
- The export can be changed according to your satisfaction

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
$ mkdir examples
$ cp basic_model.conf examples/.
$ cp basic_policy.csv examples/.
$ ./casbin-authz-plugin
```

Below should be an example of display when command above is run:
```
2017/10/21 03:47:39 Current directory: /usr/lib/docker
2017/10/21 03:47:39 Casbin model: examples/basic_model.conf
2017/10/21 03:47:39 Casbin policy: examples/basic_policy.csv
2017/10/21 03:47:39 [Model:]
2017/10/21 03:47:39 p.p: obj, act
2017/10/21 03:47:39 e.e: some(where (p_eft == allow))
2017/10/21 03:47:39 m.m: r_obj == p_obj && r_act == p_act
2017/10/21 03:47:39 r.r: obj, act
2017/10/21 03:47:39 [Policy:]
2017/10/21 03:47:39 [p :  obj, act :  [[/_ping GET] [/v1.27/images/json GET]]]
```

## Enable the authorization plugin on docker engine

### Step-1: Determine where the systemd service of the plugin is located
```bash
$ systemctl status casbin-authz-plugin

● casbin-authz-plugin.service - Docker RBAC & ABAC Authorization Plugin based on Casbin
   Loaded: loaded (/lib/systemd/system/casbin-authz-plugin.service; disabled; vendor preset: enabled)
   Active: inactive (dead)
```
- You can see the directory on the **Loaded** label

### Step-2: Add the **WorkingDirectory** of th plugin's systemd service
```bash
$ vi /lib/systemd/system/casbin-authz-plugin.service

[Service]
WorkingDirectory=/usr/lib/docker
```
- If the service directory above is different than the one that returned from the `systemctl status casbin-authz-plugin`, please use the latter
- The `WorkingDirectory` may not be the one given depending on where you put the plugin

### Step-3: Run the plugin as a systemd service

```bash
$ systemctl daemon-reload
$ systemctl enable casbin-authz-plugin
$ systemctl start casbin-authz-plugin
```

### Step-4: Edit the **Execstart** of th plugin's systemd service
```
$ systemctl edit docker

[Service]
ExecStart=
ExecStart=/usr/bin/dockerd --authorization-plugin=casbin-authz-plugin
```
- If the service directory above is different than the one that returned from the `systemctl status docker`,  please use the latter 
- Just add `--authorization-plugin=casbin-authz-plugin` if there are more options on the pre-defined `ExecStart` please retain them

### Step-5: Restart docker engine

```bash
$ systemctl daemon-reload
$ systemctl restart docker
```

### Step-6 Activate the plugin logs:

```bash
$ journalctl -xe -u casbin-authz-plugin -f
```

### STEP-7 Do a quick test
```bash
$ docker images
```
- if `docker images` is denied, simply proceed to Step-8 for the solution

### Step-8 Changing the policy
```bash
$ vi /usr/lib/docker/examples/basic_policy.csv

p, /v1.29/images/json, GET

$ systemctl restart casbin-authz-plugin
```
- take note that **versioning** is also included on the authorization. The given policy states **/v1.27/**. So edit the version in `examples/basic_policy.csv` that the docker client is throwing which is shown in `journalctl` like `obj: /v1.29/images/json, act: GET res: denied`
- you can change the `$GOPATH` to the directory where you put the plugin from `go get`
- Check the logs for more confirmation

### Step-9 Test again:

```bash
$ docker images
$ docker ps
$ docker info
```
- If `docker images` is still denied please check STEP-8 more carefully
- These should smoothly enable 

## Stop and uninstall the plugin as a systemd service

NOTE: Before doing below, remove the authorization-plugin configuration added above and restart the docker daemon.

Removing the authorization plugin on docker

```bash
$ systemctl edit docker

#[Service]
#ExecStart=
#ExecStart=/usr/bin/dockerd --authorization-plugin=casbin-authz-plugin

$ systemctl restart docker
```

Stop the plugin service:

```bash
$ systemctl stop casbin-authz-plugin
$ systemctl disable casbin-authz-plugin
```

Uninstall the plugin service:

```bash
$ cd $GOPATH/src/github.com/casbin/casbin-authz-plugin
$ make uninstall
```

## Contact

If you have any issues or feature requests, please feel free to contact me at:
- https://github.com/casbin/casbin/issues
- hsluoyz@gmail.com

## License

Apache 2.0
