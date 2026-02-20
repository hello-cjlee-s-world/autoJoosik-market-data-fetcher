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
## Strategy Engine (Gate + Score)
- 위치: `internal/strategyengine`
- 파이프라인: **Watchlist 엔진(뉴스AI/거래량/수급)** → **Entry 엔진(게이트 + 가중치 점수)** → **Exit/Risk 엔진(손절/익절/트레일링/점수붕괴/일손실 제한)**
- 확장 방식: `GateEvaluator`, `FactorEvaluator` 인터페이스로 모듈 플러그인 추가
- 지표 지원: MA(5/20/60), RSI, MACD, Bollinger, ATR, VWAP
- 설정: `internal/config/autoJoosik_market_data_fetcher_conf.yml` 의 `strategyEngine` 섹션에서 임계값/가중치 관리
