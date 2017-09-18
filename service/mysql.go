package service

import (
	"container/ring"
	"fmt"
	"io"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/silentred/toolkit/config"
)

var (
	// MySQLMaxIdle of connection
	MySQLMaxIdle = 10
	// MySQLMaxOpen of connection
	MySQLMaxOpen = 20
)

// MysqlManager for mysql connection
type MysqlManager struct {
	App          Application `inject:"app"`
	Config       config.MysqlConfig
	databases    map[string]*xorm.Engine
	readOnlyRing *ring.Ring
	master       *xorm.Engine
}

// NewMysqlManager returns a new MysqlManager
func NewMysqlManager(app Application, config config.MysqlConfig) (*MysqlManager, error) {
	var readOnlyLength int
	var err error
	var engine *xorm.Engine

	for _, instance := range config.Instances {
		if instance.ReadOnly {
			readOnlyLength++
		}
	}

	mm := &MysqlManager{
		App:          app,
		Config:       config,
		databases:    make(map[string]*xorm.Engine),
		readOnlyRing: ring.New(readOnlyLength),
	}

	for _, instance := range config.Instances {
		if instance.ReadOnly {
			engine, err = mm.newORM(instance)
			if err == nil {
				mm.readOnlyRing.Value = engine
				mm.readOnlyRing = mm.readOnlyRing.Next()
				mm.databases[instance.Name] = engine
			} else {
				log.Fatalf("WARNING when initializing MySQL: %s \n", err)
			}
		} else {
			engine, err = mm.newORM(instance)
			if err != nil {
				log.Fatalf("mysql instance: %+v, err: %v", instance, err)
			}
			mm.master = engine
			mm.databases[instance.Name] = mm.master
		}
	}

	if mm.master == nil {
		err = fmt.Errorf("Mysql master is nil")
		return mm, err
	}

	return mm, nil
}

func (mm *MysqlManager) newORM(mysql config.MysqlInstance) (*xorm.Engine, error) {
	writer := mm.App.DefaultLogger().Output()
	debug := mm.App.GetConfig().Mode == config.ModeDev
	return NewXormEngine(mysql, writer, MySQLMaxIdle, MySQLMaxOpen, debug, mm.Config.Ping)
}

func NewXormEngine(mysql config.MysqlInstance, logWriter io.Writer, idle, open int, debug, ping bool) (*xorm.Engine, error) {
	var output io.Writer = os.Stdout

	orm, err := xorm.NewEngine("mysql", mysql.String())
	if err != nil {
		return nil, err
	}
	orm.SetMaxIdleConns(idle)
	orm.SetMaxOpenConns(open)

	if logWriter != nil {
		output = logWriter
		// set Logger output
		logger := xorm.NewSimpleLogger(output)

		if debug {
			orm.ShowSQL(true)
			orm.ShowExecTime(true)
			logger.ShowSQL(true)
			logger.SetLevel(core.LOG_DEBUG)
		} else {
			logger.SetLevel(core.LOG_ERR)
		}
		orm.SetLogger(logger)
	}

	if ping {
		err = orm.Ping()
		if err != nil {
			return nil, err
		}
	}

	return orm, nil
}

// DB gets databases by name
func (mm *MysqlManager) DB(name string) *xorm.Engine {
	if engine, ok := mm.databases[name]; ok {
		return engine
	}
	return nil
}

// R get read-only mysql Engine
func (mm *MysqlManager) R() *xorm.Engine {
	if mm.readOnlyRing.Len() == 0 {
		return mm.master
	}
	if e, ok := mm.readOnlyRing.Value.(*xorm.Engine); ok {
		mm.readOnlyRing = mm.readOnlyRing.Next()
		return e
	}
	return nil
}

// W gets master mysql Engine
func (mm *MysqlManager) W() *xorm.Engine {
	return mm.master
}

func initMySQL(app Application) error {
	if app.GetConfig().Mysql.InitMySQL {
		mm, err := NewMysqlManager(app, app.GetConfig().Mysql)
		if err != nil {
			return err
		}
		app.Set("mysql", mm, nil)
	}
	return nil
}
