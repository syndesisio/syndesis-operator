package syndesis

import (
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

	for _, a := range actionPool {
		if a.CanExecute(syndesis) {
			if err := a.Execute(syndesis); err != nil {
				return err
			}
		}
	}

	return nil
}
