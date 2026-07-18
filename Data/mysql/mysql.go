package mysql // 假设你的包名

import (
	"fmt"
	"mint/errs"
	"mint/models"
	"net/http"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var db *gorm.DB

// Init 初始化 GORM MySQL 连接
func Init() error {
	// 1. 拼接 DSN
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		viper.GetString("mysql.user"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetString("mysql.port"),
		viper.GetString("mysql.dbname"),
	)

	// 2. 配置 GORM 生产级参数
	gormConfig := &gorm.Config{
		//	2.1 开启预编译 SQL 语句缓存
		PrepareStmt: true,

		//	2.2 配置日志级别：
		// 		生产环境请设为 Warn 或 Error。如果用默认的 Info，GORM 会把所有执行成功的 SQL 打印出来，严重消耗性能和磁盘。
		// 		进阶做法：你可以写一个自定义 Logger 接入上一节实现的 Zap 日志库。这里先用 GORM 自带的做掩饰。
		Logger: logger.Default.LogMode(logger.Warn),
	}

	// 3. 建立连接
	var err error
	db, err = gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		fmt.Printf("GORM connect DB failed, err:%v\n", err)
		return err
	}

	// 4. 获取底层的 *sql.DB 对象以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Printf("get sql.db failed, err:%v\n", err)
		return err
	}

	// 5. 生产级连接池配置
	// 设置打开数据库连接的最大数量，防止瞬间突发流量把 MySQL 服务端连接数打满
	sqlDB.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))

	// 设置空闲连接池中连接的最大数量，减少频繁建立/销毁连接的开销
	sqlDB.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))

	// 设置连接可复用的最大时间
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}

// Close 关闭 GORM MySQL 连接
func Close() {
	// 1. 获取底层原生的 sql.DB 对象
	sqlDB, err := db.DB()
	if err != nil {
		zap.L().Error("获取数据库底层连接失败:", zap.Error(err))
	} else {
		// 2. 关闭整个连接池，断开与 MySQL 的所有 TCP 连接
		if err := sqlDB.Close(); err != nil {
			zap.L().Error("关闭数据库连接池失败:", zap.Error(err))
		} else {
			zap.L().Info("数据库连接池已安全关闭！")
		}
	}
}

// CheckUserExist 检查用户是否存在
func CheckUserExist(username string) (bool, error) {
	var exist bool
	// 1. 执行查询语句
	err := db.Raw("SELECT EXISTS(SELECT 1 FROM user WHERE username = ?)", username).Scan(&exist).Error

	// 2. 错误处理
	if err != nil {
		return false, err
	}

	// 3. 返回结果
	return exist, nil
}

func InsertUser(user *models.User) error {
	if err := db.Create(user).Error; err != nil {
		return errs.NewAppError(http.StatusInternalServerError, errs.CodeInsertDBError, "数据库插入用户失败", err)
	}
	return nil
}
