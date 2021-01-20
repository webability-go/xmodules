// package base is the controler for all the XModules of Xamboo and is required to build any other XModule in the system.
// It controls different datasources for different sites, installed xmodules, logs, caches, databases and tables.
package base

import (
	"errors"
	"log"
	"sync"

	"golang.org/x/text/language"

	"github.com/webability-go/xconfig"
	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"

	"github.com/webability-go/xamboo/assets"

	"github.com/webability-go/xmodules/tools"
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
	// If this is a cloned datasource shell
	cloned      bool
	transaction *xdominion.XTransaction
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

func (ds *Datasource) CloneShell() assets.Datasource {
	if ds.cloned {
		return ds
	}
	nds := &Datasource{
		Name:       ds.Name,
		Config:     ds.Config,
		database:   ds.database,
		mlanguages: ds.mlanguages,
		languages:  ds.languages,
		mlogs:      ds.mlogs,
		logs:       ds.logs,
		mtables:    ds.mtables,
		tables:     ds.tables,
		mcaches:    ds.mcaches,
		caches:     ds.caches,
		mmodules:   ds.mmodules,
		modules:    ds.modules,
		cloned:     true,
	}
	return nds
}

func (ds *Datasource) StartTransaction() (*xdominion.XTransaction, error) {
	if ds.transaction != nil {
		msgerror := tools.Message(messages, "transaction.exist")
		ds.Log("main", msgerror)
		return nil, errors.New(msgerror)
	}
	tr, err := ds.database.BeginTransaction()
	if err != nil {
		ds.Log("main", tools.Message(messages, "transaction.error", err))
		return nil, err
	}
	ds.transaction = tr
	return tr, nil
}

func (ds *Datasource) GetTransaction() *xdominion.XTransaction {
	return ds.transaction
}

func (ds *Datasource) Commit() error {
	if ds.transaction == nil {
		msgerror := tools.Message(messages, "transaction.commitnone")
		ds.Log("main", msgerror)
		return errors.New(msgerror)
	}
	err := ds.transaction.Commit()
	if err != nil {
		ds.Log("main", tools.Message(messages, "transaction.error", err))
		return err
	}
	ds.transaction = nil
	return nil
}

func (ds *Datasource) Rollback() error {
	if ds.transaction == nil {
		msgerror := tools.Message(messages, "transaction.rollbacknone")
		ds.Log("main", msgerror)
		return errors.New(msgerror)
	}
	err := ds.transaction.Rollback()
	if err != nil {
		ds.Log("main", tools.Message(messages, "transaction.error", err))
		return err
	}
	ds.transaction = nil
	return nil
}
