package syndesis

import (
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	api "github.com/syndesisio/syndesis-operator/pkg/apis/syndesis/v1alpha1"
	"github.com/syndesisio/syndesis-operator/pkg/syndesis/action"
)

var (
	actionPool	[]action.InstallationAction
)

func init() {
	actionPool = append(actionPool,
		&action.Initialize{},
		&action.Install{},
		&action.Startup{},
	)
}

// Reconcile the state of the Syndesis infrastructure elements
func Reconcile(syndesis *api.Syndesis, deleted bool) error {

	if deleted {
		// No specific actions to do on deletion
		return nil
	}

	// Don't want to do anything if the syndesis resource has been updated in the meantime
	// This happens when a processing takes more tha the resync period
	if latest, err := isLatestVersion(syndesis); err != nil || !latest {
		return err
	}

	for _, a := range actionPool {
		if a.CanExecute(syndesis) {
			if err := a.Execute(syndesis); err != nil {
				return err
			}
		}
	}

	return nil
}

func isLatestVersion(syndesis *api.Syndesis) (bool, error) {
	refreshed := syndesis.DeepCopy()
	if err := sdk.Get(refreshed); err != nil {
		return false, err
	}
	return refreshed.ResourceVersion == syndesis.ResourceVersion, nil
}