package action

import (
	"github.com/operator-framework/operator-sdk/pkg/sdk"
	"github.com/sirupsen/logrus"
	"github.com/syndesisio/syndesis-operator/pkg/apis/syndesis/v1alpha1"
	"math"
	"time"
)

const (
	// Number of times a syndesis upgrade will be triggered (including the first one launched by another state)
	UpgradeMaxAttempts = 4
)

// After a failure, waits a exponential amount of time, then retries
type UpgradeBackoff struct {
	operatorVersion	string
}

func (a *UpgradeBackoff) CanExecute(syndesis *v1alpha1.Syndesis) bool {
	return syndesisInstallationStatusIs(syndesis, v1alpha1.SyndesisInstallationStatusUpgradeFailureBackoff)
}

func (a *UpgradeBackoff) Execute(syndesis *v1alpha1.Syndesis) error {

	// Check number of attempts to fail fast
	if syndesis.Status.UpgradeAttempts > UpgradeMaxAttempts {
		logrus.Info("Upgrade of Syndesis resource ", syndesis.Name, " failed too many times and will not be retried")

		target := syndesis.DeepCopy()
		target.Status.InstallationStatus = v1alpha1.SyndesisInstallationStatusUpgradeFailed
		target.Status.Reason = v1alpha1.SyndesisStatusReasonTooManyUpgradeAttempts
		target.Status.Description = "Upgrade failed too many times and will not be retried"
		target.Status.ForceUpgrade = false

		return sdk.Update(target)
	}

	now := time.Now()

	lastFailureWrapper := syndesis.Status.LastUpgradeFailure
	var lastFailure time.Time
	if lastFailureWrapper != nil {
		lastFailure = lastFailureWrapper.Time
	} else {
		lastFailure = time.Now().Add(- 8 * time.Hour)
	}

	if lastFailure.After(now) {
		lastFailure = now
	}

	power := float64(syndesis.Status.UpgradeAttempts - 1)
	if power < 0 {
		power = 0
	}

	delay := time.Duration(math.Pow(2, power)) * time.Minute

	nextAttempt := lastFailure.Add(delay)

	if now.After(nextAttempt) {
		logrus.Info("Restarting upgrade process for Syndesis resource ", syndesis.Name)

		target := syndesis.DeepCopy()
		target.Status.InstallationStatus = v1alpha1.SyndesisInstallationStatusUpgrading
		target.Status.Reason = v1alpha1.SyndesisStatusReasonMissing
		target.Status.Description = ""
		target.Status.ForceUpgrade = true

		return sdk.Update(target)
	} else {
		remaining := math.Round(nextAttempt.Sub(now).Seconds())
		logrus.Info("Upgrade of Syndesis resource ", syndesis.Name, " will be retried in ", remaining, " seconds")
		return nil
	}
}