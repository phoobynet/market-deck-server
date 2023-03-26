package facts

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

type jObject = map[string]interface{}

func parseFacts(data []byte) []Fact {

	var companyFacts jObject

	if err := json.Unmarshal(data, &companyFacts); err != nil {
		logrus.Errorf("Error unmarshalling company facts: %v", err)
	}

	facts := companyFacts["facts"].(jObject)

	factUnits := make([]Fact, 0)

	// first layer, dei, and us-gaap
	for root, roots := range facts {
		// e.g. AccountsPayableCurrent
		for conceptKey, concept := range roots.(jObject) {
			for unitType, facts := range concept.(jObject)["units"].(jObject) {

				for _, factUnitJObject := range facts.([]interface{}) {
					factUnits = append(
						factUnits, parseFactUnit(root, conceptKey, unitType, factUnitJObject.(jObject)),
					)
				}
			}
		}
	}

	return factUnits
}

func pickValue[T any](j jObject, key string) T {
	if v, ok := j[key]; ok {
		return v.(T)
	}

	var noop T

	return noop
}

func parseFactUnit(root string, conceptKey string, unitType string, factUnitJObject jObject) Fact {
	return Fact{
		Label:           pickValue[string](factUnitJObject, "label"),
		Root:            root,
		Concept:         conceptKey,
		UnitType:        unitType,
		EndDate:         pickValue[string](factUnitJObject, "end"),
		Value:           pickValue[string](factUnitJObject, "val"),
		AccessionNumber: pickValue[string](factUnitJObject, "accn"),
		FinancialYear:   int(factUnitJObject["fy"].(float64)),
		FinancialPeriod: pickValue[string](factUnitJObject, "fp"),
		Form:            pickValue[string](factUnitJObject, "form"),
		Filed:           pickValue[string](factUnitJObject, "filed"),
		Frame:           pickValue[string](factUnitJObject, "frame"),
	}
}
