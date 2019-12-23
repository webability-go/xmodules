# xmodules

Standard modules for Xamboo.

The released modules are actually:

- context: is used to contextualize a xamboo session (hit of server) to a specific set of pointers to caches, databases, logs, config files, etc. so you will use the correct set of data and tables on the correct site.

import "github.com/webability-go/xmodules/context"


- stat: tables to keep stats of use of anything, from site hits to IOT events.

import "github.com/webability-go/xmodules/stat"


- user: tables to keep all the administration users of the system, with complex profiles and access rights on the sessions

import "github.com/webability-go/xmodules/user"


- translation: a set of translation tables to keep translated words of anything, from database field to files. It supports all the known languages in UTF8.

import "github.com/webability-go/xmodules/translation"


- country: is the list of ISO official countries ready to use in a database.

import "github.com/webability-go/xmodules/country"


- metric: a full set of units to count things and convert between them.

import "github.com/webability-go/xmodules/metric"


- usda: the official USDA tables for nutrients to calculate recipes.

import "github.com/webability-go/xmodules/usda"


- ingredient: tables to manage ingredients for food and recipes.

import "github.com/webability-go/xmodules/ingredient"


- material: tables to manage materials for recipes, do it yourself, and any type of things you can build.

import "github.com/webability-go/xmodules/material"




Modules soon available: (working on them)
------------------------------------------


- client: is a set of clients that can connect to the system with basic metadata, social source and basic data

import "github.com/webability-go/xmodules/client"


- clientremote: is a table to link clients PK/FKs but from a distant client set of tables in another database/server

import "github.com/webability-go/xmodules/clientremote"


- clientsecurity: is a set of access tables to build a solid set of access rights, profiles, atomic rights, etc.

import "github.com/webability-go/xmodules/clientsecurity"


- clientp18n: personalization for clients, to add any type of connected data to the clients, from colors to navigation and AI resolutions.

import "github.com/webability-go/xmodules/clientp18n"
