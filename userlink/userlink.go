// Package userlink contains the list of administrative user for the system, copied from a controller node
// All users have accesses, into a profile and even extended access based upon table records.
// It needs context xmodule.
package userlink

import (
	"fmt"
	"strconv"

	"github.com/webability-go/xdominion"
	"github.com/webability-go/xmodules/context"
)

func SynchroUsers(sitecontext *context.Context, fromcontext *context.Context) []string {

	msg := []string{}

	// load from origin context
	table := fromcontext.GetTable("user_user")
	if table == nil {
		return []string{"Error: the origin table user_user does not exist."}
	}
	totable := sitecontext.GetTable("user_user")
	if totable == nil {
		return []string{"Error: the destination table user_user does not exist."}
	}

	users, _ := table.SelectAll(nil, xdominion.XFieldSet{"key", "status", "name", "mail"})
	keys := map[int]bool{}
	if users != nil {
		for _, u := range *users {
			key, _ := u.GetInt("key")
			_, err := totable.Upsert(key, u)
			if err != nil {
				msg = append(msg, "Error adding user:"+fmt.Sprint(err))
			}
			keys[key] = true
		}
	}
	// extra users to delete
	localusers, _ := table.SelectAll(nil, xdominion.XFieldSet{"key"})
	localkeys := map[int]bool{}
	if localusers != nil {
		for _, u := range *localusers {
			key, _ := u.GetInt("key")
			localkeys[key] = true
		}
	}
	// diff localusers - users
	for k := range keys {
		if localkeys[k] {
			delete(localkeys, k)
		}
	}
	// set status to "deleted" to X (not available anymore), they may be used by the local code
	for k := range localkeys {
		totable.Update(k, xdominion.XRecord{"status": "X"})
	}

	cnt, _ := totable.Count(nil)
	msg = append(msg, strconv.Itoa(cnt)+" admin users synchronized")
	return msg
}
