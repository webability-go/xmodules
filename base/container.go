// package base is the controler for all the XModules of Xamboo and is required to build any other XModule in the system.
// It controls different datasources for different sites, installed xmodules, logs, caches, databases and tables.
package base

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"golang.org/x/text/language"

	"github.com/webability-go/xconfig"
	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xamboo/cms/context"

	"github.com/webability-go/xmodules/tools"
)

// Container if the list of created datasources
// The XConfig file must have this syntax:
//  database.[dbid].type=[dbtype]
//  database.[dbid].username=[dbusername]
//  database.[dbid].password=[dbpassword]
//  database.[dbid].database=[dbname]
//  database.[dbid].server=[dbserver]
//  database.[dbid].ssl=[dbsslflag]
//
//  log.[logid].file=[pathtologfile]
// every line can be repeated for each dbid or logid
type Container struct {
	mdatasources sync.RWMutex
	datasources  map[string]applications.Datasource
	CoreLog      *log.Logger
}

type ContainersList map[string]*Container

func (cntl *ContainersList) AddContainer(id string, cnt *Container) {
	(*cntl)[id] = cnt
}

func (cntl *ContainersList) GetContainer(id string) *Container {
	return (*cntl)[id]
}

func (cntl *ContainersList) RegisterModule(mod applications.Module) {
	for _, cnt := range *cntl {
		cnt.RegisterModule(mod)
	}
}

func (cnt *Container) SetDatasource(id string, ds applications.Datasource) {
	cnt.mdatasources.Lock()
	cnt.datasources[id] = ds
	cnt.mdatasources.Unlock()
}

func (cnt *Container) GetDatasource(id string) applications.Datasource {
	cnt.mdatasources.RLock()
	ds := cnt.datasources[id]
	cnt.mdatasources.RUnlock()
	return ds
}

func (cnt *Container) GetDatasources() map[string]applications.Datasource {
	cnt.mdatasources.RLock()
	dss := make(map[string]applications.Datasource)
	for i, v := range cnt.datasources {
		dss[i] = v
	}
	cnt.mdatasources.RUnlock()
	return dss
}

// Createdatasource will create a new datasource, link databases and logs based on XConfig data
func (cnt *Container) CreateDatasource(name string, config *xconfig.XConfig) (applications.Datasource, error) {
	// Crear los datasourceos basados en el CoreConfig
	ds := &Datasource{
		Name:    name,
		Config:  config,
		logs:    map[string]*log.Logger{},
		tables:  map[string]*xdominion.XTable{},
		caches:  map[string]*xcore.XCache{},
		modules: map[string]string{},
	}

	// fill datasource LOGS and DATABASES with the definition of datasource Config. Caches and Tables depends on modules called
	database := config.GetConfig("database")
	if database == nil {
		// Missing Database
		return nil, errors.New(tools.Message(messages, "database.none"))
	}
	// create a connector to the database
	dbtype, _ := database.GetString("type")
	username, _ := database.GetString("username")
	password, _ := database.GetString("password")
	dbname, _ := database.GetString("database")
	host, _ := database.GetString("server")
	ssl, _ := database.GetBool("ssl")

	XBase := &xdominion.XBase{
		DBType:   dbtype,
		Username: username,
		Password: password,
		Database: dbname,
		Host:     host,
		SSL:      ssl,
	}
	XBase.Logon()
	ds.SetDatabase(XBase)

	hasmain := false
	logs := config.GetConfig("log")
	if logs != nil {
		for logname := range logs.Parameters {
			if logname == "main" {
				hasmain = true
			}
			xlog := logs.GetConfig(logname)
			if xlog != nil {
				// create a connector to the database
				filename, _ := xlog.GetString("file")

				if filename == "" {
					ds.SetLog(logname, log.New(ioutil.Discard, "", 0))
				} else {
					// open/create file log
					logw, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
					if err != nil {
						fmt.Println("Error opening LOG file in datasource: ", filename, err)
						ds.SetLog(logname, log.New(ioutil.Discard, "", 0))
					} else {
						ds.SetLog(logname, log.New(logw, name+":"+logname+": ", log.Ldate|log.Ltime|log.Lshortfile))
					}
				}
			}
		}
	}
	// "main" log is MANDATORY
	if !hasmain {
		ds.SetLog("main", log.New(ioutil.Discard, "", 0))
	}

	// languages
	languages, _ := config.GetStringCollection("languages")
	if languages != nil && len(languages) > 0 {
		for _, l := range languages {
			lt, err := language.Parse(l)
			if err == nil {
				ds.AddLanguage(lt)
			} else {
				fmt.Println("Error parsing language tag:", name, l, err)
			}
		}
	}

	// modules
	modulelist, _ := config.GetStringCollection("module")
	for _, m := range modulelist {
		ds.SetModule(m, "-")
		// If the module has been compile, we set it up
		xm := strings.Split(m, "|")
		modid := xm[0]
		modprefix := ""
		if len(xm) > 1 {
			modprefix = xm[1]
		}
		cmod := ModulesList.Get(modid)
		if cmod != nil {
			cmod.Setup(ds, modprefix)
		}
	}

	cnt.SetDatasource(name, ds)
	return ds, nil
}

func (cnt *Container) RegisterModule(mod applications.Module) {
	for _, ds := range cnt.datasources {
		ds.RegisterModule(mod)
	}
}

// TryDatasource will create a new datasource, link databases and logs based on XConfig data
func (cnt *Container) TryDatasource(ctx *context.Context, datasourcename string) applications.Datasource {

	var dsn string
	var datasource applications.Datasource
	if datasourcename != "" {
		dsn, _ = ctx.Sysparams.GetString(datasourcename)
		datasource = cnt.GetDatasource(dsn)
		if datasource != nil {
			return datasource
		}
	}

	dsn, _ = ctx.Sysparams.GetString("datasource")
	datasource = cnt.GetDatasource(dsn)
	if datasource != nil {
		return datasource
	}

	return nil
}

// Create will scan a full config file for Containers
// The XConfig file must have this syntax:
//  datasource=[datasourceid1]
//  datasource=[datasourceid2]
//  datasource=[datasourceid3]
//  datasourceid1-config=[path-to-config-file]
//  datasourceid2-config=[path-to-config-file]
//  datasourceid3-config=[path-to-config-file]
func Create(configfile string) *Container {

	CoreConfig := xconfig.New()
	CoreConfig.LoadFile(configfile)

	// Abrir CoreLog de Base
	logstr, _ := CoreConfig.GetString("logcore")
	logw, err := os.OpenFile(logstr, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening core log file xmodules::app::init:", err)
	}
	CoreLog := log.New(logw, "Core: ", log.Ldate|log.Ltime|log.Lshortfile)
	CoreLog.Println("xmodules::base::Create: Starting Core Log")

	cnt := &Container{
		datasources: map[string]applications.Datasource{},
		CoreLog:     CoreLog,
	}

	datasources, _ := CoreConfig.GetStringCollection("datasource")
	for _, datasourcename := range datasources {
		cfgpath, _ := CoreConfig.GetString(datasourcename + "+config")
		cfg := xconfig.New()
		cfg.LoadFile(cfgpath)
		_, err := cnt.CreateDatasource(datasourcename, cfg)
		if err != nil {
			fmt.Println("Error creating datasource:", err)
		}
	}
	return cnt
}

func TryDatasource(ctx *context.Context, datasourcename string) applications.Datasource {

	dscn, _ := ctx.Sysparams.GetString("datasourcecontainername")
	cnt := Containers.GetContainer(dscn)
	if cnt != nil {
		return cnt.TryDatasource(ctx, datasourcename)
	}
	return nil
}
