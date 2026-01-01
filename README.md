# auto-joosik-market-data-fetcher

# Build
# linux build / amd64
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -ldflags="-s -w" -trimpath -o .\build\auto-joosik-market-data-fetcher .\cmd\main.go

# window build / amd64
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -ldflags="-s -w" -trimpath -o .\build\auto-joosik-market-data-fetcher .\cmd\main.go


