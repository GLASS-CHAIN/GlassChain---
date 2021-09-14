[![API Reference](
https://camo.githubusercontent.com/915b7be44ada53c290eb157634330494ebe3e30a/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f6c616e672f6764646f3f7374617475732e737667
)](https://godoc.org/github.com/33cn/plugin)
[![pipeline status](https://api.travis-ci.org/33cn/plugin.svg?branch=master)](https://travis-ci.org/33cn/plugin/)
[![Go Report Card](https://goreportcard.com/badge/github.com/33cn/plugin?branch=master)](https://goreportcard.com/report/github.com/33cn/plugin)
[![Windows Build Status](https://ci.appveyor.com/api/projects/status/github/33cn/plugin?svg=true&branch=master&passingText=Windows%20-%20OK&failingText=Windows%20-%20failed&pendingText=Windows%20-%20pending)](https://ci.appveyor.com/project/33cn/plugin)
[![codecov](https://codecov.io/gh/33cn/plugin/branch/master/graph/badge.svg)](https://codecov.io/gh/33cn/plugin)

# chain33 Official Plugin System（v6.4.0）

* chain33 address: https://github.com/33cn/chain33
* chain33 Official Web: https://chain.33.cn

### Environmental Science

```
Installation requiredgolang1.13 or latest

```

#### Platform supporting make file

```
//Turn on mod function
export GO111MODULE=on

//Domestic users need to import Aliyun Agent to download dependent packages
export GOPROXY=https://mirrors.aliyun.com/goproxy

make
```
You can complete the compilation installation

```
Note: domestic users need to import a proxy before they can get dependent packages. The mod function is turned on by default in Makefile
```

## Run

```
./chain33 -f chain33.toml
```
Note that the default configuration will connect to the chain33 test network

## Note:

Using mod to manage dependency packages is primarily a wall-flipping problem

To solve the problem of package dependency wall-flipping download, we provide AliCloud Agent.


## Contribution code:

Detailed steps are available https://github.com/33cn/chain33

Here are just the simple steps:

#### Preparations:

*First click on the fork icon in the upper right corner and put chain33 fork into your own branch like mine is vipwzw/plugin

* `git clonehttps://github.com/vipwzw/plugin.git $GOPATH/src/github.com/33cn/plugin`
```
Note: Clone to $GOPATH/src/github.com/33cn/plugin is required here, otherwise go package path will not be found
```

When clone is complete, execute
```
make addupstream
```

#### Create Branches to Develop New Functions

```
make branch b=branch_dev_name
```
#### Submit Code

```
make push b=branch_dev_name m="hello world"
```
If m is not set, git commit commands will not be executed

#### Test Code
Like plugin/dapp/relay, write your own plug-ins Makefile and build.sh in the CMD directory
Write testcase and docker-compose configuration files in the build directory.
The rules for testcase refer to plugin/dapp/testcase_compose_rule.md

Users can set their own plugin's DAPP variable in travis's own project. If DAPP is set to relay, travis's run relay's testcase

