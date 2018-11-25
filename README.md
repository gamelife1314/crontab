<p align="center">
    <img alt="crontab icon" src="./crontab.jpeg" width="256px">
</p>

<p align="center">
 <a href="https://travis-ci.com/gamelife1314/crontab/"><img alt="build status" src="https://travis-ci.com/gamelife1314/crontab.svg?branch=master"></a>
 <a href="#"><img alt="language" src="https://img.shields.io/badge/language-go-orange.svg"></a>
 <a href="#"><img alt="develop status" src="https://img.shields.io/badge/status-developing-red.svg"></a>
 <a href="#"><img alt="develop status" src="https://img.shields.io/badge/version-0.1.0-green.svg"></a>
 <a href="#"><img alt="License" src="https://img.shields.io/badge/license-MIT-blue.svg"></a>
</p>

> A Distributed Task Scheduler like Cron.

** Now, It's under developing. **

### Dependencies
1. [etcd](https://coreos.com/etcd/)
2. [mongo](https://github.com/mongodb/mongo)

### Install

If you have already installed go in your local machine, you can install `crond` and `cron` by `go install`：

for cron:

    go build -o client github.com/gamelife1314/crontab/cron/main

for crond:

    go build -o master github.com/gamelife1314/crontab/crond/main

Or, you can download [prebuilt binary](https://github.com/gamelife1314/crontab/releases) files.

### Start

for crond: `./master -config  crond.yaml`

for client: `./client -config  cron.yaml`
