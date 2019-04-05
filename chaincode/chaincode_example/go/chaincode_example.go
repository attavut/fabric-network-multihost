/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

//WARNING - this chaincode's ID is hard-coded in chaincode_example04 to illustrate one way of
//calling chaincode from a chaincode. If this example is modified, chaincode_example04.go has
//to be modified as well with the new ID of chaincode_example.
//chaincode_example05 show's how chaincode ID can be passed in as a parameter instead of
//hard-coding.

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type ConsentLog struct {
	MobileID     string `json:"mobileid"`
	MobileNumber string `json:"mobilenumber"`
	AAL          string `json:"aal"`
	MobileIDIAL  string `json:"mobileidial"`
	IssuerCode   string `json:"issuercode"`
	IssuerName   string `json:"issuername"`
	VerifierCode string `json:"verifiercode"`
	VerifierName string `json:"verifiername"`
	ConsentDate  string `json:"consentdate"`
	TxCode       string `json:"txcode"`
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	function, args := stub.GetFunctionAndParameters()
	if function == "addConsent" {
		// Add MobileID Consent log
		return t.addConsent(stub, args)
	} else if function == "queryConsent" {
		// Query MobileID Consent log
		return t.queryConsent(stub, args)
	} else if function == "queryAllConsents" {
		// Query All MobileID Consents
		return t.queryAllConsents(stub, args)
	}

	return shim.Error("Invalid invoke function name. Expecting \"invoke\" \"delete\" \"query\"")
}

// Add MobileID Consent log
func (t *SimpleChaincode) addConsent(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var err error
	var mobileID string

	// if len(args) != 3 {
	// 	return shim.Error("Incorrect number of arguments. Expecting 3")
	// }

	mobileID = args[0]
	var consentLog = ConsentLog{MobileID: args[0], MobileNumber: args[1], AAL: args[2], MobileIDIAL: args[3], IssuerCode: args[4], IssuerName: args[5], VerifierCode: args[6], VerifierName: args[7], ConsentDate: args[8], TxCode: args[9]}

	consentLogAsBytes, _ := json.Marshal(consentLog)
	// Write the state back to the ledger
	err = stub.PutState("CONSENT_"+mobileID, consentLogAsBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success(nil)
}

// Query MobileID Consent log
func (t *SimpleChaincode) queryConsent(stub shim.ChaincodeStubInterface, args []string) pb.Response {

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	consentLogAsBytes, _ := stub.GetState(args[0])
	return shim.Success(consentLogAsBytes)
}

func (t *SimpleChaincode) queryAllConsents(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	startKey := "CONSENT_0"
	endKey := "CONSENT_9999999999"

	resultsIterator, err := stub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- queryAllConsents:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
