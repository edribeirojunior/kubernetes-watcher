/*
Copyright 2016 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Note: the example only works with the code within the same release/branch.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

//Container image names
type Container struct {
	Image string `json:"image"`
}

//PodObject values
type PodObject struct {
	PodName        string      `json:"pod_name"`
	Namespace      string      `json:"namespace"`
	Containers     []Container `json:"containers"`
	InitContainers []Container `json:"initcontainers,omitempty"`
}

//RunningObject standard response for api
type RunningObject struct {
	Pods []PodObject `json:"pods"`
}

var runningPods RunningObject

//ResourceF function to render in mux
func ResourceF(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(runningPods)
}

func main() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	var runningPodsObjects []PodObject

	for {
		// get pods in all the namespaces by omitting namespace
		// Or specify namespace to get pods in particular namespace
		pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

		for _, containers := range pods.Items {
			c := containers
			if c.Status.Phase == "Running" {
				var imageContainer []Container
				var imageinitContainer []Container
				for _, images := range c.Spec.Containers {
					image := Container{Image: images.Image}
					imageContainer = append(imageContainer, image)
				}

				for _, imagesinit := range c.Spec.InitContainers {
					imageinit := Container{Image: imagesinit.Image}
					imageinitContainer = append(imageinitContainer, imageinit)
				}

				v := PodObject{c.ObjectMeta.Name, c.ObjectMeta.Namespace, imageContainer, imageinitContainer}
				runningPodsObjects = append(runningPodsObjects, v)

			}
		}

		fmt.Println("Creating runningPods object...")

		runningPods = RunningObject{Pods: runningPodsObjects}
		fmt.Println("[INFO] Application is running")

		router := mux.NewRouter()
		fmt.Println("[INFO] The API is running...")
		router.HandleFunc("/", ResourceF).Methods("GET")
		log.Fatal(http.ListenAndServe(":8090", router))

		time.Sleep(10 * time.Second)

	}
}
