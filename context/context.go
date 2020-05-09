// package context is the controler for all the XModules of Xamboo and is required to build any other XModule in the system.
// It controls different contexts for different sites, installed xmodules, logs, caches, databases and tables.
package context

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
	// Only one database per context
	database *xdominion.XBase
	// Languages knows by the context
	mlanguages sync.RWMutex
	languages  []language.Tag
	// A list of loggers for the context:
	mlogs sync.RWMutex
	logs  map[string]*log.Logger
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

func (ctx *Context) SetDatabase(db *xdominion.XBase) {
	ctx.database = db
}

func (ctx *Context) GetDatabase() *xdominion.XBase {
	return ctx.database
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

func (cnt *Container) SetContext(id string, ctx *Context) {
	cnt.mcontexts.Lock()
	cnt.contexts[id] = ctx
	cnt.mcontexts.Unlock()
}

func (cnt *Container) GetContext(id string) *Context {
	cnt.mcontexts.RLock()
	ctx := cnt.contexts[id]
	cnt.mcontexts.RUnlock()
	return ctx
}

func (cnt *Container) GetContexts() map[string]*Context {
	cnt.mcontexts.RLock()
	ctxs := make(map[string]*Context)
	for i, v := range cnt.contexts {
		ctxs[i] = v
	}
	cnt.mcontexts.RUnlock()
	return ctxs
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
func (cnt *Container) CreateContext(name string, config *xconfig.XConfig) (*Context, error) {
	// Crear los contextos basados en el CoreConfig
	ctx := &Context{
		Name:    name,
		Config:  config,
		logs:    map[string]*log.Logger{},
		tables:  map[string]*xdominion.XTable{},
		caches:  map[string]*xcore.XCache{},
		modules: map[string]string{},
	}

	// fill context LOGS and DATABASES with the definition of Context Config. Caches and Tables depends on modules called
	database := config.GetConfig("database")
	if database == nil {
		// Missing Database
		return nil, errors.New("There is no available database in the context")
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
	ctx.SetDatabase(XBase)

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
		// If the module has been compile, we set it up
		xm := strings.Split(m, "|")
		modid := xm[0]
		modprefix := ""
		if len(xm) > 1 {
			modprefix = xm[1]
		}
		cmod := ModulesList.Get(modid)
		if cmod != nil {
			cmod.Setup(ctx, modprefix)
		}
	}

	cnt.SetContext(name, ctx)
	return ctx, nil
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

	cnt := &Container{
		contexts: map[string]*Context{},
		CoreLog:  CoreLog,
	}

	contexts, _ := CoreConfig.GetStringCollection("context")
	for _, context := range contexts {
		cfgpath, _ := CoreConfig.GetString(context + "+config")
		cfg := xconfig.New()
		cfg.LoadFile(cfgpath)
		_, err := cnt.CreateContext(context, cfg)
		if err != nil {
			fmt.Println("Error creating context:", err)
		}
	}
	return cnt
}

// Analyze a context and gets back the main data
func GetContextStats(ctx *Context) *xcore.XDataset {

	subdata := xcore.XDataset{}
	subdata["languages"] = ctx.GetLanguages()
	subdata["database"] = ctx.GetDatabase()
	subdata["logs"] = ctx.GetLogs()

	caches := []string{}
	for id := range ctx.GetCaches() {
		caches = append(caches, id)
	}
	subdata["xcaches"] = caches

	tables := map[string]string{}
	for id, table := range ctx.GetTables() {
		if table.Base != nil {
			db := table.Base.Database
			tables[id] = db
		} else {
			tables[id] = "N/A"
		}
	}
	subdata["tables"] = tables

	subdata["config"] = buildConfigSet(ctx.Config)

	// analiza los módulos instalados
	modules := map[string]interface{}{}
	for id, v := range ctx.GetModules() {
		md := struct {
			Version          string
			InstalledVersion string
		}{v, ModuleInstalledVersion(ctx, id)}
		modules[id] = md
	}
	subdata["modules"] = modules

	return &subdata
}
