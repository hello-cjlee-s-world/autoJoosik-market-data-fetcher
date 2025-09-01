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
                            tm               VARCHAR(20),             -- 시간 (API 그대로 저장)
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

                            created_at       TIMESTAMPTZ DEFAULT NOW() -- 로그 적재 시각
);
CREATE INDEX idx_trade_info_log_tm ON trade_info_log (tm);
CREATE INDEX idx_trade_info_log_stextp ON trade_info_log (stex_tp);
CREATE INDEX idx_trade_info_log_created_at ON trade_info_log (created_at);


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
