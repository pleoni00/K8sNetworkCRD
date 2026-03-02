package main

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	contr "graph-controller/api"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	h "graph-controller/controller/handler"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(contr.AddToScheme(scheme))
}

func main() {
	ctx := context.Background()
	var (
		config *rest.Config
		err    error
	)

	kubeconfigFilePath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	if _, err := os.Stat(kubeconfigFilePath); errors.Is(err, os.ErrNotExist) {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigFilePath)
		if err != nil {
			panic(err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	informerFactory := informers.NewSharedInformerFactory(clientset, 30*time.Second)

	// Get Pod and Node informers from the factory
	podInformer := informerFactory.Core().V1().Pods().Informer()
	nodeInformer := informerFactory.Core().V1().Nodes().Informer()

	handler := h.NewHandler(podInformer, nodeInformer)

	podInformer.AddEventHandlerWithResyncPeriod(handler, 5*time.Minute)
	nodeInformer.AddEventHandlerWithResyncPeriod(handler, 5*time.Minute)

	stopCh := make(chan struct{})

	informerFactory.Start(stopCh)
	informerFactory.WaitForCacheSync(stopCh)

	// TODO: Aggiungi il bootstrap

	go func() {
		for {
			item, quit := handler.PodQueue.Get()
			if quit {
				return
			}

			err = handler.ProcessPod(item)
			if err != nil {
				handler.PodQueue.AddRateLimited(item)
			} else {
				handler.PodQueue.Forget(item)
			}
			handler.PodQueue.Done(item)
		}
	}()

	go func() {
		for {
			item, quit := handler.NodeQueue.Get()
			if quit {
				return
			}

			err = handler.ProcessNode(item)
			if err != nil {
				handler.NodeQueue.AddRateLimited(item)
			} else {
				handler.NodeQueue.Forget(item)
			}
			handler.NodeQueue.Done(item)
		}
	}()

	go func() {
		for {
			select {
			case <-time.After(7 * time.Minute):
				handler.PeriodicFunction()
			case <-ctx.Done():
				return
			}
		}
	}()

	<-ctx.Done()
	handler.PodQueue.ShutDown()
	handler.NodeQueue.ShutDown()
	close(stopCh)
}
