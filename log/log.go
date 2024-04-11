package log

import (
	"errors"
	"time"

	"github.com/webability-go/xamboo/applications"
	"github.com/webability-go/xdominion"
	//"github.com/webability-go/xmodules/base"
)

// GetLogs -- Get logs of user
func GetLogsList(ds applications.Datasource, cond *xdominion.XConditions, order *xdominion.XOrder, quantity int, first int) *xdominion.XRecords {

	log_log := ds.GetTable("log_log")
	if log_log == nil {
		ds.Log("xmodulesxamboo::log::GetLogs: Error, the log_log table is not available on this datasource")
		return nil
	}

	data, _ := log_log.SelectAll(cond, order, quantity, first)
	return data

}

// GetcCountLogs -- Count logs of user
func GetLogsCount(ds applications.Datasource, cond *xdominion.XConditions) int {

	log_log := ds.GetTable("log_log")
	if log_log == nil {
		ds.Log("xmodulesxamboo::log::GetCountLogs: Error, the log_log table is not available on this datasource")
		return 0
	}

	cnt, _ := log_log.Count(cond)

	return cnt
}

// AddLog -- Insert data on table
func AddLog(ds applications.Datasource, userid int, object, action, keyext string) error {

	log_log := ds.GetTable("log_log")
	if log_log == nil {
		ds.Log("xmodulesxamboo::log::AddLog: Error, the log_log table is not available on this datasource")
		return errors.New("No existe la tabla")
	}

	rec := &xdominion.XRecord{
		"key":       0,
		"user":      userid,
		"object":    object,
		"action":    action,
		"timestamp": time.Now(),
		"externkey": keyext,
	}

	_, err := log_log.Insert(rec)

	return err
}
