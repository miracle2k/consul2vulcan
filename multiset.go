/*
	Call it like this:

	./multiset backend1=serverA=http://a:80 backend1=serverB=http://b:80 backend2=serverC==http://c:80

	It will insert into vulcand these servers with $TTL.

	If a backend is not yet defined, the instruction is skipped.
	You need to add the backends externally.

	We don't fully create the backends partly for performance
	reasons, partly because I've had trouble with vulcand if for
	some reason a partial backend-structure exists, i.e. ./backend
	expired, but ./ or ./servers still exists. vulcand will not
	validate such a structure.
 */

package main

import "os"
import "log"
import "fmt"
import "strconv"
import "strings"
import "github.com/coreos/go-etcd/etcd"


const (
	etcdKeyNotFound       = 100
	etcdKeyAlreadyExists  = 105
	etcdEventIndexCleared = 401
)


func main() {
	verbose := os.Getenv("VERBOSE") == "1"
	machines := []string{os.Getenv("ETCD")}
	ttl, err := strconv.Atoi(os.Getenv("TTL"))
	if (err != nil) {
		if (verbose) { log.Printf("$TTL is not a valid value (or not set), using default: 30") }
		ttl = 30
	}

	args := [][]string{}
	for _, argument := range os.Args[1:] {
		split := strings.Split(argument, "=")
		if len(split) != 3 {
			log.Fatalf("%s is not valid, needs three parts", argument)
		}

		args = append(args, split)
	}

	if verbose { log.Printf("%+v", args) }

	client := etcd.NewClient(machines)
	for _, info := range args {
		backend := info[0]
		serverId := info[1]
		serverUrl := info[2]

		// Check if a backend is defined for this
		_, err := client.Get(fmt.Sprintf("vulcand/backends/%v/backend", backend), false, false)
		switch {
		case err == nil:
		case err.(*etcd.EtcdError).ErrorCode == etcdKeyNotFound:
			if verbose { log.Printf("The backend %v is not defined yet, skipping", backend) }
			continue
		default:
			log.Printf("Failed to fetch backend info for %v, skipping: %v", backend, err)
			continue
		}

		if _, err := client.Set(fmt.Sprintf("/vulcand/backends/%v/servers/%v", backend, serverId), fmt.Sprintf("{\"URL\": \"%v\"}", serverUrl), uint64(ttl)); err != nil {
	        log.Printf("Failed to write to etcd: %s", err)
	        break
		}
		if verbose { log.Printf("Set %v (%v) to backend %v", serverId, serverUrl, backend) }
	}
}
