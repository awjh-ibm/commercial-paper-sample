package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/contractapi"
)

func assetIDToKey(ctx contractapi.TransactionContext, objectType string, assetID string) (string, error) {
	key, err := ctx.GetStub().CreateCompositeKey(objectType, []string{assetID})

	if err != nil {
		return "", fmt.Errorf("Failed to generate key for %s %s", objectType, assetID)
	}

	return key, nil
}

func serialize(data interface{}) ([]byte, error) {
	bytes, err := json.Marshal(data)
	logger.Info(string(bytes))
	if err != nil {
		return nil, fmt.Errorf("Failed to serialize")
	}

	return bytes, nil
}

func deserializeMarket(data []byte, obj *market) error {
	err := json.Unmarshal(data, obj)

	if err != nil {
		logger.Error(err)
		return fmt.Errorf("Failed to deserialize market: %s", string(data))
	}

	return nil
}

func deserializePaper(data []byte, obj *commercialPaper) error {
	err := json.Unmarshal(data, obj)

	if err != nil {
		logger.Error(err)
		return fmt.Errorf("Failed to deserialize paper: %s", string(data))
	}

	return nil
}

func get(ctx contractapi.TransactionContext, objectType string, marketID string) (string, []byte, error) {
	key, err := assetIDToKey(ctx, objectType, marketID)

	if err != nil {
		return "", nil, err
	}

	objectJSON, err := ctx.GetStub().GetState(key)

	if err != nil {
		return "", nil, errors.New("Failed to read world state")
	}

	if objectJSON == nil {
		return "", nil, fmt.Errorf("%s %s does not exist", objectType, marketID)
	}

	return key, objectJSON, nil
}

func getMarket(ctx contractapi.TransactionContext, marketID string) (*market, error) {
	key, marketJSON, err := get(ctx, "market", marketID)

	if err != nil {
		return nil, err
	}

	marketObj := new(market)
	marketObj.key = key
	err = deserializeMarket(marketJSON, marketObj)

	if err != nil {
		return nil, err
	}

	return marketObj, nil
}

func getPaper(ctx contractapi.TransactionContext, paperID string) (*commercialPaper, error) {
	key, paperJSON, err := get(ctx, "paper", paperID)

	if err != nil {
		return nil, err
	}

	paperObj := new(commercialPaper)
	paperObj.key = key
	err = deserializePaper(paperJSON, paperObj)

	if err != nil {
		return nil, err
	}

	return paperObj, nil
}
