// This file is generated. Do not edit.

package urlinterfaces

import "fmt"

type URLSearchParams interface {
	fmt.Stringer
	Size() int
	Append(string, string)
	Delete(string, string)
	Get(string) string
	GetAll(string) []string
	Has(string, string) bool
	Set(string, string)
	Sort()
}
