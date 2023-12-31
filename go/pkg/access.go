package pal

import (
	"fmt"
)

func (pal *Client) ProcessAccessRequest(handleAccess HandleAccessFunc, dataSubjectLocator Locator, dataSubjectID string) (map[string]interface{}, error) {
	fmt.Printf("Processing access request for data subject %s\n", dataSubjectID)
	if dataSubjectLocator.LocatorType != Document {
		return nil, fmt.Errorf("%s data subject locator type must be document", ACCESS_REQUEST_ERROR)
	}
	locAndObj, err := pal.dbClient.getDocument(dataSubjectLocator)
	if err != nil {
		return nil, fmt.Errorf("%s %w", ACCESS_REQUEST_ERROR, err)
	}
	dataSubject := locAndObj.Object
	data, err := pal.processAccessRequest(handleAccess, dataSubject, dataSubjectID, dataSubjectLocator)
	if err != nil {
		return nil, fmt.Errorf("%s %w", ACCESS_REQUEST_ERROR, err)
	}

	return data, nil
}

func (pal *Client) processAccessRequest(handleAccess HandleAccessFunc, dataNode DatabaseObject, dataSubjectID string, dataNodeLocator Locator) (map[string]interface{}, error) {

	data, err := handleAccess(dataSubjectID, dataNodeLocator, dataNode)
	if err != nil {
		return nil, err
	}
	report := make(map[string]interface{})

	for key, value := range data {
		if loc, ok := value.(Locator); ok {
			// if locator, recursively process
			retData, err := pal.processLocator(handleAccess, loc, dataSubjectID)
			if err != nil {
				return nil, err
			}
			report[key] = retData
		} else if locs, ok := value.([]Locator); ok {
			// if locator slice, recursively process each locator
			report[key] = make([]interface{}, 0)
			for _, loc := range locs {
				retData, err := pal.processLocator(handleAccess, loc, dataSubjectID)
				if err != nil {
					return nil, err
				}
				report[key] = append(report[key].([]interface{}), retData)
			}
		} else if locMap, ok := value.(map[string]Locator); ok {
			// if map, recursively process each locator
			report[key] = make(map[string]interface{})
			for k, loc := range locMap {
				retData, err := pal.processLocator(handleAccess, loc, dataSubjectID)
				if err != nil {
					return nil, err
				}
				report[key].(map[string]interface{})[k] = retData
			}
		} else {
			// else, directly add to report
			report[key] = value
		}
	}

	return report, nil
}

func (pal *Client) processLocator(handleAccess HandleAccessFunc, loc Locator, dataSubjectID string) (interface{}, error) {
	err := validateLocator(loc)
	if err != nil {
		return nil, err
	}
	if loc.LocatorType == Document {
		locAndObj, err := pal.dbClient.getDocument(loc)
		if err != nil {
			return nil, err
		}
		dataNode := locAndObj.Object
		retData, err := pal.processAccessRequest(handleAccess, dataNode, dataSubjectID, loc)
		if err != nil {
			return nil, err
		}
		return retData, nil
	} else if loc.LocatorType == Collection {
		locAndObjs, err := pal.dbClient.getDocuments(loc)
		if err != nil {
			return nil, err
		}

		dataNodes := make([]DatabaseObject, 0)
		for _, locAndObj := range locAndObjs {
			dataNodes = append(dataNodes, locAndObj.Object)
		}

		var retData []interface{}
		for _, dataNode := range dataNodes {
			currDataNodeData, err := pal.processAccessRequest(handleAccess, dataNode, dataSubjectID, loc)
			if err != nil {
				return nil, err
			}
			retData = append(retData, currDataNodeData)
		}
		return retData, nil

	}
	return nil, fmt.Errorf("invalid locator type")
}
