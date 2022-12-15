package controllers

import (
	"context"

	"sigs.k8s.io/controller-runtime/pkg/log"

	schedulev1 "github.com/BackAged/tdset-operator/api/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

// ConditionStatus defines TDSet condition status.
type ConditionStatus string

// Defines TDSet condition status.
const (
	TypeAvailable   ConditionStatus = "Available"
	TypeProgressing ConditionStatus = "Progressing"
	TypeDegraded    ConditionStatus = "Degraded"
)

// GetTDSet gets the TDSet from api server.
func (r *TDSetReconciler) GetTDSet(ctx context.Context, req ctrl.Request, tdSet *schedulev1.TDSet) error {
	err := r.Get(ctx, req.NamespacedName, tdSet)
	if err != nil {
		return err
	}

	return nil
}

// SetInitialCondition sets the status condition of the TDSet to available initially
// when no condition exists yet.
func (r *TDSetReconciler) SetInitialCondition(ctx context.Context, req ctrl.Request, tdSet *schedulev1.TDSet) error {
	if tdSet.Status.Conditions != nil || len(tdSet.Status.Conditions) != 0 {
		return nil
	}

	err := r.SetCondition(ctx, req, tdSet, TypeAvailable, "Starting reconciliation")

	return err
}

// SetCondition sets the status condition of the TDSet.
func (r *TDSetReconciler) SetCondition(
	ctx context.Context, req ctrl.Request,
	tdSet *schedulev1.TDSet, condition ConditionStatus,
	message string,
) error {
	log := log.FromContext(ctx)

	meta.SetStatusCondition(
		&tdSet.Status.Conditions,
		metav1.Condition{
			Type:   string(condition),
			Status: metav1.ConditionUnknown, Reason: "Reconciling",
			Message: message,
		},
	)

	if err := r.Status().Update(ctx, tdSet); err != nil {
		log.Error(err, "Failed to update TDSet status")

		return err
	}

	if err := r.Get(ctx, req.NamespacedName, tdSet); err != nil {
		log.Error(err, "Failed to re-fetch TDSet")

		return err
	}

	return nil
}
