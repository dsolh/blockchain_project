package main

import (
	"encoding/json"
	"fmt"
	
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type WizardStruct struct {
	Name string `json:name`
	Score int `json:score`
	Spell string `json:spell`
}
type SpellStruct struct {
	Spell string `json:spell`
	Category string `json:category`
}

func (s *SmartContract) WizardExists(ctx contractapi.TransactionContextInterface, wizardName string) (bool,error) {
	// if there is a data that has wizardName as a key, then return true, else then false
	wizardBytes, err := ctx.GetStub().GetState(wizardName)
	if err != nil {
		return false, err
	}

	return wizardBytes != nil, nil
}
func (s *SmartContract) SpellExists(ctx contractapi.TransactionContextInterface, spell string) (bool,error) {
	// if there is a data that has spell as a key, then return true, else then false
	spellBytes, err := ctx.GetStub().GetState(spell)
	if err != nil {
		return false, err
	}

	return spellBytes != nil, nil
}

func (s *SmartContract) RegisterWizard(ctx contractapi.TransactionContextInterface, wizardName string, country string) error {
	//var countries = ["British", "USA", "France", "German", "China"]
	
	// check if the wizard already exists
	exists, err := s.WizardExists(ctx, wizardName)
	if err != nil {
		return err
	}
	if exists != false {
		return fmt.Errorf("(%s) is already registered", wizardName)
	}
	
	wizardStruct := WizardStruct {
		Name:wizardName,
		Score:0,
	}

	wizardBytes, err := json.Marshal(wizardStruct)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(wizardName, wizardBytes)
}

func (s *SmartContract) CreateSpell(ctx contractapi.TransactionContextInterface, spell string, category string) error {
	// check if the spell already exists
	exists, err := s.SpellExists(ctx, spell)
	if err != nil {
		return err
	}
	if exists != false {
		return fmt.Errorf("spell \"(%s)\" is already existed", spell )
	}

	spellStruct := SpellStruct {
		Spell:spell,
		Category:category,
	}

	spellBytes, err := json.Marshal(spellStruct)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(spell, spellBytes)
}

func (s *SmartContract) ReadWizard(ctx contractapi.TransactionContextInterface, wizardName string) (*WizardStruct, error) {
	wizardBytes, err := ctx.GetStub().GetState(wizardName)
	if err != nil {
		return nil, err
	} else if wizardBytes == nil{
		return nil, fmt.Errorf("(%s) is not registered", wizardName)
	}

	wizardStruct := WizardStruct{}
	err = json.Unmarshal(wizardBytes, &wizardStruct)
	if err != nil {
		return nil, err
	}

	return &wizardStruct, nil
}

func (s *SmartContract) ReadSpell(ctx contractapi.TransactionContextInterface, spell string) (*SpellStruct, error) {
	spellBytes, err := ctx.GetStub().GetState(spell)
	if err != nil {
		return nil, err
	} else if spellBytes == nil{
		return nil, fmt.Errorf("spell \"(%s)\" is not existed", spell )
	}

	spellStruct := SpellStruct{}
	err = json.Unmarshal(spellBytes, &spellStruct)
	if err != nil {
		return nil, err
	}

	return &spellStruct, nil
}

func (s *SmartContract) DeleteWizard(ctx contractapi.TransactionContextInterface, wizardName string) error {
	exists, err := s.WizardExists(ctx, wizardName)
	if err != nil {
		return err
	}
	if exists == false {
		return fmt.Errorf("(%s) is not existed", wizardName)
	}

	return ctx.GetStub().DelState(wizardName)
}

func (s *SmartContract) CastSpell(ctx contractapi.TransactionContextInterface, wizardName string, spell string) error {
	// get wizard struct
	wizardStruct, err := s.ReadWizard(ctx, wizardName)

	// get spell struct
	spellStruct, err := s.ReadSpell(ctx, spell)

	// calculate score
	var scr = wizardStruct.Score
	var ctg = spellStruct.Category
	switch ctg {
	case "unforgivable_curse":
	   scr = scr - 1000
	case "curse":
	   scr = scr - 100
	case "jinx":
	   scr = scr - 50
	// 카테고리 업데이트 필요
	}

	// update info
	wizardStruct.Score = scr
	wizardStruct.Spell = spell

	wizardBytes, err := json.Marshal(wizardStruct)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(wizardName, wizardBytes) 
}

func main() {
	// peer가 chaincode를 실행시킬 수 있도록 chaincode 객체를 peer에 넘겨주는거임. 그 다음 Start() 함수로 객체를 시작하는거고.
	
	chaincode, err := contractapi.NewChaincode(new(SmartContract))

	if err != nil {
		fmt.Printf("Error create whozard chaincode: %s", err.Error())
		return
	}

	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting whozard chaincode: %s", err.Error())
	}
}
