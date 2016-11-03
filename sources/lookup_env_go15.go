// +build go1.5

package sources

import "os"

func init() {
	lookupEnv = os.LookupEnv
}
