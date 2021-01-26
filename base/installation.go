package base

import (
	"errors"

	serverassets "github.com/webability-go/xamboo/assets"

	"github.com/webability-go/xmodules/tools"
)

func VerifyNeeds(ds serverassets.Datasource, needs []string) (bool, []string) {

	result := []string{}
	flag := true

	for _, need := range needs {
		// Needed modules: context and translation
		vc := ModuleInstalledVersion(ds, need)
		if vc == "" {
			result = append(result, tools.Message(messages, "moduleneeded", MODULEID, need))
			flag = false
		} else {
			result = append(result, tools.Message(messages, "moduleok", MODULEID, need))
		}
	}

	return flag, []string{}
}

func SynchroTable(ds serverassets.Datasource, tablename string) (error, []string) {

	result := []string{}
	result = append(result, tools.Message(messages, "analyze", tablename))

	table := ds.GetTable(tablename)
	if table == nil {
		result = append(result, tools.Message(messages, "notable", tablename))
		return errors.New(tools.Message(messages, "notable", tablename)), result
	}

	_, err := table.Count(nil) // num
	if err != nil {            //|| num == 0
		//		if err != nil {
		result = append(result, tools.Message(messages, "tablenoexist", tablename, err))
		//		}
		err1 := table.Synchronize()
		if err1 != nil {
			result = append(result, tools.Message(messages, "tableerror", tablename, err1))
			return err1, result // we stop HERE, error creating the table
		}
		result = append(result, tools.Message(messages, "tablecreated", tablename))
	} else {
		result = append(result, tools.Message(messages, "tablenotmodified", tablename))
	}
	return nil, result
}
