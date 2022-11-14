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
	// instead of setting multichennel
	Country string `json:country`
}
type SpellStruct struct {
	Spell string `json:spell`
	Category string `json:category`
}

// check existance
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

//read
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

// delete
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

// main function
func (s *SmartContract) RegisterWizard(ctx contractapi.TransactionContextInterface, wizardName string, country string) error {
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
		Spell:"Didn't cast anything",
		// instead of setting multichennel
		Country:country,
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
		return fmt.Errorf("spell \"%s\" is already existed", spell )
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

func (s *SmartContract) CastSpell(ctx contractapi.TransactionContextInterface, wizardName string, spell string, country string) error {
	// get wizard struct
	wizardStruct, err := s.ReadWizard(ctx, wizardName)
	if err != nil {
		return err
	}
	
	// instead of setting multichennel
	// check country
	if wizardStruct.Country != country {
		return fmt.Errorf("country doesn't correct")
	}

	// get spell struct
	spellStruct, err := s.ReadSpell(ctx, spell)
	if err != nil {
		return err
	}

	// calculate score
	var scr = wizardStruct.Score
	var ctg = spellStruct.Category
	switch ctg {
	case "unforgivable curse":
		scr = scr - 1000
	case "curse":
		scr = scr - 100
	case "hex" :
		scr = scr - 30
	case "jinx":
		scr = scr - 30
	case "transfiguration":
		scr = scr + 10
	case "healing spell":
		scr = scr + 50
	case "counter-jinxes":
		scr = scr + 100
	case "counter-curses":
		scr = scr + 500
	case "charm":
		scr = scr
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

func (s *SmartContract) CheckIdentity(ctx contractapi.TransactionContextInterface, wizardName string, country string) (string, error) {
	// get wizard struct
	wizardStruct, err := s.ReadWizard(ctx, wizardName)
	if err != nil {
		return "", err
	}

	// instead of setting multichennel
	// check country
	if wizardStruct.Country != country {
		return "", fmt.Errorf("country doesn't correct")
	}

	// check identity
	scr := wizardStruct.Score
	identity := ""
	if scr <= -500 {
		identity = "criminal"
	} else if scr <= -100 {
		identity = "penalty"
	} else if scr < 0 {
		identity = "warning"
	} else if scr <= 200 {
		identity = "normal"
	} else if scr <= 500 {
		identity = "respectable"
	} else if scr <= 1000 {
		identity = "extraordinary"
	} else {
		identity = "grand"
	}
	
	return identity, err
}

func (s *SmartContract) RecentSpell(ctx contractapi.TransactionContextInterface, wizardName string, country string) (string, error) {
	// get wizard struct
	wizardStruct, err := s.ReadWizard(ctx, wizardName)
	if err != nil {
		return "", err
	}

	// instead of setting multichennel
	// check country
	if wizardStruct.Country != country {
		return "", fmt.Errorf("country doesn't correct")
	}

	return wizardStruct.Spell, err
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
