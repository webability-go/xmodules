# xmodules

Standard modules for Xamboo external loadable application.
All modules must be thread-safe.

The released modules are actually:

- base: is used to control datasources and modules into each of ones. It contains a specific set of pointers to caches, databases, logs, config files, etc. so you will use the correct set of data and tables on the correct site.

import "github.com/webability-go/xmodules/base"


- tools: is used to give some basic tools for keys, passwords, uuid, md5.

import "github.com/webability-go/xmodules/tools"


- stat: tables to keep stats of use of anything, from site hits to IOT events.

import "github.com/webability-go/xmodules/stat"
needs: "github.com/webability-go/xmodules/base"


- user: tables to keep all the administration users of the system, with complex profiles and access rights on the sessions

import "github.com/webability-go/xmodules/user"
needs: "github.com/webability-go/xmodules/base"


- userlink: tables to keep basic data of administration users (primary key, status, name, email) from a distant controller node that have loaded complete xmodules/user

import "github.com/webability-go/xmodules/userlink"
needs: "github.com/webability-go/xmodules/base"


- translation: a set of translation tables to keep translated words of anything, from database field to files. It supports all the known languages in UTF8.

import "github.com/webability-go/xmodules/translation"
needs: "github.com/webability-go/xmodules/base"
needs: "github.com/webability-go/xmodules/user" or "github.com/webability-go/xmodules/userlink"


- country: is the list of ISO official countries ready to use in a database.

import "github.com/webability-go/xmodules/country"
needs: "github.com/webability-go/xmodules/base"
needs: "github.com/webability-go/xmodules/translation"


- metric: a full set of units to count things and convert between them.

import "github.com/webability-go/xmodules/metric"


- usda: the official USDA tables for nutrients to calculate recipes.

import "github.com/webability-go/xmodules/usda"


- ingredient: tables to manage ingredients for food and recipes.

import "github.com/webability-go/xmodules/ingredient"


- material: tables to manage materials for recipes, do it yourself, and any type of things you can build.

import "github.com/webability-go/xmodules/material"

v2021-01-25:
- base, user, adminmenu and useradmin support now transactions to setup the modules.
- Errors control and messages enhanced during the installation of the modules.
- Separation of basic installation functions into the xmodules/base/installation.go

v2021-01-20:
- All the modules: enhancement to meet the new main structures and mmodules definition for Xamboo (use of datasource interface, use of bridge and assets modules entries)

v2020-05-25:
- user: the main admin password is now md5 encrypted

v2020-05-25:
- Change on all modules to meet new modules standard for Xamboo 1.3 (Datasources, instead of Contexts (wrong nomination), Modules new function StartContext, assets.Datasource interface.)
- The Xamboo server now controls the standard interfaces for XModules, Applications, Datasources, etc.
- Context module renamed to Base

v0.0.7:
- Modules homologation and Contexts homologation for use with Xamboo

v0.0.3:
- Added tools xmodule for basic functions
- user xmodule enhanced to log-in/log-out and controls session of a user

v0.0.2:
- Now uses xcore/v2




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
