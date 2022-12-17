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
	"time"

	appsv1 "k8s.io/api/apps/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	schedulev1 "github.com/BackAged/tdset-operator/api/v1"
)

const (
	// minute
	DefaultReconciliationInterval = 5
)

// TDSetReconciler reconciles a TDSet object
type TDSetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=schedule.rs,resources=tdsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=schedule.rs,resources=tdsets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=schedule.rs,resources=tdsets/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *TDSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	log.Info("starting reconciliation")

	tdSet := &schedulev1.TDSet{}

	// Get the TDSet
	err := r.GetTDSet(ctx, req, tdSet)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("TDSet resource not found. Ignoring since object must be deleted")

			return ctrl.Result{}, nil
		}

		log.Error(err, "Failed to get TDSet")

		return ctrl.Result{}, err
	}

	// Try to set initial condition status
	err = r.SetInitialCondition(ctx, req, tdSet)
	if err != nil {
		log.Error(err, "failed to set initial condition")

		return ctrl.Result{}, err
	}

	// TODO: Delete finalizer

	// Deployment if not exist
	ok, err := r.DeploymentIfNotExist(ctx, req, tdSet)
	if err != nil {
		log.Error(err, "failed to deploy deployment for TDSet")

		return ctrl.Result{}, err
	}

	if ok {
		return ctrl.Result{RequeueAfter: time.Minute}, nil
	}

	// Update deployment replica if mis matched.
	err = r.UpdateDeploymentReplica(ctx, req, tdSet)
	if err != nil {
		log.Error(err, "failed to update deployment for TDSet")

		return ctrl.Result{}, err
	}

	interval := DefaultReconciliationInterval
	if tdSet.Spec.IntervalMint != 0 {
		interval = int(tdSet.Spec.IntervalMint)
	}

	log.Info("ending reconciliation")

	return ctrl.Result{RequeueAfter: time.Duration(time.Minute * time.Duration(interval))}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *TDSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&schedulev1.TDSet{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
