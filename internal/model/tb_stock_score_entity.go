package model

import (
	"autoJoosik-market-data-fetcher/pkg/logger"
	"encoding/json"
	"time"
)

// TbStockScoreEntity 체결정보 로그
type TbStockScoreEntity struct {
	StkCd            string    `db:"stk_cd" json:"stk_cd"`
	ScoreTotal       float64   `db:"score_total" json:"score_total"`
	ScoreFundamental float64   `db:"score_fundamental" json:"score_fundamental"`
	ScoreMomentum    float64   `db:"score_momentum" json:"score_momentum"`
	ScoreMarket      float64   `db:"score_market" json:"score_market"`
	ScoreRisk        float64   `db:"score_risk" json:"score_risk"`
	LastPrice        float64   `db:"last_price" json:"last_price"`
	R1               float64   `db:"r1" json:"r1"`
	R2               float64   `db:"r2" json:"r2"`
	R3               float64   `db:"r3" json:"r3"`
	Volatility       float64   `db:"volatility" json:"volatility"`
	AsofTm           time.Time `db:"asof_tm" json:"asof_tm"`
	Meta             string    `db:"meta" json:"meta"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

func ToTbStockScoreEntity(str string) TbStockScoreEntity {
	var entity TbStockScoreEntity
	err := json.Unmarshal([]byte(str), &entity)
	if err != nil {
		logger.Error("While doing ToTbStockScoreEntity :: ", err.Error())
	}

	return entity
}
