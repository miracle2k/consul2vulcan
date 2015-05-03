package main

import "os"
import "log"
import "time"
//import "github.com/coreos/go-etcd/etcd"
import consulapi "github.com/hashicorp/consul/api"


const (
	// The amount of time to do a blocking query for
	defaultWaitTime = 10 * time.Second
	etcdTTL = 60
	etcdHeartbeatWait = 30
)


func main() {
	consulconfig := &consulapi.Config{
		Address: os.Getenv("CONSUL_ADDR"),
	}
    consul, _ := consulapi.NewClient(consulconfig)

	sync(consul)
}


func sync(consul *consulapi.Client) {
	// Start goroutine that will keep heartbeating all our backends
	// to etcd/vulcand.
	// ch := make(chan []*consulapi.CatalogNode, 1)
	// go writeToEtcd(ch);

	catalog := consul.Catalog()
	health := consul.Health()

	var LastIndex uint64
	LastIndex = 0

	for {
		options := &consulapi.QueryOptions{
			WaitTime:   defaultWaitTime,
			WaitIndex:  LastIndex,
		}
		nodes, qm, err := catalog.Services(options)
		if err != nil {
			log.Printf("Error getting nodes from consul: %s", err)
			time.Sleep(10 * time.Second)
			continue
		}

		// Handle returned index
		if qm.LastIndex == LastIndex {
			log.Printf("Consul response has no new data (index was the same)")
			continue
		}
		if qm.LastIndex < LastIndex {
			log.Printf("Consul response had a lower index, resetting")
			LastIndex = 0
			continue
		}
		LastIndex = qm.LastIndex

		log.Printf("Update from consul with index %d", LastIndex)

		// Now we need to fetch all healthy nodes for all services
		for serviceName, _ := range nodes {
			health.Service(serviceName, "", true, nil)
	    }

		// Tell etcd writer about new nodes
		//ch <- nodes
	}
}


// always running goroutine, keeps changing just the nodes
// func writeToEtcd(ch chan []*consulapi.CatalogNode) {
// 	var nodes []*consulapi.CatalogNode
// 	nodesChanged := false

// 	for {
// 		// Coalesce multiple quick changes into a single etcd write
// 		var timeout time.Duration
// 		if (nodesChanged) {
// 			timeout = 1
// 		} else {
// 			timeout = etcdHeartbeatWait;
// 		}

// 		// Either wait for heartbeat interval, or a node update
// 		select {
// 		case nodes = <-ch:
// 			log.Printf("Etcd writer received new nodes")
// 			nodes = nodes

// 			// Update soon, but not right away
// 			nodesChanged = true
// 			continue

// 		case <-time.After(time.Second * timeout):
// 	    }

// 	    if (nodes == nil) {
// 	    	log.Printf("Don't have any nodes yet to write to etcd")
// 	    	continue
// 	    }

// 	    // Sync all nodes to etcd
// 	    log.Printf("Writing nodes to etcd")
// 		machines := []string{os.Getenv("ETCD")}
// 	    client := etcd.NewClient(machines)

// 	    for _, node := range nodes {
// 	    	log.Printf(node.Address)
// 	    	for _, service := range node.Services {
// 	    		log.Printf("----")
// 	    		log.Printf(service.ID)
// 		    	log.Printf(service.Service)
// 		    	log.Printf(service.Port)
// 		    	log.Printf(service.Address)
// 	    	}
// 	    	log.Printf("")
// 	    	log.Printf("")


// 			if _, err := client.Set("/vulcand/backends/b1/backend", "{\"Type\": \"http\"}", etcdTTL); err != nil {
// 		        log.Printf("Failed to write to etcd: %s", err)
// 		        break
// 		    }

// 		    if _, err := client.Set("/vulcand/backends/b1/servers/srv1", "{\"URL\": \"http://localhost:5000\"}", etcdTTL); err != nil {
// 		        log.Printf("Failed to write to etcd: %s", err)
// 		        break
// 		    }
// 		}

// 		nodesChanged = false
// 	}
// }