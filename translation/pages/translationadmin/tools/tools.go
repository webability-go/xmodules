package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/webability-go/xamboo/cms/context"
	"github.com/webability-go/xcore/v2"
	"github.com/webability-go/xdominion"

	"github.com/webability-go/xmodules/base"
	"github.com/webability-go/xmodules/tools"
	"github.com/webability-go/xmodules/user/security"

	"xmodules/translation/assets"
)

var language *xcore.XLanguage

func Run(ctx *context.Context, template *xcore.XTemplate, xlanguage *xcore.XLanguage, e interface{}) interface{} {

	if language == nil {
		language = xlanguage
	}

	ok := security.Verify(ctx, security.USER, assets.ACCESS)
	if !ok {
		return ""
	}

	TRANSLATON, _ := ctx.Sysparams.GetString("translationdomains")
	TRANSLATONCONTAINER, _ := ctx.Sysparams.GetString("translationcontainer")
	digestkey, _ := ctx.Sysparams.GetString("clavedigest")
	digest := tools.GetTimeDigest(digestkey)

	params := &xcore.XDataset{
		"DIGEST":     digest,
		"SITE":       TRANSLATONCONTAINER,
		"TRANSLATON": TRANSLATON,
		"#":          language,
	}

	return template.Execute(params)
}

type lines map[string]string

// Tools

func Languagecode(ctx *context.Context, template *xcore.XTemplate, language *xcore.XLanguage, e interface{}) interface{} {

	entries := assets.GetEntries(ctx.LoggerError)
	ds := base.TryDatasource(ctx, assets.DATASOURCE)

	datascan := []string{"OK"}

	// Authorized languages and default language from the TABLE (not config, config points only "published ones")

	// Scan all the .language into the code:
	// The code directories are into the admin site configuration
	directories, _ := ctx.Sysparams.GetStringCollection("translationdirectory")
	for _, directory := range directories {
		// Search for AMY .language file into each dir, recursively
		datascan = append(datascan, "<b>Scanning: "+directory+"</b>")

		err := filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Si el archivo es un archivo regular, imprimir su ruta
			if info.Mode().IsRegular() && filepath.Ext(path) == ".language" {
				result, isok := analyse(path, []string{"es", "en", "pt", "fr"}, "es")
				if isok {
					// creates or verify entry into the database
					//	entry := translation.Get
					link := entries.GetThemeByName(ds, path)
					fmt.Println(link)

					if link == nil {
						link = &xdominion.XRecord{
							"clave":  0,
							"nombre": path,
							"fuente": 20,
						}
						entries.AddTheme(ds, link)
					}
					result += " <b style=\"color: green;\">Grabado en la tabla de temas</b><br />"
				}
				if result != "" {
					datascan = append(datascan, result)
				}
			}

			return nil
		})

		if err != nil {
			datascan = append(datascan, err.Error())
		}

	}

	//	result := bridge.JSONEncode(ctx, datascan)
	return datascan
}

// Returns a bool to TRUE if the language is ok to insert into themes to translate
func analyse(path string, languages []string, defaultlanguage string) (string, bool) {

	text := "<b style=\"color:blue;\">Analysing: " + path + "</b><br />"

	//	1. extract language
	if len(path) < 12 {
		return text + "Filename not good: " + path, false
	}
	lang := path[len(path)-12 : len(path)-8]
	dir, name := filepath.Split(path)

	type afile struct {
		name  string
		class int
		lang  *xcore.XLanguage
	}

	files := []afile{{name: path, class: 10}}
	if lang == "."+defaultlanguage+"." {
		// search for the other languages in the same path and analyse them
		for _, language := range languages {
			if language != defaultlanguage {
				newname := path[:len(path)-12] + "." + language + ".language"
				files = append(files, afile{name: newname, class: 10})
			}
		}
	} else {
		// search if default exists in same path
		// if not, error, missing default language
		newname := path[:len(path)-12] + "." + defaultlanguage + ".language"
		if _, err := os.Stat(newname); os.IsNotExist(err) {
			return text + "<b style=\"color: red;\">Error en el directorio " + dir + ": no existe el archivo del idioma por defecto " + name + "</b>", false
		}
		// Nothing to report, the default language will be analyzed apart
		return "", false
	}

	// Analyse all the files into files array
	var originlang *xcore.XLanguage
	for i, file := range files {
		lang, err := xcore.NewXLanguageFromXMLFile(file.name)
		if err != nil {
			files[i].class = 2
			if i == 0 {
				return text + "<b style=\"color: red;\">Error en el directorio " + dir + ": el archivo del idioma por defecto tiene un error " + file.name + "</b>", false
			}
		} else {
			files[i].lang = lang
			files[i].class = 1
			if i == 0 {
				originlang = lang
			}
		}
	}
	if originlang == nil {
		return text + "<b style=\"color: red;\">Error en el directorio " + dir + ": el archivo del idioma por defecto tiene un error </b>", false
	}

	// compare all the files against the default language
	original := originlang.GetEntries()
	for i, file := range files {
		if i == 0 {
			continue
		}
		if file.class != 1 {
			text += "<b style=\"color: red;\">File: " + file.name + " ERROR, file does not exists or have an error</b><br />"
			continue
		}

		numok, nummiss1, nummiss2 := comparemaps(original, file.lang.GetEntries())
		text += "File: " + file.name + " OK: " + fmt.Sprint(numok) + "; Missing in default: " + fmt.Sprint(nummiss1) + "; Missing in this: " + fmt.Sprint(nummiss2) + "<br />"
	}

	return text, true
}

func comparemaps(original, compare map[string]string) (int, int, int) {
	numok := 0
	nummiss1 := 0
	nummiss2 := 0

	for key := range original {
		if compare[key] == "" {
			nummiss1++
		} else {
			numok++
		}
	}
	for key := range compare {
		if original[key] == "" {
			nummiss2++
		}
	}

	return numok, nummiss1, nummiss2
}
