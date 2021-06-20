package airtable

import (
	"github.com/mehanizm/airtable"
)

func GetTable(api string, dbID string, table string) *airtable.Table {
	c := airtable.NewClient(api)
	return c.GetTable(dbID, table)
}
