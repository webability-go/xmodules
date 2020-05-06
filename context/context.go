// package context is the controler for all the XModules of Xamboo and is required to build any other XModule in the system.
// It controls different contexts for different sites, installed xmodules, logs, caches, databases and tables.
package context

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"golang.org/x/text/language"

	"github.com/webability-go/xconfig"
	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"
)

// Context is a portable structure containing pointer to usefull structures used in any context of sites
// Since it's thread safe and based on maps and slices, it must be accessed through Get/Set functions with mutexes
// to avoid race conditions
// The is only ONE database by context, with a set of modules and tables into this database.
type Context struct {
	// The name of the context (informative only)
	Name string
	// A configuration for the context: (does not need lock to be accessed since it's a pointer)
	Config *xconfig.XConfig
	// Languages knows by the context
	mlanguages sync.RWMutex
	languages  []language.Tag
	// A list of loggers for the context:
	mlogs sync.RWMutex
	logs  map[string]*log.Logger
	// A list of databases for the context:
	mdatabases sync.RWMutex
	databases  map[string]*xdominion.XBase
	// A list of tables for the context:
	mtables sync.RWMutex
	tables  map[string]*xdominion.XTable
	// A list of in-memory caches for the context:
	mcaches sync.RWMutex
	caches  map[string]*xcore.XCache
	// A list of linked modules id => code version
	mmodules sync.RWMutex
	modules  map[string]string
}

func (ctx *Context) AddLanguage(lang language.Tag) {
	ctx.mlanguages.Lock()
	ctx.languages = append(ctx.languages, lang)
	ctx.mlanguages.Unlock()
}

func (ctx *Context) GetLanguages() []language.Tag {
	ctx.mlanguages.RLock()
	langs := make([]language.Tag, len(ctx.languages))
	copy(langs, ctx.languages)
	ctx.mlanguages.RUnlock()
	return langs
}

func (ctx *Context) SetLog(id string, logger *log.Logger) {
	ctx.mlogs.Lock()
	ctx.logs[id] = logger
	ctx.mlogs.Unlock()
}

func (ctx *Context) GetLog(id string) *log.Logger {
	ctx.mlogs.RLock()
	l := ctx.logs[id]
	ctx.mlogs.RUnlock()
	return l
}

func (ctx *Context) GetLogs() map[string]*log.Logger {
	ctx.mlogs.RLock()
	logs := make(map[string]*log.Logger)
	for i, l := range ctx.logs {
		logs[i] = l
	}
	ctx.mlogs.RUnlock()
	return logs
}

func (ctx *Context) Log(id string, messages ...interface{}) {
	ctx.mlogs.RLock()
	l := ctx.logs[id]
	if l == nil {
		l = ctx.logs["main"]
	}
	ctx.mlogs.RUnlock()
	l.Println(messages...)
}

func (ctx *Context) SetDatabase(id string, db *xdominion.XBase) {
	ctx.mdatabases.Lock()
	ctx.databases[id] = db
	ctx.mdatabases.Unlock()
}

func (ctx *Context) GetDatabase(id string) *xdominion.XBase {
	ctx.mdatabases.RLock()
	d := ctx.databases[id]
	ctx.mdatabases.RUnlock()
	return d
}

func (ctx *Context) GetDatabases() map[string]*xdominion.XBase {
	ctx.mdatabases.RLock()
	dbs := make(map[string]*xdominion.XBase)
	for i, d := range ctx.databases {
		dbs[i] = d
	}
	ctx.mdatabases.RUnlock()
	return dbs
}

func (ctx *Context) SetTable(id string, table *xdominion.XTable) {
	ctx.mtables.Lock()
	ctx.tables[id] = table
	ctx.mtables.Unlock()
}

func (ctx *Context) GetTable(id string) *xdominion.XTable {
	ctx.mtables.RLock()
	t := ctx.tables[id]
	ctx.mtables.RUnlock()
	return t
}

func (ctx *Context) GetTables() map[string]*xdominion.XTable {
	ctx.mtables.RLock()
	tables := make(map[string]*xdominion.XTable)
	for i, t := range ctx.tables {
		tables[i] = t
	}
	ctx.mtables.RUnlock()
	return tables
}

func (ctx *Context) SetCache(id string, cache *xcore.XCache) {
	ctx.mcaches.Lock()
	ctx.caches[id] = cache
	ctx.mcaches.Unlock()
}

func (ctx *Context) GetCache(id string) *xcore.XCache {
	ctx.mcaches.RLock()
	c := ctx.caches[id]
	ctx.mcaches.RUnlock()
	return c
}

func (ctx *Context) GetCaches() map[string]*xcore.XCache {
	ctx.mcaches.RLock()
	caches := make(map[string]*xcore.XCache)
	for i, c := range ctx.caches {
		caches[i] = c
	}
	ctx.mcaches.RUnlock()
	return caches
}

func (ctx *Context) SetModule(moduleid string, moduleversion string) {
	ctx.mmodules.Lock()
	ctx.modules[moduleid] = moduleversion
	ctx.mmodules.Unlock()
}

func (ctx *Context) GetModule(moduleid string) string {
	ctx.mmodules.RLock()
	m := ctx.modules[moduleid]
	ctx.mmodules.RUnlock()
	return m
}

func (ctx *Context) GetModules() map[string]string {
	ctx.mmodules.RLock()
	modules := make(map[string]string)
	for i, v := range ctx.modules {
		modules[i] = v
	}
	ctx.mmodules.RUnlock()
	return modules
}

// Container if the list of created contexts
type Container struct {
	mcontexts sync.RWMutex
	contexts  map[string]*Context
	CoreLog   *log.Logger
}

func (cs *Container) SetContext(id string, context *Context) {
	cs.mcontexts.Lock()
	cs.contexts[id] = context
	cs.mcontexts.Unlock()
}

func (cs *Container) GetContext(id string) *Context {
	cs.mcontexts.RLock()
	c := cs.contexts[id]
	cs.mcontexts.RUnlock()
	return c
}

func (cs *Container) GetContexts() map[string]*Context {
	cs.mcontexts.RLock()
	contexts := make(map[string]*Context)
	for i, v := range cs.contexts {
		contexts[i] = v
	}
	cs.mcontexts.RUnlock()
	return contexts
}

// CreateContainer will create a new container for contexts  from am XConfig data
// The XConfig file must have this syntax:
//  context=[contextid1]
//  context=[contextid2]
//  context=[contextid3]
//  contextid1-config=[path-to-config-file]
//  contextid2-config=[path-to-config-file]
//  contextid3-config=[path-to-config-file]
func CreateContainer(contextconfig *xconfig.XConfig) *Container {

	cc := &Container{
		contexts: map[string]*Context{},
	}

	contexts, _ := contextconfig.GetStringCollection("context")
	for _, context := range contexts {
		cfgpath, _ := contextconfig.GetString(context + "+config")
		cfg := xconfig.New()
		cfg.LoadFile(cfgpath)
		cc.CreateContext(context, cfg)
	}
	return cc
}

// CreateContext will create a new context, link databases and logs based on XConfig data
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
func (cs *Container) CreateContext(name string, config *xconfig.XConfig) *Context {
	// Crear los contextos basados en el CoreConfig
	ctx := &Context{
		Name:      name,
		Config:    config,
		logs:      map[string]*log.Logger{},
		databases: map[string]*xdominion.XBase{},
		tables:    map[string]*xdominion.XTable{},
		caches:    map[string]*xcore.XCache{},
		modules:   map[string]string{},
	}

	// fill context LOGS and DATABASES with the definition of Context Config. Caches and Tables depends on modules called
	databases := config.GetConfig("database")
	if databases != nil {
		for dbcname := range databases.Parameters {
			database := databases.GetConfig(dbcname)
			if database != nil {
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
				ctx.SetDatabase(dbcname, XBase)
			}
		}
	}

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
					ctx.SetLog(logname, log.New(ioutil.Discard, "", 0))
				} else {
					// open/create file log
					logw, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
					if err != nil {
						fmt.Println("Error opening LOG file in context: ", filename, err)
						ctx.SetLog(logname, log.New(ioutil.Discard, "", 0))
					} else {
						ctx.SetLog(logname, log.New(logw, name+":"+logname+": ", log.Ldate|log.Ltime|log.Lshortfile))
					}
				}
			}
		}
	}
	// "main" log is MANDATORY
	if !hasmain {
		ctx.SetLog("main", log.New(ioutil.Discard, "", 0))
	}

	// languages
	languages, _ := config.GetStringCollection("languages")
	if languages != nil && len(languages) > 0 {
		for _, l := range languages {
			lt, err := language.Parse(l)
			if err == nil {
				ctx.AddLanguage(lt)
			} else {
				fmt.Println("Error parsing language tag:", name, l, err)
			}
		}
	}

	// modules
	modulelist, _ := config.GetStringCollection("module")
	for _, m := range modulelist {
		ctx.SetModule(m, "-")
	}

	cs.SetContext(name, ctx)
	return ctx
}

// Create will scan a full config file for Containers
// The XConfig file must have this syntax:
//  context=[contextid1]
//  context=[contextid2]
//  context=[contextid3]
//  contextid1-config=[path-to-config-file]
//  contextid2-config=[path-to-config-file]
//  contextid3-config=[path-to-config-file]
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
	CoreLog.Println("xmodules::context::Create: Starting Core Log")

	cc := &Container{
		contexts: map[string]*Context{},
		CoreLog:  CoreLog,
	}

	contexts, _ := CoreConfig.GetStringCollection("context")
	for _, context := range contexts {
		cfgpath, _ := CoreConfig.GetString(context + "+config")
		cfg := xconfig.New()
		cfg.LoadFile(cfgpath)
		cc.CreateContext(context, cfg)
	}
	return cc
}

// Analyze a context and gets back the main data
func GetContextStats(sitecontext *Context) *xcore.XDataset {

	subdata := xcore.XDataset{}
	subdata["languages"] = sitecontext.GetLanguages()
	subdata["databases"] = sitecontext.GetDatabases()
	subdata["logs"] = sitecontext.GetLogs()

	caches := []string{}
	for id := range sitecontext.GetCaches() {
		caches = append(caches, id)
	}
	subdata["xcaches"] = caches

	tables := map[string]string{}
	for id, table := range sitecontext.GetTables() {
		if table.Base != nil {
			db := table.Base.Database
			tables[id] = db
		} else {
			tables[id] = "N/A"
		}
	}
	subdata["tables"] = tables

	subdata["config"] = buildConfigSet(sitecontext.Config)

	// analiza los módulos instalados
	modules := map[string]interface{}{}
	for id, v := range sitecontext.GetModules() {
		md := struct {
			Version          string
			InstalledVersion string
		}{v, ModuleInstalledVersion(sitecontext, id)}
		modules[id] = md
	}
	subdata["modules"] = modules

	return &subdata
}
