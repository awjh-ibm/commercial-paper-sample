package main

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/hyperledger/fabric/contractapi"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

var logger = shim.NewLogger("commercial_paper")

type market struct {
	key          string
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	ListedPapers []listing `json:"listedPapers"`
}

type listing struct {
	Paper    *commercialPaper `json:"paper"`
	Discount int              `json:"discount"`
}

type commercialPaper struct {
	key      string
	ID       string `json:"id"`
	Maturity int    `json:"maturity"`
	Par      int    `json:"par"`
}

// CommercialPaperContract - contains the business rules for commercial papers
type CommercialPaperContract struct {
	contractapi.Contract
}

// Setup - creates the initial market
func (cpc *CommercialPaperContract) Setup(ctx contractapi.TransactionContext) error {
	defaultName := "US_BLUE_ONE"

	key, err := assetIDToKey(ctx, "market", defaultName)

	if err != nil {
		return err
	}

	marketData := market{key, defaultName, fmt.Sprintf("%s trading", defaultName), []listing{}}

	serializedMarket, err := serialize(marketData)

	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(key, serializedMarket)

	return err
}

// CreatePaper - creates a new commercial paper and stores it in the world state
func (cpc *CommercialPaperContract) CreatePaper(ctx contractapi.TransactionContext, CUSIP string, maturity string, par string) error {
	maturityInt, err := strconv.Atoi(maturity)

	if err != nil {
		return fmt.Errorf("Maturity should be an integer. Passed maturity was %s", maturity)
	}

	if maturityInt < 1 || maturityInt > 364 {
		return fmt.Errorf("Maturity must be between 1 and 364. Passed maturity was %d", maturityInt)
	}

	parInt, err := strconv.Atoi(par)

	if err != nil {
		return fmt.Errorf("Par should be an integer. Passed par was %s", par)
	}

	if parInt < 1 {
		return fmt.Errorf("Par must be greater than 1. Passed par was ")
	}

	key, err := assetIDToKey(ctx, "paper", CUSIP)

	if err != nil {
		return err
	}

	paper := commercialPaper{key, CUSIP, maturityInt, parInt}

	serializedPaper, err := serialize(paper)

	if err != nil {
		return err
	}

	err = ctx.GetStub().PutState(key, serializedPaper)

	return nil
}

// ListOnMarket - Adds the papers passed to the market passed and marks their discount
func (cpc *CommercialPaperContract) ListOnMarket(ctx contractapi.TransactionContext, marketID string, discount string, papersToList []string) error {
	stub := ctx.GetStub()

	discountInt, err := strconv.Atoi(discount)

	if err != nil {
		return fmt.Errorf("Discount should be an integer. Passed discount was %s", discount)
	}

	marketObj, err := getMarket(ctx, marketID)

	if err != nil {
		return err
	}

	for _, paper := range papersToList {
		logger.Info(paper)

		paperObj, err := getPaper(ctx, paper)

		if err != nil {
			return err
		}

		listing := listing{paperObj, discountInt}

		marketObj.ListedPapers = append(marketObj.ListedPapers, listing)
	}

	newMarketData, err := serialize(marketObj)

	if err != nil {
		return err
	}

	err = stub.PutState(marketObj.key, newMarketData)

	if err != nil {
		return errors.New("Failed to update market")
	}

	return nil
}

// RetrieveMarket - retrieves the current state of a stored market in JSON format
func (cpc *CommercialPaperContract) RetrieveMarket(ctx contractapi.TransactionContext, marketID string) (string, error) {
	_, marketJSON, err := get(ctx, "market", marketID)

	if err != nil {
		return "", err
	}

	return string(marketJSON), nil
}

func main() {
	cpc := new(CommercialPaperContract)

	if err := contractapi.CreateNewChaincode(cpc); err != nil {
		fmt.Printf("Error starting SimpleAsset chaincode: %s", err)
	}
}
