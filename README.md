# Nas-Mall
[中文](README_cn.md)

![build](https://github.com/NASKIDS/nas-mall/actions/workflows/go.yml/badge.svg)


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
