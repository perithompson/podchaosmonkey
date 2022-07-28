/*
Copyright 2022.

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

package controllers

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	podchaosv1alpha1 "github.com/perithompson/podchaosmonkey/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

// MonkeyReconciler reconciles a Monkey object
type MonkeyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=podchaos.podchaosmonkey.pt,resources=monkeys,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=podchaos.podchaosmonkey.pt,resources=monkeys/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=podchaos.podchaosmonkey.pt,resources=monkeys/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;delete;watch;

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *MonkeyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	monkeySay := ctrl.Log.WithName("controller").WithName("monkey")
	monkey := &podchaosv1alpha1.Monkey{}

	monkeySay.Info(fmt.Sprintf("Looking for monkey: %v", req.NamespacedName))
	if err := r.Get(ctx, req.NamespacedName, monkey); err != nil {
		if !apierrors.IsNotFound(err) {
			monkeySay.Error(err, fmt.Sprintf("Unable to fetch monkey: %v", req.NamespacedName))
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if len(monkey.Status.Conditions) == 0 {
		registeredCondition := metav1.Condition{
			Type:               "Registered",
			Status:             metav1.ConditionTrue,
			ObservedGeneration: 0,
			LastTransitionTime: metav1.Now(),
			Reason:             "Registered",
			Message:            "",
		}
		monkey.Status.Conditions = append(monkey.Status.Conditions, registeredCondition)
		return r.UpdateStatus(ctx, monkey)
	}
	return r.PerformExperiment(ctx, monkey)
}

//GetTarget chooses 1 pod that matches the namespace and labelselector provided to be deleted
func (r *MonkeyReconciler) GetTarget(ctx context.Context, namespace string, labelSelector metav1.LabelSelector) (corev1.Pod, error) {
	var list corev1.PodList

	rand.Seed(time.Now().UnixNano())

	selector, err := metav1.LabelSelectorAsSelector(&labelSelector)
	if err != nil {
		return corev1.Pod{}, err
	}

	if err := r.List(ctx, &list, &client.ListOptions{LabelSelector: selector}); err != nil {
		return corev1.Pod{}, err
	}
	if len(list.Items) > 0 {
		max := len(list.Items)
		randomID := rand.Intn(max)
		return list.Items[randomID], nil
	}
	return corev1.Pod{}, nil
}

//int64ToPointerint64 returns pointer of int64
func int64ToPointerint64(in int64) *int64 {
	return &in
}

//PerformExperiment deletes 1 pod that matches the namespace and labelselector provided
func (r *MonkeyReconciler) PerformExperiment(ctx context.Context, monkey *podchaosv1alpha1.Monkey) (ctrl.Result, error) {
	monkeySay := ctrl.Log.WithName("controller").WithName("monkey")
	target, err := r.GetTarget(ctx, monkey.Spec.Namespace, monkey.Spec.Selector)
	if err != nil {
		return ctrl.Result{}, err
	}
	requeueInterval, err := GetMinInterval(monkey.Spec.Interval)
	if err != nil {
		return ctrl.Result{RequeueAfter: requeueInterval}, err
	}
	if target.GetUID() != "" {
		podname := client.ObjectKeyFromObject(&target)
		if monkey.Spec.Noop {
			monkeySay.Info(fmt.Sprintf("No Operation specified ==== Would have deleted pod: %s", podname))
			return ctrl.Result{RequeueAfter: requeueInterval}, nil
		}
		if err := r.Delete(ctx, &target, &client.DeleteOptions{GracePeriodSeconds: int64ToPointerint64(0)}); err != nil {
			return ctrl.Result{RequeueAfter: requeueInterval}, err
		}
		monkeySay.Info(fmt.Sprintf("Deleted Pod: %s", podname))
	}
	return ctrl.Result{RequeueAfter: requeueInterval}, nil
}

//UpdateStatus updates the status of the Monkey Object
func (r *MonkeyReconciler) UpdateStatus(ctx context.Context, monkey *podchaosv1alpha1.Monkey) (ctrl.Result, error) {
	monkeySay := ctrl.Log.WithName("controller").WithName("monkey")
	newKey := client.ObjectKeyFromObject(monkey)
	newObject := &podchaosv1alpha1.Monkey{}
	if err := r.Get(ctx, newKey, newObject); err != nil {
		if !apierrors.IsNotFound(err) {
			monkeySay.Error(err, fmt.Sprintf("Unable to fetch monkey: %v", newKey))
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	newObject.Status = monkey.Status
	if err := r.Status().Update(ctx, newObject); err != nil {
		monkeySay.Error(err, fmt.Sprintf("Unable to update Monkey: %v", newObject.GetName()))
	}
	requeueInterval, err := GetMinInterval(monkey.Spec.Interval)
	if err != nil {
		return ctrl.Result{RequeueAfter: requeueInterval}, err
	}
	return ctrl.Result{RequeueAfter: requeueInterval}, err
}

//GetMinInterval Gets the minimal intervals for Chaos to occur
func GetMinInterval(interval string) (time.Duration, error) {
	if interval == "" {
		return time.ParseDuration("30s")
	}
	return time.ParseDuration(interval)
}

// SetupWithManager sets up the controller with the Manager.
func (r *MonkeyReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&podchaosv1alpha1.Monkey{}).
		Complete(r)
}
