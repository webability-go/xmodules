package base

import (
	// "embed"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/webability-go/xamboo/applications"

	"github.com/webability-go/xmodules/tools"
)

func VerifyNeeds(ds applications.Datasource, needs []string) (bool, []string) {

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

	return flag, result
}

func SynchroTable(ds applications.Datasource, tablename string) (error, []string) {

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

func SynchroFiles(origin fs.FS, destination string) (error, []string) {

	result := []string{}

	err := fs.WalkDir(origin, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			result = append(result, "Error:")
			return err
		}

		// if subdirectory, do nothing, the fs already have all the files into the structure
		if d.IsDir() {
			return nil
		}

		// read file in buffer
		result = append(result, "Copy file: "+path+" to "+destination+path)
		data, err := fs.ReadFile(origin, path)
		if err != nil {
			return err
		}

		dir := filepath.Dir(path)
		if dir != "." {
			os.MkdirAll(destination+dir, 0700)
		}
		err = os.WriteFile(destination+path, data, 0644)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err, result
	}

	return nil, result
}
