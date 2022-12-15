package controllers

import (
	"context"
	"time"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log"

	schedulev1 "github.com/BackAged/tdset-operator/api/v1"
)

func (r *TDSetReconciler) GetExpectedReplica(ctx context.Context, req ctrl.Request, tdSet *schedulev1.TDSet) (int32, error) {
	log := log.FromContext(ctx)

	if tdSet.Spec.SchedulingConfig != nil && len(tdSet.Spec.SchedulingConfig) != 0 {
		now := time.Now()
		hour := now.Hour()

		log.Info("current server", "hour", hour, "time", now)

		for _, config := range tdSet.Spec.SchedulingConfig {
			if hour >= config.StartTime && hour < config.EndTime {
				return int32(config.Replica), nil
			}
		}
	}

	return tdSet.Spec.DefaultReplica, nil
}
