/*
Copyright 2025.

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

package controller

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/cluster-api/util"
	"sigs.k8s.io/cluster-api/util/annotations"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/aniruddha2000/cluster-api-provider-stackit/api/v1alpha1"
	infrastructurev1alpha1 "github.com/aniruddha2000/cluster-api-provider-stackit/api/v1alpha1"
	"github.com/aniruddha2000/cluster-api-provider-stackit/pkg/cloud/scope"
	secretManager "github.com/aniruddha2000/cluster-api-provider-stackit/pkg/secret"
	"github.com/stackitcloud/stackit-sdk-go/core/config"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
)

// STACKITClusterReconciler reconciles a STACKITCluster object
type STACKITClusterReconciler struct {
	client.Client
	APIReader client.Reader
	Scheme    *runtime.Scheme
}

// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=stackitclusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=stackitclusters/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infrastructure.cluster.x-k8s.io,resources=stackitclusters/finalizers,verbs=update

// Reconcile reconciles STACKIT Cluster object.
func (r *STACKITClusterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx).WithValues("controller", "STACKITCluster")

	stackitCluster := &v1alpha1.STACKITCluster{}
	if err := r.Client.Get(ctx, req.NamespacedName, stackitCluster); err != nil {
		if apierrors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}

		return reconcile.Result{}, fmt.Errorf("failed to get STACKIT Cluster: %w", err)
	}

	cluster, err := util.GetOwnerCluster(ctx, r.Client, stackitCluster.ObjectMeta)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get owner cluster: %w", err)
	}

	if cluster == nil {
		log.Info("Cluster Controller has not yet set OwnerRef")
		return ctrl.Result{}, nil
	}

	if annotations.IsPaused(cluster, stackitCluster) {
		log.Info("GCPCluster of linked Cluster is marked as paused. Won't reconcile")
		return ctrl.Result{}, nil
	}

	manager := secretManager.NewSecretManager(log, r.Client, r.APIReader)
	token, _, err := getAndValidateSTACKITToken(ctx, req.Namespace, stackitCluster, manager)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get STACKIT token: %w", err)
	}

	loadbalancerClient, err := loadbalancer.NewAPIClient(config.WithToken(token))
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to create loadbalancer client: %w", err)
	}

	clusterScope, err := scope.NewClusterScope(ctx, scope.ClusterScopeParams{
		STACKITCluster:            stackitCluster,
		Cluster:                   cluster,
		STACKITLoadbalancerClient: loadbalancerClient,
	})
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to create cluster scope: %w", err)
	}

	return ctrl.Result{}, nil
}

func getAndValidateSTACKITToken(ctx context.Context, namespace string, stackitCluster *v1alpha1.STACKITCluster, secretManager *secretManager.SecretManager) (string, *corev1.Secret, error) {
	secretNamspacedName := types.NamespacedName{Namespace: namespace, Name: stackitCluster.Spec.STACKITToken.Name}

	stackitSecret, err := secretManager.AcquireSecret(
		ctx,
		secretNamspacedName,
		stackitCluster,
		false,
		stackitCluster.DeletionTimestamp.IsZero(),
	)
	if err != nil {
		if apierrors.IsNotFound(err) {
			return "", nil, fmt.Errorf("The STACKIT secret %s does not exist: %w", secretNamspacedName, &err)
		}
		return "", nil, err
	}

	stackitToken := string(stackitSecret.Data[stackitCluster.Spec.STACKITToken.Key])

	// Validate token
	if stackitToken == "" {
		return "", nil, fmt.Errorf("invalid token: empty")
	}

	return stackitToken, stackitSecret, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *STACKITClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrastructurev1alpha1.STACKITCluster{}).
		Named("stackitcluster").
		Complete(r)
}
