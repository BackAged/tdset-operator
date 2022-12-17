package controllers

import (
	"context"
	"fmt"

	schedulev1 "github.com/BackAged/tdset-operator/api/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (r *TDSetReconciler) Deployment(
	ctx context.Context, req ctrl.Request,
	tdSet *schedulev1.TDSet,
) (*appsv1.Deployment, error) {
	log := log.FromContext(ctx)

	replicas, err := r.GetExpectedReplica(ctx, req, tdSet)
	if err != nil {
		log.Error(err, "failed to get expected replica")

		return nil, err
	}

	labels := map[string]string{
		"app.kubernetes.io/name":       "TDSet",
		"app.kubernetes.io/instance":   tdSet.Name,
		"app.kubernetes.io/version":    "v1",
		"app.kubernetes.io/part-of":    "tdset-operator",
		"app.kubernetes.io/created-by": "controller-manager",
	}

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      tdSet.Name,
			Namespace: tdSet.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image:           tdSet.Spec.Container.Image,
						Name:            tdSet.Name,
						ImagePullPolicy: corev1.PullIfNotPresent,
						Ports: []corev1.ContainerPort{{
							ContainerPort: int32(tdSet.Spec.Container.Port),
							Name:          "tdset",
						}},
					}},
				},
			},
		},
	}

	// Set the ownerRef for the Deployment
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/owners-dependents/
	if err := ctrl.SetControllerReference(tdSet, dep, r.Scheme); err != nil {
		log.Error(err, "failed to set controller owner reference")
		return nil, err
	}

	return dep, nil
}

func (r *TDSetReconciler) DeploymentIfNotExist(
	ctx context.Context, req ctrl.Request,
	tdSet *schedulev1.TDSet,
) (bool, error) {
	log := log.FromContext(ctx)

	dep := &appsv1.Deployment{}

	err := r.Get(ctx, types.NamespacedName{Name: tdSet.Name, Namespace: tdSet.Namespace}, dep)
	if err != nil && apierrors.IsNotFound(err) {
		dep, err := r.Deployment(ctx, req, tdSet)
		if err != nil {
			log.Error(err, "Failed to define new Deployment resource for TDSet")

			err = r.SetCondition(
				ctx, req, tdSet, TypeAvailable,
				fmt.Sprintf("Failed to create Deployment for TDSet (%s): (%s)", tdSet.Name, err),
			)
			if err != nil {
				return false, err
			}
		}

		log.Info(
			"Creating a new Deployment",
			"Deployment.Namespace", dep.Namespace,
			"Deployment.Name", dep.Name,
		)

		err = r.Create(ctx, dep)
		if err != nil {
			log.Error(
				err, "Failed to create new Deployment",
				"Deployment.Namespace", dep.Namespace,
				"Deployment.Name", dep.Name,
			)

			return false, err
		}

		err = r.GetTDSet(ctx, req, tdSet)
		if err != nil {
			log.Error(err, "Failed to re-fetch TDSet")
			return false, err
		}

		err = r.SetCondition(
			ctx, req, tdSet, TypeProgressing,
			fmt.Sprintf("Created Deployment for the TDSet: (%s)", tdSet.Name),
		)
		if err != nil {
			return false, err
		}

		return true, nil
	}

	if err != nil {
		log.Error(err, "Failed to get Deployment")

		return false, err
	}

	return false, nil
}

func (r *TDSetReconciler) UpdateDeploymentReplica(
	ctx context.Context, req ctrl.Request,
	tdSet *schedulev1.TDSet,
) error {
	log := log.FromContext(ctx)

	dep := &appsv1.Deployment{}

	err := r.Get(ctx, types.NamespacedName{Name: tdSet.Name, Namespace: tdSet.Namespace}, dep)
	if err != nil {
		log.Error(err, "Failed to get Deployment")

		return err
	}

	replicas, err := r.GetExpectedReplica(ctx, req, tdSet)
	if err != nil {
		log.Error(err, "failed to get expected replica")

		return err
	}

	if replicas == *dep.Spec.Replicas {
		return nil
	}

	log.Info(
		"Updating a Deployment replica",
		"Deployment.Namespace", dep.Namespace,
		"Deployment.Name", dep.Name,
	)

	dep.Spec.Replicas = &replicas

	err = r.Update(ctx, dep)
	if err != nil {
		log.Error(
			err, "Failed to update Deployment",
			"Deployment.Namespace", dep.Namespace,
			"Deployment.Name", dep.Name,
		)

		err = r.GetTDSet(ctx, req, tdSet)
		if err != nil {
			log.Error(err, "Failed to re-fetch TDSet")
			return err
		}

		err = r.SetCondition(
			ctx, req, tdSet, TypeProgressing,
			fmt.Sprintf("Failed to update replica for the TDSet (%s): (%s)", tdSet.Name, err),
		)
		if err != nil {
			return err
		}

		return nil
	}

	err = r.GetTDSet(ctx, req, tdSet)
	if err != nil {
		log.Error(err, "Failed to re-fetch TDSet")
		return err
	}

	err = r.SetCondition(
		ctx, req, tdSet, TypeProgressing,
		fmt.Sprintf("Updated replica for the TDSet (%s)", tdSet.Name),
	)
	if err != nil {
		return err
	}

	return nil
}
