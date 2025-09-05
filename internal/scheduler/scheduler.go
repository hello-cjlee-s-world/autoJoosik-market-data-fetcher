package scheduler

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type TaskConfig struct {
	Name     string // 작업 이름
	Schedule string // cron 표현식 또는 "every 10s" 같은 형태
	TaskType string // 어떤 작업을 실행할지 구분
}

type ScheduleInfo struct {
	ID        int
	Name      string
	Schedule  string
	TaskType  string
	Enabled   bool
	CreatedAt time.Time
}

func LoadTaskConfigs(ctx context.Context, pool *pgxpool.Pool) []TaskConfig {
	rows, _ := pool.Query(ctx, "SELECT name, schedule, task_type FROM schedule_info")
	var configs []TaskConfig
	for rows.Next() {
		var name, schedule, taskType string
		rows.Scan(&name, &schedule, &taskType)
		configs = append(configs, TaskConfig{name, schedule, taskType})
	}
	return configs
}
