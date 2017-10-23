// Package oraclegen helps with generating SQL statements based on a given schema and additional parameters
package oraclegen

// FixLastInsertIDbug is a nasty hack to deal with some bugs I found in rana's
// ora.v4 oracle driver. FIXME Add more information here.
func FixLastInsertIDbug() bool {
	return true
}
