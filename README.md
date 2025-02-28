# Nas-Mall
![build](https://github.com/NASKIDS/nas-mall/actions/workflows/go.yml/badge.svg)
![GitHub stars](https://img.shields.io/github/stars/NASKIDS/nas-mall/?style=social)
![GitHub forks](https://img.shields.io/github/forks/NASKIDS/nas-mall/?style=social)
![GitHub watchers](https://img.shields.io/github/watchers/NASKIDS/nas-mall/?style=social)
![GitHub repo size](https://img.shields.io/github/repo-size/NASKIDS/nas-mall/)
![GitHub language count](https://img.shields.io/github/languages/count/NASKIDS/nas-mall/)
![GitHub top language](https://img.shields.io/github/languages/top/NASKIDS/nas-mall/)
![GitHub last commit](https://img.shields.io/github/last-commit/NASKIDS/nas-mall/?color=red)
[中文](README_cn.md)

our project for Youth Training Camp

## Technology Stack
| technology | introduce |
|---------------|----|
| cwgo          | -  |
| kitex         | -  |
| [bootstrap](https://getbootstrap.com/docs/5.3/getting-started/introduction/) | Bootstrap is a powerful, feature-packed frontend toolkit. Build anything—from prototype to production—in minutes.  |
| Hertz         | -  |
| MySQL         | -  |
| Redis         | -  |
| ES            | -  |
| Prometheus    | -  |
| Jaeger        | -  |
| K8S        | -  |


## Biz Logic


## How to use
### Prepare 
List required
- Go
- IDE / Code Editor
- Docker
- [cwgo](https://github.com/cloudwego/cwgo)
- kitex `go install github.com/cloudwego/kitex/tool/cmd/kitex@latest`
- [Air](https://github.com/cosmtrek/air)
- ...

### Clone code
```
git clone ...
```

### Copy `.env` file
```
make init
```
*Note:*`You must generate and input SESSION_SECRET random value for session`

### Download go module
```
make tidy
```

### Start Docker Compose
```
make env-start
```
if you want to stop their docker application,you can run `make env-stop`.

### Run Service
This cmd must appoint a service.

*Note:* `Run the Go server using air. So it must be installed`
```
make run svc=`svcName`
```
### View Gomall Website
```
make open-gomall
```
### Check Registry
```
make open-consul
```
### Make Usage
```
make
```
## Contributors
- [rogerogers](https://github.com/rogerogers)
- [baiyutang](https://github.com/baiyutang)
