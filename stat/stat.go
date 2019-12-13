package stat

import (
  "time"
  "strconv"

  "golang.org/x/text/language"

  "github.com/webability-go/xdominion"

  "xmodules/context"
)

func InitStat(ctx *context.Context, prefix string, databasename string) error {

  // open 12 tables for each file
  for i := 1; i < 13; i ++ {
    m := ""
    if i < 10 { m = "0" }
    m += strconv.Itoa(i)
    
    ctx.Tables[prefix + "stat_" + m] = stat_stat(prefix, m)
    ctx.Tables[prefix + "stat_" + m].SetBase(ctx.Databases[databasename])
    ctx.Tables[prefix + "stat_" + m].SetLanguage(language.Spanish)
  }

  return nil
}

func SynchronizeDatabase(ctx *context.Context, prefix string) {

  for i := 1; i < 13; i ++ {
    m := ""
    if i < 10 { m = "0" }
    m += strconv.Itoa(i)

    // alguna protección para saber si existe la tabla y no tronarla si tiene datos?
    // hacer un select count
    num, err := ctx.Tables[prefix + "stat_" + m].Count(nil)
    if err != nil || num == 0 {
      ctx.Logs["main"].Println("The table " + prefix + "stat_" + m + " was created (again)")
      ctx.Tables[prefix + "stat_" + m].Synchronize()
    } else {
      ctx.Logs["main"].Println("The table " + prefix + "stat_" + m + " was not created because it contains data")
    }
  }
}

func getMonth() string {
  currentTime := time.Now()
  return currentTime.Format("01")
}

func RegisterStat(ctx *context.Context, prefix string, data xdominion.XRecord ) {
  data.Set("clave", 0)
  _, err := ctx.Tables[prefix + "stat_" + getMonth()].Insert(data)
  if err != nil {
    ctx.Logs["main"].Println("Error insertando el log:", err)
  }
}

// TODO(phil) hacer las funciones RegisterMiss y RegisterSys



