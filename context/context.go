package context

import (
  "fmt"
  "os"
  "io/ioutil"
  "log"

  "golang.org/x/text/language"

  "github.com/webability-go/xcore"
  "github.com/webability-go/xconfig"
  "github.com/webability-go/xdominion"
)

// The context is a portable structre containing pointer to usefull structures used in any context of sites
type Context struct {
  // The name of the context (informative only)
  Name string
	// Languages knows by the context
	Languages []language.Tag
  // A configuration for the context:
  Config *xconfig.XConfig
  // A list of loggers for the context:
  Logs map[string]*log.Logger
  // A list of databases for the context:
  Databases map[string]*xdominion.XBase
  // A list of tables for the context:
  Tables map[string]*xdominion.XTable
  // A list of in-memory caches for the context:
  Caches map[string]*xcore.XCache
}
// The type for a list of contexts
type ContextContainer map[string]*Context

func CreateContextContainer(contextconfig *xconfig.XConfig) ContextContainer {

  cc := ContextContainer{}

  contexts, _ := contextconfig.GetStringCollection("context")
  for _, context := range contexts {
    cfgpath, _ := contextconfig.GetString(context + "+config")
    cfg := xconfig.New()
    cfg.LoadFile(cfgpath)
    cc.CreateContext(context, cfg)
  }
  return cc
}

func (cs ContextContainer)CreateContext(name string, config *xconfig.XConfig) *Context {
  // Crear los contextos basados en el CoreConfig
  ctx := &Context{
    Name: name,
    Config: config,
    Logs: map[string]*log.Logger{},
    Databases: map[string]*xdominion.XBase{},
    Tables: map[string]*xdominion.XTable{},
    Caches: map[string]*xcore.XCache{},
  }

  // fill context LOGS and DATABASES with the definition of Context Config. Caches and Tables depends on modules called
  databases := config.GetConfig("database")
  if databases != nil {
    for dbcname, _ := range databases.Parameters {
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
          DBType: dbtype,
          Username: username,
          Password: password,
          Database: dbname,
          Host: host,
          SSL: ssl,
        }
        XBase.Logon()
        ctx.Databases[dbcname] = XBase
      }
    }
  }

  hasmain := false
  logs := config.GetConfig("log")
  if logs != nil {
    for logname, _ := range logs.Parameters {
      if logname == "main" {
        hasmain = true
      }
      xlog := logs.GetConfig(logname)
      if xlog != nil {
        // create a connector to the database
        filename, _ := xlog.GetString("file")

        if filename == "" {
          ctx.Logs[logname] = log.New(ioutil.Discard, "", 0)
        } else {
          // open/create file log
          logw, err := os.OpenFile(filename, os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
          if err != nil {
            fmt.Println("Error opening LOG file in context: ", filename, err)
            ctx.Logs[logname] = log.New(ioutil.Discard, "", 0)
          } else {
            ctx.Logs[logname] = log.New(logw, name + ":" + logname + ": ", log.Ldate|log.Ltime|log.Lshortfile)
          }
        }
      }
    }
  }
  // "main" log is MANDATORY
  if !hasmain {
    ctx.Logs["main"] = log.New(ioutil.Discard, "", 0)
  }

	// languages
  languages, _ := config.GetStringCollection("languages")
  if languages != nil && len(languages) > 0 {
	  for _, l := range languages {
  	  lt, err := language.Parse(l)
			if err == nil {
  			ctx.Languages = append(ctx.Languages, lt)
			} else {
  			fmt.Println("Error parsing language tag:", name, l, err)
			}
	  }
	}

  cs[name] = ctx
  return ctx
}

func (cs ContextContainer)GetContext(name string) *Context {
  return cs[name]
}

func (c *Context)AddLog(name string, logger *log.Logger) error {
  c.Logs[name] = logger
  return nil
}

func (c *Context)GetLog(name string) *log.Logger {
  return c.Logs[name]
}

func (c *Context)AddCache(name string, cache *xcore.XCache) error {
  c.Caches[name] = cache
  return nil
}

func (c *Context)GetCache(name string) *xcore.XCache {
  return c.Caches[name]
}

func (c *Context)AddDatabase(name string, database *xdominion.XBase) error {
  c.Databases[name] = database
  return nil
}

func (c *Context)GetDatabase(name string) *xdominion.XBase {
  return c.Databases[name]
}

func (c *Context)AddTable(name string, table *xdominion.XTable) error {
  c.Tables[name] = table
  return nil
}

func (c *Context)GetTable(name string) *xdominion.XTable {
  return c.Tables[name]
}
