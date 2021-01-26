package useradmin

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
		"MODULENAME":     "Users Administration tools",
		"modulemodified": "The entry %s was modified successfully in the base_module table.",
		"commit":         "Installation successfull.",
		"rollback":       "Installation aborted with error: %s",

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
		"MODULENAME":     "Herramientas para administrar usuarios",
		"modulemodified": "La entrada %s fue modificada con exito en la tabla base_module.",
		"commit":         "Instalación exitosa.",
		"rollback":       "Instalación con error: %s",

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
		"MODULENAME":     "Outils pour administrer les utilisateurs",
		"modulemodified": "L'entrée %s a été modifiée avec succès dans la table base_module.",
		"commit":         "Instalation réussie.",
		"rollback":       "Instalation avec erreur: %s",

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
