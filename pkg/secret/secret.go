package secret

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

// SecretManager is a type for fetching Secrets whether or not they are in the
// client cache, labelling so that they will be included in the client cache,
// and optionally setting an owner reference.
type SecretManager struct {
	log       logr.Logger
	client    client.Client
	apiReader client.Reader
}

// NewSecretManager returns a new SecretManager.
func NewSecretManager(log logr.Logger, cacheClient client.Client, apiReader client.Reader) *SecretManager {
	return &SecretManager{
		log:       log.WithName("secret_manager"),
		client:    cacheClient,
		apiReader: apiReader,
	}
}

// AcquireSecret retrieves a Secret and ensures that it has a label that will
// ensure it is present in the cache (and that we can watch for changes), and
// that it has a particular owner reference. The owner reference may optionally
// be a controller reference.
func (sm *SecretManager) AcquireSecret(ctx context.Context, key types.NamespacedName, owner client.Object, ownerIsController, addFinalizer bool) (*corev1.Secret, error) {
	if owner == nil {
		panic("AcquireSecret called with no owner")
	}

	secret, err := sm.findSecret(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to find secret: %w", err)
	}

	err = sm.claimSecret(ctx, secret, owner, ownerIsController, addFinalizer)

	return secret, err
}

// findSecret retrieves a Secret from the cache if it is available, and from the
// k8s API if not.
func (sm *SecretManager) findSecret(ctx context.Context, key types.NamespacedName) (secret *corev1.Secret, err error) {
	secret = &corev1.Secret{}

	// Look for secret in the filtered cache
	err = sm.client.Get(ctx, key, secret)
	if err == nil {
		return secret, nil
	}
	if !apierrors.IsNotFound(err) {
		return nil, err
	}

	// Secret not in cache; check API directly for unlabelled Secret
	err = sm.apiReader.Get(ctx, key, secret)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

// claimSecret ensures that the Secret has a label that will ensure it is
// present in the cache (and that we can watch for changes), and optionally
// that it has a particular owner reference.
func (sm *SecretManager) claimSecret(ctx context.Context, secret *corev1.Secret, owner client.Object, ownerIsController, addFinalizer bool) error {
	needsUpdate := false
	if !metav1.HasLabel(secret.ObjectMeta, LabelEnvironmentName) {
		metav1.SetMetaDataLabel(&secret.ObjectMeta, LabelEnvironmentName, LabelEnvironmentValue)
		needsUpdate = true
	}
	if owner != nil {
		if ownerIsController {
			if !metav1.IsControlledBy(secret, owner) {
				if err := controllerutil.SetControllerReference(owner, secret, sm.client.Scheme()); err != nil {
					return fmt.Errorf("failed to set secret controller reference: %w", err)
				}
				needsUpdate = true
			}
		} else {
			alreadyOwned := false
			ownerUID := owner.GetUID()
			for _, ref := range secret.GetOwnerReferences() {
				if ref.UID == ownerUID {
					alreadyOwned = true
					break
				}
			}
			if !alreadyOwned {
				if err := controllerutil.SetOwnerReference(owner, secret, sm.client.Scheme()); err != nil {
					return fmt.Errorf("failed to set secret owner reference: %w", err)
				}
				needsUpdate = true
			}
		}
	}

	if addFinalizer && (controllerutil.AddFinalizer(secret, SecretFinalizer) ||
		controllerutil.RemoveFinalizer(secret, DeprecatedSecretFinalizer)) {
		needsUpdate = true
	}

	if needsUpdate {
		if err := sm.client.Update(ctx, secret); err != nil {
			return fmt.Errorf("failed to update secret %s in namespace %s: %w", secret.ObjectMeta.Name, secret.ObjectMeta.Namespace, err)
		}
	}

	return nil
}
