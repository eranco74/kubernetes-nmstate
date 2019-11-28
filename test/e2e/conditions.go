package e2e

import (
	"context"
	"time"

	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/types"
	yaml "sigs.k8s.io/yaml"

	framework "github.com/operator-framework/operator-sdk/pkg/test"

	corev1 "k8s.io/api/core/v1"

	nmstatev1alpha1 "github.com/nmstate/kubernetes-nmstate/pkg/apis/nmstate/v1alpha1"
)

type expectedConditionsStatus struct {
	Node       string
	conditions nmstatev1alpha1.ConditionList
}

func conditionsToYaml(conditions nmstatev1alpha1.ConditionList) string {
	manifest, err := yaml.Marshal(conditions)
	if err != nil {
		panic(err)
	}
	return string(manifest)
}

func nodeNetworkConfigurationEnactment(key types.NamespacedName) nmstatev1alpha1.NodeNetworkConfigurationEnactment {
	state := nmstatev1alpha1.NodeNetworkConfigurationEnactment{}
	Eventually(func() error {
		return framework.Global.Client.Get(context.TODO(), key, &state)
	}, ReadTimeout, ReadInterval).ShouldNot(HaveOccurred())
	return state
}

func enactmentConditionsStatus(node string) nmstatev1alpha1.ConditionList {
	//TODO: Take the format from pkg
	key := types.NamespacedName{Name: node + "-" + TestPolicy}
	enactment := nodeNetworkConfigurationEnactment(key)
	enactmentsConditionTypes := []nmstatev1alpha1.ConditionType{
		nmstatev1alpha1.NodeNetworkConfigurationEnactmentConditionAvailable,
		nmstatev1alpha1.NodeNetworkConfigurationEnactmentConditionFailing,
		nmstatev1alpha1.NodeNetworkConfigurationEnactmentConditionProgressing,
	}
	obtainedConditions := nmstatev1alpha1.ConditionList{}
	for _, enactmentsConditionType := range enactmentsConditionTypes {
		obtainedCondition := enactment.Status.Conditions.Find(enactmentsConditionType)
		obtainedConditionStatus := corev1.ConditionUnknown
		if obtainedCondition != nil {
			obtainedConditionStatus = obtainedCondition.Status
		}
		obtainedConditions = append(obtainedConditions, nmstatev1alpha1.Condition{
			Type:   enactmentsConditionType,
			Status: obtainedConditionStatus,
		})
	}
	return obtainedConditions
}

func enactmentConditionsStatusEventually(node string) AsyncAssertion {
	return Eventually(func() nmstatev1alpha1.ConditionList {
		return enactmentConditionsStatus(node)
	}, 180*time.Second, 1*time.Second)
}

func enactmentConditionsStatusConsistently(node string) AsyncAssertion {
	return Consistently(func() nmstatev1alpha1.ConditionList {
		return enactmentConditionsStatus(node)
	}, 5*time.Second, 1*time.Second)
}