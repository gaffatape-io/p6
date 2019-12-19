package main

import (
	"context"
	"cloud.google.com/go/firestore"
	"flag"
	"k8s.io/klog"
	"github.com/gaffatape-io/p6"
)

var (
	firestoreProjectID = flag.String("firestore_project_id", "dev-p6", "Firestore project to use")
)

func main() {	
	klog.InitFlags(nil)
	klog.Info("p6 server starting")
	ctx := context.Background()
	fs, err := firestore.NewClient(ctx, *firestoreProjectID)
	if err != nil {
		klog.Fatalf("firestore.NewClient(%q) failed; err:%+v", *firestoreProjectID, err)
	}
	store := &crud.Store{fs}
	
	klog.Infof("firestore:%q connected", *firestoreProjectID)
	klog.Info(store)
	klog.Info("p6 server started")
}
