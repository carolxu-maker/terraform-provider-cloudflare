package r2_bucket_lifecycle

import (
	"encoding/json"
	"time"

	"github.com/cloudflare/terraform-provider-cloudflare/internal/apijson"
	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
)

func (m R2BucketLifecycleModel) marshalCustom() (data []byte, err error) {
	if data, err = apijson.MarshalRoot(m); err != nil {
		return
	}
	if data, err = m.marshalRulesWithFullTimestamps(data); err != nil {
		return
	}
	return
}

func (m R2BucketLifecycleModel) marshalCustomForUpdate(state R2BucketLifecycleModel) (data []byte, err error) {
	if data, err = apijson.MarshalForUpdate(m, state); err != nil {
		return
	}
	if data, err = m.marshalRulesWithFullTimestamps(data); err != nil {
		return
	}
	return
}

func (m R2BucketLifecycleModel) marshalRulesWithFullTimestamps(b []byte) (data []byte, err error) {
	// If there are no rules, return the original data
	if m.Rules == nil {
		return b, nil
	}

	// Parse the existing JSON
	var jsonData map[string]interface{}
	if err = json.Unmarshal(b, &jsonData); err != nil {
		return nil, err
	}

	// Extract the rules array
	rulesInterface, exists := jsonData["rules"]
	if !exists {
		return b, nil
	}

	rules, ok := rulesInterface.([]interface{})
	if !ok {
		return b, nil
	}

	// Process each rule to ensure date fields have full timestamp precision
	for i, ruleInterface := range rules {
		rule, ok := ruleInterface.(map[string]interface{})
		if !ok {
			continue
		}

		// Check if we have the corresponding model rule
		if i >= len(*m.Rules) {
			continue
		}
		modelRule := (*m.Rules)[i]

		// Handle delete objects transition condition date
		if modelRule.DeleteObjectsTransition != nil && modelRule.DeleteObjectsTransition.Condition != nil {
			if err = m.updateDateFieldInCondition(rule, "deleteObjectsTransition", "condition", modelRule.DeleteObjectsTransition.Condition.Date); err != nil {
				return nil, err
			}
		}

		// Handle storage class transitions condition dates
		if modelRule.StorageClassTransitions != nil {
			storageTransitions, exists := rule["storageClassTransitions"]
			if exists {
				storageTransitionsArray, ok := storageTransitions.([]interface{})
				if ok && len(storageTransitionsArray) == len(*modelRule.StorageClassTransitions) {
					for j, transitionInterface := range storageTransitionsArray {
						transition, ok := transitionInterface.(map[string]interface{})
						if !ok {
							continue
						}
						modelTransition := (*modelRule.StorageClassTransitions)[j]
						if modelTransition.Condition != nil {
							if err = m.updateDateFieldInCondition(transition, "condition", "", modelTransition.Condition.Date); err != nil {
								return nil, err
							}
						}
					}
				}
			}
		}
	}

	return json.Marshal(jsonData)
}

func (m R2BucketLifecycleModel) updateDateFieldInCondition(rule map[string]interface{}, conditionPath string, subPath string, dateField timetypes.RFC3339) error {
	// Navigate to the condition object
	var condition map[string]interface{}
	var exists bool

	if subPath != "" {
		// For nested paths like deleteObjectsTransition.condition
		parentObj, parentExists := rule[conditionPath]
		if !parentExists {
			return nil
		}
		parentMap, ok := parentObj.(map[string]interface{})
		if !ok {
			return nil
		}
		conditionObj, conditionExists := parentMap[subPath]
		if !conditionExists {
			return nil
		}
		condition, exists = conditionObj.(map[string]interface{})
	} else {
		// For direct paths like condition
		conditionObj, conditionExists := rule[conditionPath]
		if !conditionExists {
			return nil
		}
		condition, exists = conditionObj.(map[string]interface{})
	}

	if !exists {
		return nil
	}

	// Update the date field if it exists and the model has a non-null date
	if !dateField.IsNull() && !dateField.IsUnknown() {
		dateValue, diags := dateField.ValueRFC3339Time()
		if diags.HasError() {
			return nil
		}
		// Use RFC3339Nano to preserve full timestamp precision
		condition["date"] = dateValue.Format(time.RFC3339Nano)
	}

	return nil
}
