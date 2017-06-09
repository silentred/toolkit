package db

import (
	"container/ring"
	"io"
	"log"

	"fmt"

	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	"github.com/silentred/kassadin"
)

const (
	MaxIdle = 10
	MaxOpen = 20
)

type MysqlManager struct {
	Application  *kassadin.App `inject`
	Config       kassadin.MysqlConfig
	databases    map[string]*xorm.Engine
	readOnlyRing *ring.Ring
	master       *xorm.Engine
}

// NewMysqlManager returns a new MysqlManager
func NewMysqlManager(app *kassadin.App, config kassadin.MysqlConfig) (*MysqlManager, error) {
	var readOnlyLength int
	var err error
	var engine *xorm.Engine

	for _, instance := range config.Instances {
		if instance.ReadOnly {
			readOnlyLength++
		}
	}

	mm := &MysqlManager{
		Application:  app,
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
				fmt.Printf("WARNING when initializing MySQL: %s \n", err)
			}
		} else {
			engine, err = mm.newORM(instance)
			if err != nil {
				log.Fatal(err)
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

func (mm *MysqlManager) newORM(mysql kassadin.MysqlInstance) (*xorm.Engine, error) {
	var output io.Writer = os.Stdout

	orm, err := xorm.NewEngine("mysql", mysql.String())
	if err != nil {
		return nil, err
	}
	orm.SetMaxIdleConns(MaxIdle)
	orm.SetMaxOpenConns(MaxOpen)

	if mm.Application != nil {
		output = mm.Application.Logger("default").Output()
		// set Logger output
		logger := xorm.NewSimpleLogger(output)

		if mm.Application.Config.Mode == kassadin.ModeDev {
			orm.ShowSQL(true)
			orm.ShowExecTime(true)
			logger.ShowSQL(true)
			logger.SetLevel(core.LOG_DEBUG)
		} else {
			logger.SetLevel(core.LOG_ERR)
		}
		orm.SetLogger(logger)
	}

	err = orm.Ping()
	if err != nil {
		log.Fatal(err)
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
