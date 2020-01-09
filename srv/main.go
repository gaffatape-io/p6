package main

import (
	"cloud.google.com/go/firestore"
	"context"
	"flag"
	"github.com/gaffatape-io/p6/crud"
	"github.com/gaffatape-io/p6/fe"
	"github.com/gaffatape-io/p6/okrs"
	"github.com/gaffatape-io/p6/rest"
	"k8s.io/klog"
	"net/http"
)

var (
	firestoreProjectID = flag.String("firestore_project_id", "dev-p6", "Firestore project to use")
	ipPort             = flag.String("port", ":8081", "Server ip:port moniker")
	developmentMode    = flag.Bool("development_mode", false, "Set to true to enable development mode")
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

	serve := http.NewServeMux()

	api := rest.NewMux(store, &okrs.Objectives{
		Objectives: store,
		RunTx:      store.RunTx,
	})

	web := fe.NewMux(*developmentMode)

	serve.Handle("/", web)
	serve.Handle("/api", api)

	klog.Info("p6 server listen starting")
	err = http.ListenAndServe(*ipPort, serve)
	if err != nil {
		klog.Fatalf("http.ListenAndServe(%q) failed; err:%+v", *ipPort, err)
	}
}
