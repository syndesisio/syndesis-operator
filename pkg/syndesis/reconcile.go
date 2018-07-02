package syndesis

import api "github.com/syndesisio/syndesis-operator/pkg/apis/syndesis/v1alpha1"

// Reconcile the state of the Syndesis infrastructure elements
// For a new installation: Instantiate the Syndesis Template with the parameters set in the resource
// For an update: execute the Update Pod for the desired version.
func Reconcile(syndesis *api.Syndesis) (err error) {
	return nil
}