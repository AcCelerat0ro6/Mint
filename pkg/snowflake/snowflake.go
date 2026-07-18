package snowflake

import (
	"fmt"
	"time"

	sf "github.com/bwmarrin/snowflake"
	"github.com/spf13/viper"
)

var node *sf.Node

// Init 初始化雪花算法节点，需在 viper 已加载配置后调用
func Init() error {
	// 1. 解析起始时间，设置为雪花算法的 Epoch（毫秒级 Unix 时间戳）
	startTimeStr := viper.GetString("snowflake.start_time")
	st, err := time.Parse("2006-01-02", startTimeStr)
	if err != nil {
		return fmt.Errorf("parse snowflake.start_time failed: %w", err)
	}
	sf.Epoch = st.UnixMilli()

	// 2. 读取机器 ID 并创建节点
	machineID := viper.GetInt64("snowflake.machine_id")
	node, err = sf.NewNode(machineID)
	if err != nil {
		return fmt.Errorf("create snowflake node failed: %w", err)
	}
	return nil
}

// GenID 生成一个全局唯一的 int64 分布式 ID
func GenID() int64 {
	return node.Generate().Int64()
}
