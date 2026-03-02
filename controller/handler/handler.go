package handler

import (
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type Handler struct {
	PodQueue     workqueue.TypedRateLimitingInterface[*corev1.Pod]
	NodeQueue    workqueue.TypedRateLimitingInterface[*corev1.Node]
	PodInformer  cache.SharedIndexInformer
	NodeInformer cache.SharedIndexInformer
}

func (h *Handler) OnAdd(obj interface{}, b bool) {
	switch obj.(type) {
	case *corev1.Pod:
		h.PodQueue.Add(obj.(*corev1.Pod))
	case *corev1.Node:
		h.NodeQueue.Add(obj.(*corev1.Node))
	}
}

func (h *Handler) OnUpdate(old, new interface{}) {
	switch new.(type) {
	case *corev1.Pod:
		h.PodQueue.Add(new.(*corev1.Pod))
	case *corev1.Node:
		h.NodeQueue.Add(new.(*corev1.Node))
	}
}

func (h *Handler) OnDelete(obj interface{}) {
	switch obj.(type) {
	case *corev1.Pod:
		h.PodQueue.Add(obj.(*corev1.Pod))
	case *corev1.Node:
		h.NodeQueue.Add(obj.(*corev1.Node))
	}
}

func NewHandler(podInformer cache.SharedIndexInformer, nodeInformer cache.SharedIndexInformer) *Handler {
	return &Handler{
		PodQueue: workqueue.NewTypedRateLimitingQueue(
			workqueue.DefaultTypedControllerRateLimiter[*corev1.Pod](),
		),
		NodeQueue: workqueue.NewTypedRateLimitingQueue(
			workqueue.DefaultTypedControllerRateLimiter[*corev1.Node](),
		),
		PodInformer:  podInformer,
		NodeInformer: nodeInformer,
	}
}

func (h *Handler) ProcessNode(item *corev1.Node) error {
	fmt.Println("Processing Node:", item)
	return errors.New("")
}

func (h *Handler) ProcessPod(item *corev1.Pod) error {
	fmt.Println("Processing Pod:", item)
	return errors.New("")
}

func (h *Handler) PeriodicFunction() {
	podCache := h.PodInformer.GetStore()
	podCache.List()

	nodeCache := h.NodeInformer.GetStore()
	nodeCache.List()
}
