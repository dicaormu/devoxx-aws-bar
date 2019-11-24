package dao

import (
	"bartenderAsFunction/model"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iot"
	"github.com/aws/aws-sdk-go/service/iotdataplane"
)

type IotConnection struct {
	Iot *iot.IoT
}

var policyDrunkClient = os.Getenv("POLICY_DRUNK_CLIENT")

type IotConnectionInterface interface {
	RegisterDevice(drunkClient *model.DrunkClient) error
	UpdateShadow(idClient, status string) error
}

//TODO 4 Implement this method to register device
func (con *IotConnection) RegisterDevice(drunkClient *model.DrunkClient) error {
	var input iot.CreateKeysAndCertificateInput
	input.SetSetAsActive(true)

	//TODO 4.1 create keys and certificate for the thing using CreateKeysAndCertificate of con.Iot (github.com/aws/aws-sdk-go/service/iot.Iot)
	//output, errCert := con.Iot.CreateKeysAndCertificate(pointer to the input object)
	err := fmt.Errorf("to delete")
	if err != nil {
		return err
	}
	//TODO 4.2 assign the output containing the certificate, private key, public key and pem to the "thing"(drunkClient) to create the thing
	//drunkClient.CertificateArn = *output.CertificateArn
	//drunkClient.PrivateKey = *output.KeyPair.PrivateKey
	//drunkClient.PublicKey = *output.KeyPair.PublicKey
	//drunkClient.CertificatePem = *output.CertificatePem

	//TODO 4.3 explore this function to verify the thing creation
	return createLock(*drunkClient, con)
}

func createLock(drunkClient model.DrunkClient, con *IotConnection) error {
	input := iot.CreateThingInput{ThingName: &drunkClient.IdClient}
	_, err := con.Iot.CreateThing(&input)
	if err != nil {
		fmt.Println(err)
		return err
	}
	var attachPolicyInput iot.AttachPolicyInput
	attachPolicyInput.SetPolicyName(policyDrunkClient)
	attachPolicyInput.SetTarget(drunkClient.CertificateArn)
	_, errAttachPolicy := con.Iot.AttachPolicy(&attachPolicyInput)
	if errAttachPolicy != nil {
		fmt.Println(errAttachPolicy)
		return errAttachPolicy
	}
	if drunkClient.CertificateArn != "" {
		var attachCertToThingInput iot.AttachThingPrincipalInput
		attachCertToThingInput.SetThingName(drunkClient.IdClient)
		attachCertToThingInput.SetPrincipal(drunkClient.CertificateArn)
		_, errAttachPrinc := con.Iot.AttachThingPrincipal(&attachCertToThingInput)
		if errAttachPrinc != nil {
			fmt.Println(errAttachPrinc)
			return errAttachPrinc
		}
	}
	return nil
}

func (con *IotConnection) UpdateShadow(idClient string, desiredStatus string) error {
	var desiredShadow model.ClientObjectState
	desiredShadow.BarStatus = desiredStatus

	shadow := model.IotShadowDoc{
		State: model.IotShadowState{
			Desired: desiredShadow,
		},
	}
	payload, _ := json.Marshal(shadow)
	fmt.Println(payload)
	//TODO create an object iotdataplane.UpdateThingShadowInput and set the thing name and the payload
	input := iotdataplane.UpdateThingShadowInput{}
	//input.SetThingName(idClient)
	//input.SetPayload(payload)
	iotShadowConn := initializeIotDataClient(con.Iot)
	_, err := iotShadowConn.UpdateThingShadow(&input)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			fmt.Println("Error updating shadow")
			fmt.Println(awsErr)
		} else {
			fmt.Println(err)
		}
		return err
	}
	return nil
}

func CreateIotConnection() IotConnectionInterface {
	return &IotConnection{initializeIotClient()}
}

func initializeIotDataClient(iotSvc *iot.IoT) *iotdataplane.IoTDataPlane {
	res, err := iotSvc.DescribeEndpoint(&iot.DescribeEndpointInput{})
	sessionVar, err := session.NewSession(&iotSvc.Config)
	if err != nil {
		fmt.Println("error during aws session initialization :  " + err.Error())
		os.Exit(1)
	}
	return iotdataplane.New(sessionVar, &aws.Config{Endpoint: res.EndpointAddress})
}

func initializeIotClient() *iot.IoT {
	sessionVar, err := session.NewSession(&aws.Config{Region: aws.String("eu-west-1")})
	if err != nil {
		fmt.Println("error during aws session initialization :  ", err)
		os.Exit(1)
	}
	return iot.New(sessionVar, aws.NewConfig().WithLogLevel(aws.LogDebugWithHTTPBody))
}
