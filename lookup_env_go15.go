// +build go1.5

package configr

import "os"

func init() {
	lookupEnv = os.LookupEnv
}
