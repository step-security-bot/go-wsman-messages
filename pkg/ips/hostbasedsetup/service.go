/*********************************************************************
 * Copyright (c) Intel Corporation 2023
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/

package hostbasedsetup

import (
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/internal/wsman"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/ips/actions"
)

type Response struct {
	XMLName xml.Name     `xml:"Envelope"`
	Header  wsman.Header `xml:"Header"`
	Body    Body         `xml:"Body"`
}

type Body struct {
	XMLName                   xml.Name                  `xml:"Body"`
	Setup_OUTPUT              Setup_OUTPUT              `xml:"Setup_OUTPUT"`
	AdminSetup_OUTPUT         AdminSetup_OUTPUT         `xml:"AdminSetup_OUTPUT"`
	AddNextCertInChain_OUTPUT AddNextCertInChain_OUTPUT `xml:"AddNextCertInChain_OUTPUT"`
	IPS_HostBasedSetupService HostBasedSetupService     `xml:"IPS_HostBasedSetupService"`
}

type HostBasedSetupService struct {
	XMLName                 xml.Name `xml:"IPS_HostBasedSetupService"`
	ElementName             string
	SystemCreationClassName string
	SystemName              string
	CreationClassName       string
	Name                    string
	CurrentControlMode      int
	AllowedControlModes     int
	ConfigurationNonce      string
	CertChainStatus         int
}

type AdminPassEncryptionType int

const (
	AdminPassEncryptionTypeNone AdminPassEncryptionType = iota
	AdminPassEncryptionTypeOther
	AdminPassEncryptionTypeHTTPDigestMD5A1
)

type SigningAlgorithm int

const (
	SigningAlgorithmNone SigningAlgorithm = iota
	SigningAlgorithmOther
	SigningAlgorithmRSASHA2256
)

type Service struct {
	base wsman.Base
}

const IPS_HostBasedSetupService = "IPS_HostBasedSetupService"

// NewHostBasedSetupService returns a new instance of the HostBasedSetupService struct.
func NewHostBasedSetupService(wsmanMessageCreator *wsman.WSManMessageCreator) Service {
	return Service{
		base: wsman.NewBase(wsmanMessageCreator, string(IPS_HostBasedSetupService)),
	}
}

// Get retrieves the representation of the instance
func (b Service) Get() string {
	return b.base.Get(nil)
}

// Enumerates the instances of this class
func (b Service) Enumerate() string {
	return b.base.Enumerate()
}

// Pulls instances of this class, following an Enumerate operation
func (b Service) Pull(enumerationContext string) string {
	return b.base.Pull(enumerationContext)
}

type AddNextCertInChain struct {
	XMLName           xml.Name `xml:"h:AddNextCertInChain_INPUT"`
	H                 string   `xml:"xmlns:h,attr"`
	NextCertificate   string   `xml:"h:NextCertificate"`
	IsLeafCertificate bool     `xml:"h:IsLeafCertificate"`
	IsRootCertificate bool     `xml:"h:IsRootCertificate"`
}

type AddNextCertInChain_OUTPUT struct {
	ReturnValue int
}

// Add a certificate to the provisioning certificate chain, to be used by AdminSetup or UpgradeClientToAdmin methods.
func (b Service) AddNextCertInChain(cert string, isLeaf bool, isRoot bool) string {
	header := b.base.WSManMessageCreator.CreateHeader(string(actions.AddNextCertInChain), string(IPS_HostBasedSetupService), nil, "", "")
	body := b.base.WSManMessageCreator.CreateBody("AddNextCertInChain_INPUT", string(IPS_HostBasedSetupService), AddNextCertInChain{
		H:                 "http://intel.com/wbem/wscim/1/ips-schema/1/IPS_HostBasedSetupService",
		NextCertificate:   cert,
		IsLeafCertificate: isLeaf,
		IsRootCertificate: isRoot,
	})
	return b.base.WSManMessageCreator.CreateXML(header, body)
}

type AdminSetup struct {
	XMLName                    xml.Name `xml:"h:AdminSetup_INPUT"`
	H                          string   `xml:"xmlns:h,attr"`
	NetAdminPassEncryptionType int      `xml:"h:NetAdminPassEncryptionType"`
	DigestRealm                string	`xml:"h:DigestRealm"`
	NetworkAdminPassword       string   `xml:"h:NetworkAdminPassword"`
	McNonce                    string   `xml:"h:McNonce"`
	SigningAlgorithm           int      `xml:"h:SigningAlgorithm"`
	DigitalSignature           string   `xml:"h:DigitalSignature"`
}

type AdminSetup_OUTPUT struct {
	ReturnValue int
}

// Setup Intel(R) AMT from the local host, resulting in Admin Setup Mode. Requires OS administrator rights, and moves Intel(R) AMT from "Pre Provisioned" state to "Post Provisioned" state. The control mode after this method is run will be "Admin".
func (b Service) AdminSetup(adminPassEncryptionType AdminPassEncryptionType, digestRealm string, adminPassword string, mcNonce string, signingAlgorithm SigningAlgorithm, digitalSignature string) string {
	hashInHex := createMD5Hash(adminPassword, digestRealm)
	header := b.base.WSManMessageCreator.CreateHeader(string(actions.AdminSetup), string(IPS_HostBasedSetupService), nil, "", "")
	body := b.base.WSManMessageCreator.CreateBody("AdminSetup_INPUT", string(IPS_HostBasedSetupService), AdminSetup{
		H:                          "http://intel.com/wbem/wscim/1/ips-schema/1/IPS_HostBasedSetupService",
		NetAdminPassEncryptionType: int(adminPassEncryptionType),
		NetworkAdminPassword:       string(hashInHex),
		McNonce:                    mcNonce,
		SigningAlgorithm:           int(signingAlgorithm),
		DigitalSignature:           digitalSignature,
	})
	return b.base.WSManMessageCreator.CreateXML(header, body)
}

type Setup struct {
	XMLName                    xml.Name `xml:"h:Setup_INPUT"`
	H                          string   `xml:"xmlns:h,attr"`
	NetAdminPassEncryptionType int      `xml:"h:NetAdminPassEncryptionType"`
	NetworkAdminPassword       string   `xml:"h:NetworkAdminPassword"`
}
type Setup_OUTPUT struct {
	ReturnValue int
}

func (b Service) Setup(adminPassEncryptionType AdminPassEncryptionType, digestRealm, adminPassword string) string {
	hashInHex := createMD5Hash(adminPassword, digestRealm)
	header := b.base.WSManMessageCreator.CreateHeader(string(actions.Setup), string(IPS_HostBasedSetupService), nil, "", "")
	body := b.base.WSManMessageCreator.CreateBody("Setup_INPUT", string(IPS_HostBasedSetupService), Setup{
		H:                          "http://intel.com/wbem/wscim/1/ips-schema/1/IPS_HostBasedSetupService",
		NetAdminPassEncryptionType: int(adminPassEncryptionType),
		NetworkAdminPassword:       string(hashInHex),
	})
	return b.base.WSManMessageCreator.CreateXML(header, body)
}

func createMD5Hash(adminPassword string, digestRealm string) string {
	// Create an md5 hash.
	setupPassword := "admin:" + digestRealm + ":" + adminPassword
	hash := md5.New()
	_, _ = io.WriteString(hash, setupPassword)
	hashInHex := fmt.Sprintf("%x", hash.Sum(nil))
	return hashInHex
}

type UpgradeClientToAdmin struct {
	XMLName          xml.Name `xml:"h:UpgradeClientToAdmin_INPUT"`
	H                string   `xml:"xmlns:h,attr"`
	McNonce          string   `xml:"h:McNonce"`
	SigningAlgorithm int      `xml:"h:SigningAlgorithm"`
	DigitalSignature string   `xml:"h:DigitalSignature"`
}

// Upgrade Intel(R) AMT from Client to Admin Control Mode.
func (b Service) UpgradeClientToAdmin(mcNonce string, signingAlgorithm SigningAlgorithm, digitalSignature string) string {
	header := b.base.WSManMessageCreator.CreateHeader(string(actions.UpgradeClientToAdmin), string(IPS_HostBasedSetupService), nil, "", "")
	body := b.base.WSManMessageCreator.CreateBody("UpgradeClientToAdmin_INPUT", string(IPS_HostBasedSetupService), UpgradeClientToAdmin{
		H:                "http://intel.com/wbem/wscim/1/ips-schema/1/IPS_HostBasedSetupService",
		McNonce:          mcNonce,
		SigningAlgorithm: int(signingAlgorithm),
		DigitalSignature: digitalSignature,
	})
	return b.base.WSManMessageCreator.CreateXML(header, body)
}
