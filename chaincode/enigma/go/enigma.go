/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

/*
 * The sample smart contract for documentation topic:
 * Trade Finance Use Case - WORK IN  PROGRESS
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// Define the vehicle information request
type VehicleInfoRequest struct {
	VehicleNumber	string		`json:"vehicleNumber"`
	ChasisNumber	string		`json:"chasisNumber"`
}

type VehicleResponse struct {
	VehicleNumber	string		`json:"vehicleNumber"`
	ChasisNumber	string		`json:"chasisNumber"`
	InsuranceStatus	string		`json:"insuranceStatus"`
	CurrentAddress	string		`json:"currentAddress"`
	TotalClaims		string		`json:"totalClaims"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {

	// Retrieve the requested Smart Contract function and arguments
	function, args := APIstub.GetFunctionAndParameters()
	// Route to the appropriate handler function to interact with the ledger appropriately
	if function == "requestVehicleInfo" {
		return s.requestVehicleInfo(APIstub, args)
	} else if function == "responseVehicleInfo" {
		return s.responseVehicleInfo(APIstub, args)
	} else if function == "getVehicleHistory" {
		return s.getVehicleHistory(APIstub, args)
	}

	return shim.Error("Invalid Smart Contract function name.")
}

// This function is initiate by RTA 
func (s *SmartContract) requestVehicleInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {	
	vehicleNumber := args[0];
	chasisNumber := args[1];

	VI := VehicleInfoRequest{VehicleNumber: vehicleNumber, ChasisNumber: chasisNumber}
	VIBytes, err := json.Marshal(VI)

	if err != nil {
		return shim.Error("Issue with Request VehicleInformation json marshaling")
	}

    APIstub.PutState(vehicleNumber, VIBytes)
	fmt.Println("VI Requested -> ", VI)

	return shim.Success(nil)
}

// This function is initiate by IRDA
func (s *SmartContract) responseVehicleInfo(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	vehicleNumber := args[0];

	VIAsBytes, _ := APIstub.GetState(vehicleNumber)

	var vi VehicleResponse

	err := json.Unmarshal(VIAsBytes, &vi)

	if err != nil {
		return shim.Error("Issue with VehicleInformation json unmarshaling")
	}

	VI := VehicleResponse{VehicleNumber: vi.VehicleNumber, ChasisNumber: vi.ChasisNumber, InsuranceStatus: vi.InsuranceStatus, CurrentAddress: vi.CurrentAddress, TotalClaims: vi.TotalClaims}

	VIBytes, err := json.Marshal(VI)

	if err != nil {
		return shim.Error("Issue with VehicleInformation json marshaling")
	}

    APIstub.PutState(vi.ChasisNumber,VIBytes)
	fmt.Println("Vehicle Information Provided -> ", VI)

	return shim.Success(nil)
}

// func (s *SmartContract) acceptRequest(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

// 	vehicleNumber := args[0];
// 	VIAsBytes, _ := APIstub.GetState(vehicleNumber)

// 	var lc VehicleInfoRequest

// 	err := json.Unmarshal(VIAsBytes, &lc)

// 	if err != nil {
// 		return shim.Error("Issue with LC json unmarshaling")
// 	}

// 	LC := VehicleInfoRequest{LCId: lc.LCId, ExpiryDate: lc.ExpiryDate, Buyer: lc.Buyer, Bank: lc.Bank, Seller: lc.Seller, Amount: lc.Amount, Status: "Accepted"}
// 	LCBytes, err := json.Marshal(LC)

// 	if err != nil {
// 		return shim.Error("Issue with LC json marshaling")
// 	}

//     APIstub.PutState(lc.LCId,LCBytes)
// 	fmt.Println("LC Accepted -> ", LC)


	

// 	return shim.Success(nil)
// }

func (s *SmartContract) getVehicleHistory(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

	vehicleNumber := args[0];

	resultsIterator, err := APIstub.GetHistoryForKey(vehicleNumber)
	if err != nil {
		return shim.Error("Error retrieving Vehicle history.")
	}
	defer resultsIterator.Close()

	// buffer is a JSON array containing historic values for the marble
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for resultsIterator.HasNext() {
		response, err := resultsIterator.Next()
		if err != nil {
			return shim.Error("Error retrieving Vehicle history.")
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"TxId\":")
		buffer.WriteString("\"")
		buffer.WriteString(response.TxId)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Value\":")
		// if it was a delete operation on given key, then we need to set the
		//corresponding value null. Else, we will write the response.Value
		//as-is (as the Value itself a JSON marble)
		if response.IsDelete {
			buffer.WriteString("null")
		} else {
			buffer.WriteString(string(response.Value))
		}

		buffer.WriteString(", \"Timestamp\":")
		buffer.WriteString("\"")
		buffer.WriteString(time.Unix(response.Timestamp.Seconds, int64(response.Timestamp.Nanos)).String())
		buffer.WriteString("\"")

		buffer.WriteString(", \"IsDelete\":")
		buffer.WriteString("\"")
		buffer.WriteString(strconv.FormatBool(response.IsDelete))
		buffer.WriteString("\"")

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	fmt.Printf("- getVehicleHistory returning:\n%s\n", buffer.String())

	return shim.Success(buffer.Bytes())
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {

	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
