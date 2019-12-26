package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"flag"
	"github.com/gaffatape-io/p6/crud"
	"github.com/gaffatape-io/p6/rest"
	"k8s.io/klog"
	"net/http"
)

var (
	firestoreProjectID = flag.String("firestore_project_id", "dev-p6", "Firestore project to use")
	ipPort             = flag.String("port", ":8081", "Server ip:port moniker")
)

func main() {
	flag.Parse()
	klog.InitFlags(nil)
	klog.Info("p6 server starting")
	ctx := context.Background()
	fs, err := firestore.NewClient(ctx, *firestoreProjectID)
	if err != nil {
		klog.Fatalf("firestore.NewClient(%q) failed; err:%+v", *firestoreProjectID, err)
	}
	store := &crud.Store{fs}
	klog.Infof("firestore:%q connected", *firestoreProjectID)

	api := rest.NewMux(store)
	klog.Info("p6 server listen starting")
	err = http.ListenAndServe(*ipPort, api)
	if err != nil {
		klog.Fatalf("http.ListenAndServe() failed; err:%+v", err)
	}
}
