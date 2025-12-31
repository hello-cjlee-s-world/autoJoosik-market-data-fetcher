DROP TABLE stock_info;
DROP TABLE trade_info_log;
DROP TABLE orderbook_log;
DROP TABLE stock_daily_log;
DROP TABLE stock_tick_log;
DROP TABLE account_profit_log;
DROP TABLE atn_stk_infr;
DROP TABLE schedule_info;
DROP TABLE tb_virtual_asset;
DROP TABLE tb_virtual_account;
DROP TABLE tb_virtual_account;
DROP TABLE tb_virtual_order;
DROP TABLE tb_virtual_trade_log;

DROP INDEX IF EXISTS idx_trade_info_log_tm;
DROP INDEX IF EXISTS idx_trade_info_log_stextp;
DROP INDEX IF EXISTS idx_trade_info_log_created_at;
DROP INDEX IF EXISTS idx_orderbook_log_tm;
DROP INDEX IF EXISTS idx_orderbook_log_created_at;
DROP INDEX IF EXISTS idx_stock_daily_log_cd_date;
DROP INDEX IF EXISTS idx_stock_tick_log_cd_tm;
DROP INDEX IF EXISTS idx_account_profit_dt_cd;
DROP INDEX IF EXISTS uq_virtual_asset;
DROP INDEX IF EXISTS uq_virtual_account_user_name;
DROP INDEX IF EXISTS uq_virtual_order_client;
DROP INDEX IF EXISTS idx_trade_log_order;
DROP INDEX IF EXISTS idx_trade_log_account_stk;


    CREATE TABLE stock_info ( --주식기본정보요청
                            stk_cd              VARCHAR(20) PRIMARY KEY,   -- 종목코드
                            stk_nm              VARCHAR(40),              -- 종목명
                            setl_mm             VARCHAR(20),              -- 결산월
                            fav                 VARCHAR(20),              -- 액면가
                            cap                 VARCHAR(20),              -- 자본금
                            flo_stk             VARCHAR(20),              -- 상장주식
                            crd_rt              VARCHAR(20),              -- 신용비율
                            oyr_hgst            VARCHAR(20),              -- 연중최고
                            oyr_lwst            VARCHAR(20),              -- 연중최저
                            mac                 VARCHAR(20),              -- 시가총액
                            mac_wght            VARCHAR(20),              -- 시가총액비중
                            for_exh_rt          VARCHAR(20),              -- 외인소진률
                            repl_pric           VARCHAR(20),              -- 대용가
                            per                 VARCHAR(20),              -- PER
                            eps                 VARCHAR(20),              -- EPS
                            roe                 VARCHAR(20),              -- ROE
                            pbr                 VARCHAR(20),              -- PBR
                            ev                  VARCHAR(20),              -- EV
                            bps                 VARCHAR(20),              -- BPS
                            sale_amt            VARCHAR(20),              -- 매출액
                            bus_pro             VARCHAR(20),              -- 영업이익
                            cup_nga             VARCHAR(20),              -- 당기순이익
                            "250hgst"           VARCHAR(20),              -- 250최고
                            "250lwst"           VARCHAR(20),              -- 250최저
                            high_pric           VARCHAR(20),              -- 고가
                            open_pric           VARCHAR(20),              -- 시가
                            low_pric            VARCHAR(20),              -- 저가
                            upl_pric            VARCHAR(20),              -- 상한가
                            lst_pric            VARCHAR(20),              -- 하한가
                            base_pric           VARCHAR(20),              -- 기준가
                            exp_cntr_pric       VARCHAR(20),              -- 예상체결가
                            exp_cntr_qty        VARCHAR(20),              -- 예상체결수량
                            "250hgst_pric_dt"   VARCHAR(20),              -- 250최고가일
                            "250hgst_pric_pre_rt" VARCHAR(20),            -- 250최고가대비율
                            "250lwst_pric_dt"   VARCHAR(20),              -- 250최저가일
                            "250lwst_pric_pre_rt" VARCHAR(20),            -- 250최저가대비율
                            cur_prc             VARCHAR(20),              -- 현재가
                            pre_sig             VARCHAR(20),              -- 대비기호
                            pred_pre            VARCHAR(20),              -- 전일대비
                            flu_rt              VARCHAR(20),              -- 등락율
                            trde_qty            VARCHAR(20),              -- 거래량
                            trde_pre            VARCHAR(20),              -- 거래대비
                            fav_unit            VARCHAR(20),              -- 액면가단위
                            dstr_stk            VARCHAR(20),              -- 유통주식
                            dstr_rt             VARCHAR(20)               -- 유통비율
);

CREATE TABLE trade_info_log ( --체결정보
                            tm               TIMESTAMP,             -- 시간 (API 그대로 저장)
                            cur_prc          NUMERIC(18,2),           -- 현재가
                            pred_pre         NUMERIC(18,2),           -- 전일대비
                            pre_rt           NUMERIC(10,4),           -- 대비율 (%)
                            pri_sel_bid_unit NUMERIC(18,2),           -- 우선매도호가단위
                            pri_buy_bid_unit NUMERIC(18,2),           -- 우선매수호가단위
                            cntr_trde_qty    BIGINT,                  -- 체결거래량
                            sign             VARCHAR(20),             -- sign
                            acc_trde_qty     BIGINT,                  -- 누적거래량
                            acc_trde_prica   NUMERIC(20,2),           -- 누적거래대금
                            cntr_str         NUMERIC(10,2),           -- 체결강도
                            stex_tp          VARCHAR(10),             -- 거래소구분 (KRX, NXT, 통합)
                            stk_cd           VARCHAR(20) ,            -- 종목코드
                            created_at       TIMESTAMPTZ DEFAULT NOW() -- 로그 적재 시각
);
CREATE INDEX idx_trade_info_log_tm ON trade_info_log (tm);
CREATE INDEX idx_trade_info_log_stextp ON trade_info_log (stex_tp);
CREATE INDEX idx_trade_info_log_created_at ON trade_info_log (created_at);
ALTER TABLE trade_info_log ADD CONSTRAINT trade_info_log_unique UNIQUE (stk_cd, tm, cur_prc, cntr_trde_qty);

CREATE TABLE orderbook_log ( -- 주식호가
                               id                   BIGSERIAL PRIMARY KEY,     -- 내부 PK
                               bid_req_base_tm      VARCHAR(20),               -- 호가잔량기준시간 (호가시간)

    -- 매도호가 (10차~1차)
                               sel_10th_pre_req_pre NUMERIC(18,2),
                               sel_10th_pre_req     BIGINT,
                               sel_10th_pre_bid     NUMERIC(18,2),
                               sel_9th_pre_req_pre  NUMERIC(18,2),
                               sel_9th_pre_req      BIGINT,
                               sel_9th_pre_bid      NUMERIC(18,2),
                               sel_8th_pre_req_pre  NUMERIC(18,2),
                               sel_8th_pre_req      BIGINT,
                               sel_8th_pre_bid      NUMERIC(18,2),
                               sel_7th_pre_req_pre  NUMERIC(18,2),
                               sel_7th_pre_req      BIGINT,
                               sel_7th_pre_bid      NUMERIC(18,2),
                               sel_6th_pre_req_pre  NUMERIC(18,2),
                               sel_6th_pre_req      BIGINT,
                               sel_6th_pre_bid      NUMERIC(18,2),
                               sel_5th_pre_req_pre  NUMERIC(18,2),
                               sel_5th_pre_req      BIGINT,
                               sel_5th_pre_bid      NUMERIC(18,2),
                               sel_4th_pre_req_pre  NUMERIC(18,2),
                               sel_4th_pre_req      BIGINT,
                               sel_4th_pre_bid      NUMERIC(18,2),
                               sel_3th_pre_req_pre  NUMERIC(18,2),
                               sel_3th_pre_req      BIGINT,
                               sel_3th_pre_bid      NUMERIC(18,2),
                               sel_2th_pre_req_pre  NUMERIC(18,2),
                               sel_2th_pre_req      BIGINT,
                               sel_2th_pre_bid      NUMERIC(18,2),
                               sel_1th_pre_req_pre  NUMERIC(18,2),
                               sel_fpr_req          BIGINT,                    -- 매도최우선잔량
                               sel_fpr_bid          NUMERIC(18,2),             -- 매도최우선호가

    -- 매수호가 (1차~10차)
                               buy_fpr_bid          NUMERIC(18,2),             -- 매수최우선호가
                               buy_fpr_req          BIGINT,                    -- 매수최우선잔량
                               buy_1th_pre_req_pre  NUMERIC(18,2),
                               buy_2th_pre_bid      NUMERIC(18,2),
                               buy_2th_pre_req      BIGINT,
                               buy_2th_pre_req_pre  NUMERIC(18,2),
                               buy_3th_pre_bid      NUMERIC(18,2),
                               buy_3th_pre_req      BIGINT,
                               buy_3th_pre_req_pre  NUMERIC(18,2),
                               buy_4th_pre_bid      NUMERIC(18,2),
                               buy_4th_pre_req      BIGINT,
                               buy_4th_pre_req_pre  NUMERIC(18,2),
                               buy_5th_pre_bid      NUMERIC(18,2),
                               buy_5th_pre_req      BIGINT,
                               buy_5th_pre_req_pre  NUMERIC(18,2),
                               buy_6th_pre_bid      NUMERIC(18,2),
                               buy_6th_pre_req      BIGINT,
                               buy_6th_pre_req_pre  NUMERIC(18,2),
                               buy_7th_pre_bid      NUMERIC(18,2),
                               buy_7th_pre_req      BIGINT,
                               buy_7th_pre_req_pre  NUMERIC(18,2),
                               buy_8th_pre_bid      NUMERIC(18,2),
                               buy_8th_pre_req      BIGINT,
                               buy_8th_pre_req_pre  NUMERIC(18,2),
                               buy_9th_pre_bid      NUMERIC(18,2),
                               buy_9th_pre_req      BIGINT,
                               buy_9th_pre_req_pre  NUMERIC(18,2),
                               buy_10th_pre_bid     NUMERIC(18,2),
                               buy_10th_pre_req     BIGINT,
                               buy_10th_pre_req_pre NUMERIC(18,2),

    -- 총잔량
                               tot_sel_req_jub_pre  NUMERIC(18,2),             -- 총매도잔량직전대비
                               tot_sel_req          BIGINT,                    -- 총매도잔량
                               tot_buy_req          BIGINT,                    -- 총매수잔량
                               tot_buy_req_jub_pre  NUMERIC(18,2),             -- 총매수잔량직전대비

    -- 시간외 잔량
                               ovt_sel_req_pre      NUMERIC(18,2),             -- 시간외매도잔량대비
                               ovt_sel_req          BIGINT,                    -- 시간외매도잔량
                               ovt_buy_req          BIGINT,                    -- 시간외매수잔량
                               ovt_buy_req_pre      NUMERIC(18,2),             -- 시간외매수잔량대비

                               created_at           TIMESTAMPTZ DEFAULT NOW()  -- 로그 적재 시각
);
CREATE INDEX idx_orderbook_log_tm ON orderbook_log (bid_req_base_tm);
CREATE INDEX idx_orderbook_log_created_at ON orderbook_log (created_at);


CREATE TABLE stock_daily_log ( --주식일주월시분
                                 id BIGSERIAL PRIMARY KEY,          -- 로그 PK
                                 stk_cd VARCHAR(20) NOT NULL,       -- 종목코드 (외래키로도 연결 가능)
                                 date DATE NOT NULL,                -- 날짜

                                 open_pric NUMERIC(18,2),           -- 시가
                                 high_pric NUMERIC(18,2),           -- 고가
                                 low_pric NUMERIC(18,2),            -- 저가
                                 close_pric NUMERIC(18,2),          -- 종가
                                 pre NUMERIC(18,2),                 -- 대비
                                 flu_rt NUMERIC(8,4),               -- 등락률(%)
                                 trde_qty BIGINT,                   -- 거래량
                                 trde_prica NUMERIC(20,2),          -- 거래대금

                                 for_poss BIGINT,                   -- 외인보유
                                 for_wght NUMERIC(8,4),             -- 외인비중(%)
                                 for_netprps BIGINT,                -- 외인순매수
                                 orgn_netprps BIGINT,               -- 기관순매수
                                 ind_netprps BIGINT,                -- 개인순매수
                                 crd_remn_rt NUMERIC(8,4),          -- 신용잔고율(%)
                                 frgn BIGINT,                       -- 외국계
                                 prm BIGINT,                        -- 프로그램

                                 created_at TIMESTAMP DEFAULT now() -- 데이터 적재 시각
);
CREATE INDEX idx_stock_daily_log_cd_date ON stock_daily_log(stk_cd, date);

CREATE TABLE stock_tick_log ( -- 주식 틱차트
                                id BIGSERIAL PRIMARY KEY,          -- 로그 PK
                                stk_cd VARCHAR(20) NOT NULL,       -- 종목코드
                                last_tic_cnt INT,                  -- 마지막 틱 갯수 (서버 조회시 기준)

                                cur_prc NUMERIC(18,2),             -- 현재가
                                trde_qty BIGINT,                   -- 거래량
                                cntr_tm TIMESTAMP,                 -- 체결시간
                                open_pric NUMERIC(18,2),           -- 시가
                                high_pric NUMERIC(18,2),           -- 고가
                                low_pric NUMERIC(18,2),            -- 저가

                                upd_stkpc_tp INT,                  -- 수정주가구분 (비트 플래그: 1,2,4,8,...)
                                upd_rt NUMERIC(8,4),               -- 수정비율
                                bic_inds_tp VARCHAR(20),           -- 대업종구분
                                sm_inds_tp VARCHAR(20),            -- 소업종구분
                                stk_infr VARCHAR(50),              -- 종목정보
                                upd_stkpc_event VARCHAR(50),       -- 수정주가이벤트
                                pred_close_pric NUMERIC(18,2),     -- 전일종가

                                created_at TIMESTAMP DEFAULT now() -- 적재 시간
);
CREATE INDEX idx_stock_tick_log_cd_tm ON stock_tick_log(stk_cd, cntr_tm);

CREATE TABLE account_profit_log ( -- 계좌수익률
                                    id BIGSERIAL PRIMARY KEY,           -- 로그 PK
                                    dt DATE NOT NULL,                   -- 일자
                                    stk_cd VARCHAR(20) NOT NULL,        -- 종목코드
                                    stk_nm VARCHAR(40),                 -- 종목명

                                    cur_prc NUMERIC(18,2),              -- 현재가
                                    pur_pric NUMERIC(18,2),             -- 매입가
                                    pur_amt NUMERIC(20,2),              -- 매입금액
                                    rmnd_qty BIGINT,                    -- 보유수량

                                    tdy_sel_pl NUMERIC(20,2),           -- 당일매도손익
                                    tdy_trde_cmsn NUMERIC(20,2),        -- 당일매매수수료
                                    tdy_trde_tax NUMERIC(20,2),         -- 당일매매세금

                                    crd_tp VARCHAR(20),                 -- 신용구분
                                    loan_dt DATE,                       -- 대출일
                                    setl_remn NUMERIC(20,2),            -- 결제잔고
                                    clrn_alow_qty BIGINT,               -- 청산가능수량
                                    crd_amt NUMERIC(20,2),              -- 신용금액
                                    crd_int NUMERIC(20,2),              -- 신용이자
                                    expr_dt DATE,                       -- 만기일

                                    created_at TIMESTAMP DEFAULT now()  -- 적재 시간
);
CREATE INDEX idx_account_profit_dt_cd ON account_profit_log(dt, stk_cd);

CREATE TABLE atn_stk_infr ( --관심종목정보요청
                              stk_cd          VARCHAR(20),   -- 종목코드
                              stk_nm          VARCHAR(40),   -- 종목명
                              cur_prc         VARCHAR(20),   -- 현재가
                              base_pric       VARCHAR(20),   -- 기준가
                              pred_pre        VARCHAR(20),   -- 전일대비
                              pred_pre_sig    VARCHAR(20),   -- 전일대비기호
                              flu_rt          VARCHAR(20),   -- 등락율
                              trde_qty        VARCHAR(20),   -- 거래량
                              trde_prica      VARCHAR(20),   -- 거래대금
                              cntr_qty        VARCHAR(20),   -- 체결량
                              cntr_str        VARCHAR(20),   -- 체결강도
                              pred_trde_qty_pre VARCHAR(20), -- 전일거래량대비
                              sel_bid         VARCHAR(20),   -- 매도호가
                              buy_bid         VARCHAR(20),   -- 매수호가
                              sel_1th_bid     VARCHAR(20),   -- 매도1차호가
                              sel_2th_bid     VARCHAR(20),   -- 매도2차호가
                              sel_3th_bid     VARCHAR(20),   -- 매도3차호가
                              sel_4th_bid     VARCHAR(20),   -- 매도4차호가
                              sel_5th_bid     VARCHAR(20),   -- 매도5차호가
                              buy_1th_bid     VARCHAR(20),   -- 매수1차호가
                              buy_2th_bid     VARCHAR(20),   -- 매수2차호가
                              buy_3th_bid     VARCHAR(20),   -- 매수3차호가
                              buy_4th_bid     VARCHAR(20),   -- 매수4차호가
                              buy_5th_bid     VARCHAR(20),   -- 매수5차호가
                              upl_pric        VARCHAR(20),   -- 상한가
                              lst_pric        VARCHAR(20),   -- 하한가
                              open_pric       VARCHAR(20),   -- 시가
                              high_pric       VARCHAR(20),   -- 고가
                              low_pric        VARCHAR(20),   -- 저가
                              close_pric      VARCHAR(20),   -- 종가
                              cntr_tm         VARCHAR(20),   -- 체결시간
                              exp_cntr_pric   VARCHAR(20),   -- 예상체결가
                              exp_cntr_qty    VARCHAR(20),   -- 예상체결량
                              cap             VARCHAR(20),   -- 자본금
                              fav             VARCHAR(20),   -- 액면가
                              mac             VARCHAR(20),   -- 시가총액
                              stkcnt          VARCHAR(20),   -- 주식수
                              bid_tm          VARCHAR(20),   -- 호가시간
                              dt              VARCHAR(20),   -- 일자
                              pri_sel_req     VARCHAR(20),   -- 우선매도잔량
                              pri_buy_req     VARCHAR(20),   -- 우선매수잔량
                              pri_sel_cnt     VARCHAR(20),   -- 우선매도건수
                              pri_buy_cnt     VARCHAR(20),   -- 우선매수건수
                              tot_sel_req     VARCHAR(20),   -- 총매도잔량
                              tot_buy_req     VARCHAR(20),   -- 총매수잔량
                              tot_sel_cnt     VARCHAR(20),   -- 총매도건수
                              tot_buy_cnt     VARCHAR(20),   -- 총매수건수
                              prty            VARCHAR(20),   -- 패리티
                              gear            VARCHAR(20),   -- 기어링
                              pl_qutr         VARCHAR(20),   -- 손익분기
                              cap_support     VARCHAR(20),   -- 자본지지
                              elwexec_pric    VARCHAR(20),   -- ELW행사가
                              cnvt_rt         VARCHAR(20),   -- 전환비율
                              elwexpr_dt      VARCHAR(20),   -- ELW만기일
                              cntr_engg       VARCHAR(20),   -- 미결제약정
                              cntr_pred_pre   VARCHAR(20),   -- 미결제전일대비
                              theory_pric     VARCHAR(20),   -- 이론가
                              innr_vltl       VARCHAR(20),   -- 내재변동성
                              delta           VARCHAR(20),   -- 델타
                              gam             VARCHAR(20),   -- 감마
                              theta           VARCHAR(20),   -- 쎄타
                              vega            VARCHAR(20),   -- 베가
                              law             VARCHAR(20)    -- 로
);

CREATE TABLE schedule_info ( -- 스케줄 목록
                               id SERIAL PRIMARY KEY,
                               name VARCHAR(100) NOT NULL,       -- 작업 이름
                               schedule VARCHAR(50) NOT NULL,    -- 실행 스케줄 ("every 10s", "09:00" 등)
                               task_type VARCHAR(50) NOT NULL,   -- 실행할 작업 타입
                               enabled BOOLEAN DEFAULT true,     -- 활성화 여부
                               created_at TIMESTAMP DEFAULT now()
);
INSERT INTO public.schedule_info
(id, "name", schedule, task_type, enabled, created_at)
VALUES(1, 'trade-info', 'every 10s', 'GetTradeInfoLog', true, '2025-09-05 17:04:14.463');

INSERT INTO public.schedule_info
(id, "name", schedule, task_type, enabled, created_at)
VALUES(2, 'stock-info', 'every 3s', 'UpsertStockInfo', true, '2025-09-10 14:54:06.291');

CREATE TABLE tb_virtual_asset ( --가상자산 테이블
                                  asset_id         BIGSERIAL PRIMARY KEY,              -- 포지션 고유 ID

                                  user_id          BIGINT       NOT NULL,              -- 사용자 ID
                                  account_id       BIGINT       NOT NULL,              -- 계좌 ID (FK)

                                  stk_cd           VARCHAR(20)  NOT NULL,              -- 종목 코드
                                  market           VARCHAR(10)  NOT NULL,              -- 시장 구분(KOSPI/KOSDAQ 등)

                                  position_side    CHAR(1)      NOT NULL,              -- 'B'=매수(롱), 'S'=매도(숏)

                                  qty              NUMERIC(18,4) NOT NULL,             -- 총 보유 수량
                                  available_qty    NUMERIC(18,4) NOT NULL,             -- 주문 가능 수량

                                  avg_price        NUMERIC(18,2) NOT NULL,             -- 평균 매입단가
                                  last_price       NUMERIC(18,2),                      -- 최근 평가 기준 가격
                                  highest_price       NUMERIC(18,2),                   -- 매수 이후 최고가

                                  invested_amount  NUMERIC(18,2) NOT NULL,             -- 총 매입금액
                                  eval_amount      NUMERIC(18,2),                      -- 평가금액
                                  eval_pl          NUMERIC(18,2),                      -- 평가손익
                                  eval_pl_rate     NUMERIC(9,4),                       -- 평가수익률(%)

                                  today_buy_qty    NUMERIC(18,4) DEFAULT 0,            -- 당일 매수 총 수량
                                  today_sell_qty   NUMERIC(18,4) DEFAULT 0,            -- 당일 매도 총 수량

                                  status           VARCHAR(20)  NOT NULL DEFAULT 'ACTIVE', -- ACTIVE / CLOSED 등

                                  last_eval_at     TIMESTAMPTZ,                        -- 마지막 평가 시각
                                  created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                  updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- 한 계좌에서 같은 종목 + 같은 방향 포지션은 하나만 존재
-- CREATE UNIQUE INDEX uq_asset_account_stkcd_side
--     ON tb_virtual_asset (account_id, stk_cd, position_side);
CREATE UNIQUE INDEX IF NOT EXISTS uq_virtual_asset
    ON tb_virtual_asset (account_id, stk_cd, market, position_side);




CREATE TABLE tb_virtual_account ( -- 가상 계좌 테이블
                                    account_id       BIGSERIAL PRIMARY KEY,          -- 계좌 ID (PK)
                                    user_id          BIGINT     NOT NULL,            -- 사용자 ID

                                    account_name     VARCHAR(50) NOT NULL,           -- 계좌 이름 (예: "모의투자 계좌1")

                                    cash_balance     NUMERIC(18,2) NOT NULL,         -- 현재 사용 가능한 현금
                                    total_invested   NUMERIC(18,2) DEFAULT 0,        -- 총 투자금액 (포지션들의 invested_amount 합)
                                    total_eval       NUMERIC(18,2) DEFAULT 0,        -- 총 평가금액 (포지션들의 eval_amount 합)
                                    total_pl         NUMERIC(18,2) DEFAULT 0,        -- 전체 평가손익
                                    total_pl_rate    NUMERIC(9,4),                   -- 전체 수익률

                                    deposit_amount   NUMERIC(18,2) DEFAULT 0,        -- 총 입금액(시작 자금)
                                    withdraw_amount  NUMERIC(18,2) DEFAULT 0,        -- 총 출금액(있다면)

                                    status           VARCHAR(20) DEFAULT 'ACTIVE',   -- 계좌 상태

                                    created_at       TIMESTAMPTZ DEFAULT NOW(),
                                    updated_at       TIMESTAMPTZ DEFAULT NOW()
);
INSERT INTO tb_virtual_account VALUES(
       0, 0, 'test_account_1', 1000000, 0, 1000000, 0, 0, 1000000, 0, 'ACTIVE', now(), now()
);
-- 유저 하나가 여러 가상 계좌를 가질 수 있음 → 이름 중복 방지
CREATE UNIQUE INDEX uq_virtual_account_user_name
    ON tb_virtual_account (user_id, account_name);



CREATE TABLE tb_virtual_order ( -- 주문 테이블
                                  order_id        BIGSERIAL      PRIMARY KEY,            -- 주문 ID (PK)

                                  user_id         BIGINT         NOT NULL,               -- 사용자 ID
                                  account_id      BIGINT         NOT NULL,               -- 계좌 ID (tb_virtual_account FK 대상)

                                  stk_cd          VARCHAR(20)    NOT NULL,               -- 종목 코드
                                  market          VARCHAR(10)    NOT NULL,               -- 시장 구분 (KOSPI, KOSDAQ, ETF, US 등)

                                  side            CHAR(1)        NOT NULL,               -- 'B'=매수, 'S'=매도
                                  order_type      VARCHAR(20)    NOT NULL,               -- 'MARKET', 'LIMIT' 등
                                  time_in_force   VARCHAR(20),                           -- 'DAY', 'IOC', 'FOK' 등 (옵션)

                                  price           NUMERIC(18,2),                         -- 주문 가격 (시장가면 NULL 또는 0)
                                  qty             NUMERIC(18,4)  NOT NULL,               -- 주문 수량

                                  filled_qty      NUMERIC(18,4)  NOT NULL DEFAULT 0,     -- 누적 체결 수량
                                  remaining_qty   NUMERIC(18,4)  NOT NULL,               -- 남은 수량 = qty - filled_qty

                                  status          VARCHAR(20)    NOT NULL,               -- 'NEW','OPEN','PARTIAL','FILLED','CANCELED','REJECTED'

                                  client_order_id VARCHAR(50),                           -- 클라이언트에서 부여하는 주문 ID (선택)
                                  reason          VARCHAR(255),                          -- 거절/취소 사유 등 (옵션)

                                  created_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW(), -- 주문 생성 시각
                                  updated_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW()  -- 마지막 상태 변경 시각
);
-- 계좌 + 클라이언트 주문 ID 기준으로 유니크하게 관리하고 싶을 때 사용
CREATE UNIQUE INDEX uq_virtual_order_client
    ON tb_virtual_order (account_id, client_order_id);


-- 체결 로그 테이블
CREATE TABLE tb_virtual_trade_log (
                                      trade_id        BIGSERIAL      PRIMARY KEY,            -- 체결 ID (PK)

                                      order_id        BIGINT         NOT NULL,               -- 어떤 주문의 체결인지 (FK 후보)
                                      user_id         BIGINT         NOT NULL,               -- 사용자 ID
                                      account_id      BIGINT         NOT NULL,               -- 계좌 ID
                                      stk_cd          VARCHAR(20)    NOT NULL,               -- 종목 코드
                                      market          VARCHAR(10)    NOT NULL,               -- 시장 구분

                                      side            CHAR(1)        NOT NULL,               -- 'B'=매수, 'S'=매도

                                      filled_qty      NUMERIC(18,4)  NOT NULL,               -- 이번 체결 수량
                                      filled_price    NUMERIC(18,2)  NOT NULL,               -- 이번 체결 단가
                                      filled_amount   NUMERIC(18,2)  NOT NULL,               -- 체결 금액 = filled_qty * filled_price

                                      fee_amount      NUMERIC(18,2)  DEFAULT 0,              -- 수수료 (옵션)
                                      tax_amount      NUMERIC(18,2)  DEFAULT 0,              -- 거래세 (옵션)

                                      created_at      TIMESTAMPTZ    NOT NULL DEFAULT NOW()  -- 체결 발생 시각
);

-- 주문별 체결 조회 자주 할 거라면 인덱스
CREATE INDEX idx_trade_log_order
    ON tb_virtual_trade_log (order_id);

-- 계좌, 종목 기준 체결 조회용 인덱스 (옵션)
CREATE INDEX idx_trade_log_account_stk
    ON tb_virtual_trade_log (account_id, stk_cd, created_at);



CREATE TABLE IF NOT EXISTS tb_stock_score (
    stk_cd            varchar(20) PRIMARY KEY,   -- 종목코드 (ex: 005930)

-- 총점 + 구성점수
    score_total       numeric(9,4)  NOT NULL,    -- 최종 점수 (0~100 권장)
    score_fundamental numeric(9,4)  DEFAULT 0,   -- 기업(재무) 점수
    score_momentum    numeric(9,4)  DEFAULT 0,   -- 모멘텀 점수
    score_market      numeric(9,4)  DEFAULT 0,   -- (선택) 시장/섹터 점수
    score_risk        numeric(9,4)  DEFAULT 0,   -- (감점) 리스크 점수(변동성 등)

-- 판단 근거로 남길 지표들(나중에 디버깅/튜닝할 때 필수)
    last_price        numeric(18,4),
    r1                numeric(18,8),             -- 최근 1틱 수익률(%)
    r2                numeric(18,8),             -- 최근 10틱 수익률(%)
    r3                numeric(18,8),             -- 최근 30틱 수익률(%)
    volatility        numeric(18,8),             -- stddev_pop(returns) (%)

-- 스코어 산정 시점 (이 시점 데이터로 계산했다)
    asof_tm           timestamptz NOT NULL DEFAULT now(),

    -- 설명/디버깅용 (어떤 규칙이 적용됐는지)
    meta              jsonb        DEFAULT '{}'::jsonb,

    created_at        timestamptz NOT NULL DEFAULT now(),
    updated_at        timestamptz NOT NULL DEFAULT now()
    );

-- 조회 성능용 인덱스 (점수 상위 N개 뽑을 때)
CREATE INDEX IF NOT EXISTS ix_stock_score_total
    ON tb_stock_score (score_total DESC);

-- 최신 업데이트 순으로 보기 좋게
CREATE INDEX IF NOT EXISTS ix_stock_score_updated
    ON tb_stock_score (updated_at DESC);


-- 종목별 변동성 조회 function
CREATE OR REPLACE FUNCTION fn_volatility(
  in_stk_cd text,
  in_minutes integer DEFAULT 300
)
RETURNS numeric
LANGUAGE sql
STABLE
AS $$
SELECT stddev_pop(ret)::numeric
FROM (
         SELECT
             ((cur_prc - LAG(cur_prc) OVER (ORDER BY tm))
                  / NULLIF(LAG(cur_prc) OVER (ORDER BY tm), 0) * 100) AS ret
         FROM trade_info_log
         WHERE stk_cd = in_stk_cd
           AND tm >= now() - make_interval(mins => in_minutes)
     ) t
WHERE ret IS NOT NULL;
$$;
