# auto-joosik-market-data-fetcher
 
# Build
# linux build / amd64
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -ldflags="-s -w" -trimpath -o .\build\autoJoosik-market-data-fetcher .\cmd\main.go

# window build / amd64
$env:GOOS="windows"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0" 
go build -ldflags="-s -w" -trimpath -o .\build\autoJoosik-market-data-fetcher.exe .\cmd\main.go 
# 설정파일 autoJoosik_market_data_fetcher_conf.yml 생성 후 실행 명령어 
auto-joosik-market-data-fetcher.exe -config "autoJoosik_market_data_fetcher_conf.yml"


# 설정파일 항목
server:
    port: 1323
    title: autoJoosik_market_data_fetcher
    contextPath: /auto

logging:
    level: debug
    filename: ./logs/app.log
    maxSize: 10
    maxBackups: 10
    maxAge: 10
    compress: false
    consoleOutput: true

database:
    user: postgres
    password: test
    host: localhost
    port: 12345
    name: test
    maxIdleConn: 10
    connMaxLifetime: 10
    sslMode: disable
    maximumPoolSize: 1

kiwoomApi:
    appKey: appKey
    secretKey: secretKey