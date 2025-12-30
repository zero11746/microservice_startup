package database

/*import (
	"database/sql"
	"entgo.io/ent/dialect"
	dialectsql "entgo.io/ent/dialect/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
	"user/config"
	"user/internal/ent"
)

var (
	mysqlClientConn *sql.DB   // 主库连接（写）
	slaveDBs        []*sql.DB // 从库连接（读）
	slaveEntDrvs    []*dialectsql.Driver
	entClient       *ent.Client // Ent 客户端实例
)

func InitMysqlConnect() (err error) {
	mysqlConfig := config.GetConfig().Mysql
	var (
		DefaultMaxConnectLifeTime = 10  // 最长连接生命周期（s）
		DefaultMaxIdleConnect     = 50  // 最大空闲连接数
		DefaultMaxOpenConnect     = 100 // 最大打开连接数
	)

	// 设置连接池（复用原有逻辑）
	setDBPool := func(db *sql.DB) {
		maxOpenConnect := DefaultMaxOpenConnect
		if mysqlConfig.MaxOpenCount > 0 {
			maxOpenConnect = mysqlConfig.MaxOpenCount
		}
		db.SetMaxOpenConns(maxOpenConnect)

		maxIdleConnect := DefaultMaxIdleConnect
		if mysqlConfig.MaxIdleCount > 0 {
			maxIdleConnect = mysqlConfig.MaxIdleCount
		}
		db.SetMaxIdleConns(maxIdleConnect)

		maxConnectLifeTime := time.Second * time.Duration(DefaultMaxConnectLifeTime)
		if mysqlConfig.ConnMaxLifeTime > 0 {
			maxConnectLifeTime = time.Second * time.Duration(mysqlConfig.ConnMaxLifeTime)
		}
		db.SetConnMaxLifetime(maxConnectLifeTime)
	}

	// 单库模式（无读写分离）
	if !mysqlConfig.Separation {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
			mysqlConfig.Master.Username,
			mysqlConfig.Master.Password,
			mysqlConfig.Master.Host,
			mysqlConfig.Master.Port,
			mysqlConfig.Master.DBName,
			mysqlConfig.DSNParams)

		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return fmt.Errorf("init single mysql connect failed: %v", err)
		}
		setDBPool(db)
		mysqlClientConn = db
		slaveDBs = []*sql.DB{}

		// 初始化 Ent 客户端（单库模式） )
		drv := dialectsql.OpenDB(dialect.MySQL, db)
		entClient = ent.NewClient(ent.Driver(drv), ent.Debug())
		return nil
	}

	// 读写分离模式：初始化主库（写）
	masterDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		mysqlConfig.Master.Username,
		mysqlConfig.Master.Password,
		mysqlConfig.Master.Host,
		mysqlConfig.Master.Port,
		mysqlConfig.Master.DBName,
		mysqlConfig.DSNParams)
	masterDB, err := sql.Open("mysql", masterDSN)
	if err != nil {
		return fmt.Errorf("init master mysql connect failed: %v", err)
	}
	setDBPool(masterDB)
	mysqlClientConn = masterDB

	// 读写分离模式：初始化从库（读）
	slaveDBs = []*sql.DB{}
	slaveEntDrvs = []*dialectsql.Driver{} // 初始化从库 Ent 驱动切片
	for i, slave := range mysqlConfig.Slaves {
		slaveDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
			slave.Username,
			slave.Password,
			slave.Host,
			slave.Port,
			slave.DBName,
			mysqlConfig.DSNParams)
		slaveDB, err := sql.Open("mysql", slaveDSN)
		if err != nil {
			return fmt.Errorf("init slave %d mysql connect failed: %v", i+1, err)
		}
		setDBPool(slaveDB)
		slaveDBs = append(slaveDBs, slaveDB)

		slaveEntDrv := dialectsql.OpenDB(dialect.MySQL, slaveDB)
		slaveEntDrvs = append(slaveEntDrvs, slaveEntDrv)
	}

	// 初始化 Ent 客户端（读写分离模式，默认连接主库）
	drv := dialectsql.OpenDB(dialect.MySQL, masterDB)
	entClient = ent.NewClient(ent.Driver(drv), ent.Debug())

	return nil
}*/
