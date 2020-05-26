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

	"github.com/webability-go/xamboo/assets"
)

// Datasource is a portable structure containing pointer to usefull structures used in any datasource of sites. It must be compilant with assets.Datasource interface
// Since it's thread safe and based on maps and slices, it must be accessed through Get/Set functions with mutexes
// to avoid race conditions
// The is only ONE database by datasource, with a set of modules and tables into this database.
type Datasource struct {
	// The name of the datasource (informative only)
	Name string
	// A configuration for the datasource: (does not need lock to be accessed since it's a pointer)
	Config *xconfig.XConfig
	// Only one database per datasource
	database *xdominion.XBase
	// Languages knows by the datasource
	mlanguages sync.RWMutex
	languages  []language.Tag
	// A list of loggers for the datasource:
	mlogs sync.RWMutex
	logs  map[string]*log.Logger
	// A list of tables for the datasource:
	mtables sync.RWMutex
	tables  map[string]*xdominion.XTable
	// A list of in-memory caches for the datasource:
	mcaches sync.RWMutex
	caches  map[string]*xcore.XCache
	// A list of linked modules id => code version
	mmodules sync.RWMutex
	modules  map[string]string
}

func (ds *Datasource) GetName() string {
	return ds.Name
}

func (ds *Datasource) AddLanguage(lang language.Tag) {
	ds.mlanguages.Lock()
	ds.languages = append(ds.languages, lang)
	ds.mlanguages.Unlock()
}

func (ds *Datasource) GetLanguages() []language.Tag {
	ds.mlanguages.RLock()
	langs := make([]language.Tag, len(ds.languages))
	copy(langs, ds.languages)
	ds.mlanguages.RUnlock()
	return langs
}

func (ds *Datasource) SetLog(id string, logger *log.Logger) {
	ds.mlogs.Lock()
	ds.logs[id] = logger
	ds.mlogs.Unlock()
}

func (ds *Datasource) GetLog(id string) *log.Logger {
	ds.mlogs.RLock()
	l := ds.logs[id]
	ds.mlogs.RUnlock()
	return l
}

func (ds *Datasource) GetLogs() map[string]*log.Logger {
	ds.mlogs.RLock()
	logs := make(map[string]*log.Logger)
	for i, l := range ds.logs {
		logs[i] = l
	}
	ds.mlogs.RUnlock()
	return logs
}

func (ds *Datasource) Log(id string, messages ...interface{}) {
	ds.mlogs.RLock()
	l := ds.logs[id]
	if l == nil {
		l = ds.logs["main"]
	}
	ds.mlogs.RUnlock()
	l.Println(messages...)
}

func (ds *Datasource) SetDatabase(db *xdominion.XBase) {
	ds.database = db
}

func (ds *Datasource) GetDatabase() *xdominion.XBase {
	return ds.database
}

func (ds *Datasource) SetTable(id string, table *xdominion.XTable) {
	ds.mtables.Lock()
	ds.tables[id] = table
	ds.mtables.Unlock()
}

func (ds *Datasource) GetTable(id string) *xdominion.XTable {
	ds.mtables.RLock()
	t := ds.tables[id]
	ds.mtables.RUnlock()
	return t
}

func (ds *Datasource) GetTables() map[string]*xdominion.XTable {
	ds.mtables.RLock()
	tables := make(map[string]*xdominion.XTable)
	for i, t := range ds.tables {
		tables[i] = t
	}
	ds.mtables.RUnlock()
	return tables
}

func (ds *Datasource) SetCache(id string, cache *xcore.XCache) {
	ds.mcaches.Lock()
	ds.caches[id] = cache
	ds.mcaches.Unlock()
}

func (ds *Datasource) GetCache(id string) *xcore.XCache {
	ds.mcaches.RLock()
	c := ds.caches[id]
	ds.mcaches.RUnlock()
	return c
}

func (ds *Datasource) GetCaches() map[string]*xcore.XCache {
	ds.mcaches.RLock()
	caches := make(map[string]*xcore.XCache)
	for i, c := range ds.caches {
		caches[i] = c
	}
	ds.mcaches.RUnlock()
	return caches
}

func (ds *Datasource) SetModule(moduleid string, moduleversion string) {
	ds.mmodules.Lock()
	ds.modules[moduleid] = moduleversion
	ds.mmodules.Unlock()
}

func (ds *Datasource) GetModule(moduleid string) string {
	ds.mmodules.RLock()
	m := ds.modules[moduleid]
	ds.mmodules.RUnlock()
	return m
}

func (ds *Datasource) GetModules() map[string]string {
	ds.mmodules.RLock()
	modules := make(map[string]string)
	for i, v := range ds.modules {
		modules[i] = v
	}
	ds.mmodules.RUnlock()
	return modules
}

func (ds *Datasource) IsModuleAuthorized(id string) bool {
	return ds.GetModule(id) != ""
}

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
	datasources  map[string]assets.Datasource
	CoreLog      *log.Logger
}

func (cnt *Container) SetDatasource(id string, ds assets.Datasource) {
	cnt.mdatasources.Lock()
	cnt.datasources[id] = ds
	cnt.mdatasources.Unlock()
}

func (cnt *Container) GetDatasource(id string) assets.Datasource {
	cnt.mdatasources.RLock()
	ds := cnt.datasources[id]
	cnt.mdatasources.RUnlock()
	return ds
}

func (cnt *Container) GetDatasources() map[string]assets.Datasource {
	cnt.mdatasources.RLock()
	dss := make(map[string]assets.Datasource)
	for i, v := range cnt.datasources {
		dss[i] = v
	}
	cnt.mdatasources.RUnlock()
	return dss
}

// Createdatasource will create a new datasource, link databases and logs based on XConfig data
func (cnt *Container) CreateDatasource(name string, config *xconfig.XConfig) (assets.Datasource, error) {
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
		return nil, errors.New("There is no available database in the datasource")
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

// Createdatasource will create a new datasource, link databases and logs based on XConfig data
func (cnt *Container) TryDatasource(ctx *assets.Context, datasourcename string) assets.Datasource {

	dsn, _ := ctx.Sysparams.GetString(datasourcename)
	datasource := cnt.GetDatasource(dsn)
	if datasource != nil {
		return datasource
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
		datasources: map[string]assets.Datasource{},
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

// Analyze a datasource and gets back the main data
func GetDatasourceStats(ds *Datasource) *xcore.XDataset {

	subdata := xcore.XDataset{}
	subdata["languages"] = ds.GetLanguages()
	subdata["database"] = ds.GetDatabase()
	subdata["logs"] = ds.GetLogs()

	caches := []string{}
	for id := range ds.GetCaches() {
		caches = append(caches, id)
	}
	subdata["xcaches"] = caches

	tables := map[string]string{}
	for id, table := range ds.GetTables() {
		if table.Base != nil {
			db := table.Base.Database
			tables[id] = db
		} else {
			tables[id] = "N/A"
		}
	}
	subdata["tables"] = tables

	subdata["config"] = buildConfigSet(ds.Config)

	// analiza los módulos instalados
	modules := map[string]interface{}{}
	for id, v := range ds.GetModules() {
		md := struct {
			Version          string
			InstalledVersion string
		}{v, ModuleInstalledVersion(ds, id)}
		modules[id] = md
	}
	subdata["modules"] = modules

	return &subdata
}
