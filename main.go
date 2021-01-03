package main

import (
	"log"
	"os"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	deleteOlderThan, err := time.ParseDuration(os.Args[1])
	if err != nil {
		log.Fatalln("parse error", err)
	}

	old := metav1.NewTime(time.Now().Add(-1 * deleteOlderThan))
	log.Println("deleting older than", old)

	config, err := clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		log.Panicln("clientcmd", err)
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Panicln("clientset", err)
	}

	nodeList, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		log.Panicln("node list", err)
	}

	for _, node := range nodeList.Items {
		for _, taint := range node.Spec.Taints {
			if taint.Key != "node.kubernetes.io/unreachable" {
				continue
			}

			log.Println(node.Name, taint.TimeAdded, taint.TimeAdded.Before(&old))
			clientset.CoreV1().Nodes().Delete(node.Name, &metav1.DeleteOptions{})
		}
	}

}
