# xmodules

Standard modules for Xamboo.

The released modules are actually:

- context: is used to contextualize a xamboo session (hit of server) to a specific set of pointers to caches, databases, logs, config files, etc.

import "github.com/webability-go/xmodules/context"


- structure: is used to build a basic data structure to use in the system. The structure is generally a complex set of data build from many database tables, for example an invoice, a recipe, an accounting policy, a client, etc.

import "github.com/webability-go/xmodules/structure"


- usda: the official USDA tables for nutrients to calculate recipes.

import "github.com/webability-go/xmodules/usda"


- metrics: a full set of units to count things and convert between them.

import "github.com/webability-go/xmodules/metrics"


- translation: a set of translation tables to keep translated words of anything, from database field to files. It supports anl the known languages in UTF8.

import "github.com/webability-go/xmodules/translation"


- ingredient: tables to manage ingredients for food and recipes.

import "github.com/webability-go/xmodules/ingredient"


- material: tables to manage materials for recipes, do it yourself, and any type of things you can build.

import "github.com/webability-go/xmodules/material"


- stat: tables to keep stats of use of anything, from site hits to IOT events.

import "github.com/webability-go/xmodules/stat"


Modules soon available: (working on them)

- country: is the list of ISO official countries ready to use in a database.

import "github.com/webability-go/xmodules/country"


- client: is a set of clients that can connect to the system with basic metadata, social source and FindAllStringSubmatch

import "github.com/webability-go/xmodules/client"


- clientsecurity: is a set of access tables to build a solid set of access rights, profiles, atomic rights, etc.

import "github.com/webability-go/xmodules/clientsecurity"


- clientp18n: personalization for clients, to add any type of connected data to the clients, from colors to navigation and AI resolutions.

import "github.com/webability-go/xmodules/clientp18n"
