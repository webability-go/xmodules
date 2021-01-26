package base

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
		"MODULENAME":     "XModules base",
		"modulemodified": "The entry %s was modified successfully in the base_module table.",
		"commit":         "Installation successfull.",
		"rollback":       "Installation aborted with error: %s",
		// installatin.go
		"moduleneeded":     "Verifying '%s': The module '%s' is not installed: ERROR.",
		"moduleok":         "Verifying '%s': The module '%s' is installed: PASSED.",
		"analyze":          "Analysing %s table.",
		"notable":          "Critical Error: the module 'base' table '%s' does not exist.",
		"tablenoexist":     "The table %s does not exist in the database: %s",
		"tableerror":       "The table %s was not created: %s",
		"tablecreated":     "The table %s was created (again).",
		"tablenotmodified": "The table %s was not created because it contains data.",

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
		"MODULENAME":     "base XModules",
		"modulemodified": "La entrada %s fue modificada con exito en la tabla base_module.",
		"commit":         "Instalación exitosa.",
		"rollback":       "Instalación con error: %s",
		// installatin.go
		"moduleneeded":     "Verificando '%s': El módulo '%s' no esta instalado: ERROR.",
		"moduleok":         "Verificando '%s': El módulo '%s' esta instalado: OK.",
		"analyze":          "Analizando la tabla %s.",
		"notable":          "Error crítico: la tabla del módulo 'base', '%s' no existe.",
		"tablenoexist":     "La tabla %s no existe en base de datos: %s",
		"tableerror":       "La tabla %s no pudo ser creada: %s",
		"tablecreated":     "La tabla %s fue creada (de nuevo).",
		"tablenotmodified": "La tabla %s no fue creada porque ya existe y contiene datos.",

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
		"MODULENAME":     "base XModules",
		"modulemodified": "L'entrée %s a été modifiée avec succès dans la table base_module.",
		"commit":         "Instalation réussie.",
		"rollback":       "Instalation avec erreur: %s",
		// installatin.go
		"moduleneeded":     "Verification '%s': Le module '%s' n'est pas installé: ERREUR.",
		"moduleok":         "Verification '%s': Le module '%s' est installé: OK.",
		"analyze":          "Analyze de la table %s.",
		"notable":          "Erreur critique: la table du module 'base', '%s' n'existe pas.",
		"tablenoexist":     "La table %s n'existe pas dans la base de données: %s",
		"tableerror":       "La table %s pe peux pas être créée: %s",
		"tablecreated":     "La table %s a été créée (de nouveau).",
		"tablenotmodified": "La table %s n'a pas été créée car elle existe déjà et contient des données.",

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
