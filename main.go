package main

/*
#cgo LDFLAGS: -lfdb_c -lm
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <pthread.h>
#define FDB_API_VERSION 21
#include <foundationdb/fdb_c.h>

extern void GoCallback(void* chanPointer);

// TODO -- clean these up (or get rid of them)
uint8_t* convert(char* test) {
	return (uint8_t*)(test);
}

char* convert2(uint8_t* test) {
	return (char*)(test);
}

void callbackHelper(FDBFuture* f, void* chanPointer) {
	fdb_future_set_callback(f, (FDBCallback)GoCallback, chanPointer);
}

*/
import "C"

import "fmt"
import "unsafe"

//func futureGoroutine(f, *C.FDBFuture, c chan bool) {
//
//}

func main() {
	C.fdb_select_api_version_impl(C.FDB_API_VERSION, C.FDB_API_VERSION)

	C.fdb_setup_network()

	// Not sure if we have to stop this goroutine later...
	go C.fdb_run_network()

	/* == Cluster == */
	clusterFuture := C.fdb_create_cluster(C.CString("/etc/foundationdb/fdb.cluster"))

	clusterFutureChan := make(chan bool, 1)
	C.callbackHelper(clusterFuture, unsafe.Pointer(&clusterFutureChan))
	<-clusterFutureChan

	var cluster *C.FDBCluster
	C.fdb_future_get_cluster(clusterFuture, &cluster)
	C.fdb_future_destroy(clusterFuture)

	/* == Database == */
	dbFuture := C.fdb_cluster_create_database(cluster, C.convert(C.CString("DB")), (C.int)(C.strlen(C.CString("DB"))))

	dbFutureChan := make(chan bool, 1)
	C.callbackHelper(dbFuture, unsafe.Pointer(&dbFutureChan))
	<-dbFutureChan

	var db *C.FDBDatabase
	C.fdb_future_get_database(dbFuture, &db)
	C.fdb_future_destroy(dbFuture)

	/* == Transaction == */
	var tr *C.FDBTransaction

	C.fdb_database_create_transaction(db, &tr)

	getFuture := C.fdb_transaction_get(tr, C.convert(C.CString("TestKey")), (C.int)(C.strlen(C.CString("TestKey"))), 0)

	getFutureChan := make(chan bool, 1)
	C.callbackHelper(getFuture, unsafe.Pointer(&getFutureChan))
	<-getFutureChan

	var valuePresent C.fdb_bool_t
	var value *C.uint8_t
	var valueLength C.int

	C.fdb_future_get_value(getFuture, &valuePresent, &value, &valueLength)

	C.fdb_future_destroy(getFuture)
	C.fdb_transaction_destroy(tr)

	val := C.GoStringN(C.convert2(value), valueLength)
	fmt.Println(val)
}
