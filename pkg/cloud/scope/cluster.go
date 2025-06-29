package scope

import (
	"context"
	"errors"
	"fmt"

	"github.com/aniruddha2000/cluster-api-provider-stackit/api/v1alpha1"
	"github.com/stackitcloud/stackit-sdk-go/services/loadbalancer"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/patch"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ClusterScopeParams defines the input parameters used to create a new Scope.
type ClusterScopeParams struct {
	Client client.Client

	Cluster        *clusterv1.Cluster
	STACKITCluster *v1alpha1.STACKITCluster

	STACKITLoadbalancerClient *loadbalancer.APIClient
}

// ClusterScope defines the basic context for an actuator to operate upon.
type ClusterScope struct {
	client      client.Client
	patchHelper *patch.Helper

	Cluster        *clusterv1.Cluster
	STACKITCluster *v1alpha1.STACKITCluster
}

// NewClusterScope creates a new Scope from the supplied parameters.
// This is meant to be called for each reconcile iteration.
func NewClusterScope(ctx context.Context, params ClusterScopeParams) (*ClusterScope, error) {
	if params.Cluster == nil {
		return nil, errors.New("failed to generate new scope from nil Cluster")
	}
	if params.STACKITCluster == nil {
		return nil, errors.New("failed to generate new scope from nil GCPCluster")
	}

	helper, err := patch.NewHelper(params.STACKITCluster, params.Client)
	if err != nil {
		return nil, fmt.Errorf("failed to init patch helper: %w", err)
	}

	return &ClusterScope{
		client:         params.Client,
		Cluster:        params.Cluster,
		STACKITCluster: params.STACKITCluster,
		patchHelper:    helper,
	}, nil
}
