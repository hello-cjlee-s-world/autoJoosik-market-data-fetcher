# auto-joosik-market-data-fetcher

# Build
# linux build / amd64
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -ldflags="-s -w" -trimpath -o .\build\auto-joosik-market-data-fetcher .\cmd\main.go

# window build / amd64
$env:GOOS="windows"; $env:GOARCH="amd64"; $env:CGO_ENABLED="0" 
go build -ldflags="-s -w" -trimpath -o .\build\auto-joosik-market-data-fetcher.exe .\cmd\main.go
# 설정파일 autoJoosik_market_data_fetcher_conf.yml 생성 후 실행 명령어 
auto-joosik-market-data-fetcher.exe -config "autoJoosik_market_data_fetcher_conf.yml"


