package main

import (
	"context"
	"log"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	ctx := context.Background()

	taintKey := "node.kubernetes.io/unreachable"

	deleteOlderThan, err := time.ParseDuration(os.Args[1])
	if err != nil {
		log.Fatalln("parse error", err)
	}

	old := metav1.NewTime(time.Now().Add(-1 * deleteOlderThan))
	log.Println("deleting nodes with taint", taintKey, "which are older than", old)

	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		log.Panicln("clientcmd", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicln("clientset", err)
	}

	nodeList, err := clientset.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		log.Panicln("node list", err)
	}

	for _, node := range nodeList.Items {
		if len(node.Spec.Taints) == 0 {
			log.Println("node has no taints", node.Name)
			continue
		}

		for _, taint := range node.Spec.Taints {
			if taint.Key != taintKey {
				log.Println("node ", node.Name, "does not have taint", taintKey)
				continue
			}
			if !taint.TimeAdded.Before(&old) {
				log.Println("node ", node.Name, "has taint, but ", taint.TimeAdded, "is newer than", old)
				continue
			}

			log.Println("deleting ", node.Name)
			clientset.CoreV1().Nodes().Delete(ctx, node.Name, metav1.DeleteOptions{})
			log.Println("deleted", node.Name)
		}
	}

	log.Println("done")
}
