package adminmenu

import (
	"golang.org/x/text/language"

	"github.com/webability-go/xcore/v2"
)

// Do no forget to call tools.BuildMessages from the init
var messages *map[language.Tag]*xcore.XLanguage

var smessages = map[language.Tag]map[string]string{
	language.English: {
		// Module installation messages
		// init.go
		"MODULENAME":     "Administration menu",
		"modulemodified": "The entry %s was modified successfully in the base_module table.",
		"commit":         "Installation successfull.",
		"rollback":       "Installation aborted with error: %s",

		// util.go
		"MAINMENU":                     "Default Administration menu",
		"accessgroup.name":             "Access group for menu administration",
		"accessgroup.description":      "Group of all accesses for menu administration",
		"access.name":                  "Access for menu administration",
		"access.description":           "Access for menu administration",
		"menufolder.name":              "Menu constructor",
		"menufolder.description1":      "Click on this line to see the different options for the administration of menus",
		"adminmenugroup.name":          "Menu groups",
		"adminmenugroup.description1":  "Menu groups",
		"adminmenuoption.name":         "Menu options",
		"adminmenuoption.description1": "Menu options",

		"moduleerror": "Error modifying the entry %s in the base_module table: %s",
		// Datasources transactions
		"transaction.exist":        "Error creating a transaction: There is already a started transaction.",
		"transaction.none":         "Error searching the transaction: There is no available transaction.",
		"transaction.commitnone":   "Error searching the transaction to commit: There is no available transaction.",
		"transaction.rollbacknone": "Error searching the transaction to rollback: There is no available transaction.",
		"transaction.error":        "Error in the transaction: %s",
		// Containers
		"database.none": "There is no available database in the datasource",
	},
	language.Spanish: {
		// Module installation messages
		// init.go
		"MODULENAME":     "Menu de administración",
		"modulemodified": "La entrada %s fue modificada con exito en la tabla base_module.",
		"commit":         "Instalación exitosa.",
		"rollback":       "Instalación con error: %s",

		// util.go
		"MAINMENU":                     "Menú de administración por defecto",
		"accessgroup.name":             "Grupos de accesos de la administración de menús",
		"accessgroup.description":      "Grupos de accesos de la administración de menús",
		"access.name":                  "Acceso de la administración de menús",
		"access.description":           "Acceso de la administración de menús",
		"menufolder.name":              "Constructor de menús",
		"menufolder.description1":      "Haz clic sobre esta linea para ver las diferentes opciones para la administración de menús.",
		"adminmenugroup.name":          "Grupos de menús",
		"adminmenugroup.description1":  "Grupos de menús",
		"adminmenuoption.name":         "Opciones de menús",
		"adminmenuoption.description1": "Opciones de menús",

		"moduleerror": "Error modificando la entrada %s en la tabla base_module: %s",
		// Datasources transactions
		"transaction.exist":        "Error creating a transaction: There is already a started transaction.",
		"transaction.none":         "Error searching the transaction: There is no available transaction.",
		"transaction.commitnone":   "Error searching the transaction to commit: There is no available transaction.",
		"transaction.rollbacknone": "Error searching the transaction to rollback: There is no available transaction.",
		"transaction.error":        "Error in the transaction: %s",
		// Containers
		"database.none": "There is no available database in the datasource",
	},
	language.French: {
		// Module installation messages
		// init.go
		"MODULENAME":     "Menu pour l'administration",
		"modulemodified": "L'entrée %s a été modifiée avec succès dans la table base_module.",
		"commit":         "Instalation réussie.",
		"rollback":       "Instalation avec erreur: %s",

		// util.go
		"MAINMENU":                     "Menu d'administration par défaut",
		"accessgroup.name":             "Groupe des accès d'administration de menus",
		"accessgroup.description":      "Groupe des accès d'administration de menus",
		"access.name":                  "Accès d'administration de menus",
		"access.description":           "Accès d'administration de menus",
		"menufolder.name":              "Constructeur de menus",
		"menufolder.description1":      "Constructeur de menus",
		"adminmenugroup.name":          "Groupes de menus",
		"adminmenugroup.description1":  "Groupes de menus",
		"adminmenuoption.name":         "Options de menus",
		"adminmenuoption.description1": "Options de menus",

		"moduleerror": "Erreur en modifiant l'entrée %s dans la table base_module: %s",
		// Datasources transactions
		"transaction.exist":        "Error creating a transaction: There is already a started transaction.",
		"transaction.none":         "Error searching the transaction: There is no available transaction.",
		"transaction.commitnone":   "Error searching the transaction to commit: There is no available transaction.",
		"transaction.rollbacknone": "Error searching the transaction to rollback: There is no available transaction.",
		"transaction.error":        "Error in the transaction: %s",
		// Containers
		"database.none": "There is no available database in the datasource",
	},
}
