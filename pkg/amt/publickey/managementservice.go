/*********************************************************************
 * Copyright (c) Intel Corporation 2023
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/

package publickey

import (
	"encoding/xml"

	"github.com/open-amt-cloud-toolkit/go-wsman-messages/internal/wsman"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/amt/actions"
	"github.com/open-amt-cloud-toolkit/go-wsman-messages/pkg/cim/models"
)

const AMT_PublicKeyManagementService = "AMT_PublicKeyManagementService"

type Response struct {
	XMLName xml.Name     `xml:"Envelope"`
	Header  wsman.Header `xml:"Header"`
	Body    Body         `xml:"Body"`
}
type Body struct {
	XMLName                          xml.Name                     `xml:"Body"`
	AddTrustedRootCertificate_OUTPUT AddTrustedCertificate_OUTPUT `xml:"AddTrustedRootCertificate_OUTPUT,omitempty"`
	AddTrustedCertificate_OUTPUT     AddTrustedCertificate_OUTPUT `xml:"AddCertificate_OUTPUT,omitempty"`
}
type AddTrustedCertificate_OUTPUT struct {
	CreatedCertificate CreatedCertificate
	ReturnValue        int
}
type AddTrustedRootCertificate_OUTPUT struct {
	CreatedCertificate CreatedCertificate
	ReturnValue        int
}
type CreatedCertificate struct {
	XMLName             xml.Name                          `xml:"CreatedCertificate"`
	Address             string                            `xml:"Address,omitempty"`
	ReferenceParameters models.ReferenceParameters_OUTPUT `xml:"ReferenceParameters,omitempty"`
}
type AddCertificate_INPUT struct {
	XMLName         xml.Name `xml:"h:AddCertificate_INPUT"`
	H               string   `xml:"xmlns:h,attr"`
	CertificateBlob string   `xml:"h:CertificateBlob"`
}
type AddTrustedRootCertificate_INPUT struct {
	XMLName         xml.Name `xml:"h:AddTrustedRootCertificate_INPUT"`
	H               string   `xml:"xmlns:h,attr"`
	CertificateBlob string   `xml:"h:CertificateBlob"`
}

type GenerateKeyPair_INPUT struct {
	XMLName      xml.Name     `xml:"h:GenerateKeyPair_INPUT"`
	H            string       `xml:"xmlns:h,attr"`
	KeyAlgorithm KeyAlgorithm `xml:"h:KeyAlgorithm"`
	KeyLength    KeyLength    `xml:"h:KeyLength"`
}
type KeyAlgorithm int

const (
	RSA KeyAlgorithm = 0
)

type KeyLength int

const (
	KeyLength2048 KeyLength = 2048
)

type PKCS10Request struct {
	XMLName                      xml.Name         `xml:"h:GeneratePKCS10RequestEx_INPUT"`
	H                            string           `xml:"xmlns:h,attr"`
	KeyPair                      string           `xml:"h:KeyPair"`
	NullSignedCertificateRequest string           `xml:"h:NullSignedCertificateRequest"`
	SigningAlgorithm             SigningAlgorithm `xml:"h:SigningAlgorithm"`
}
type SigningAlgorithm int

const (
	SHA1RSA SigningAlgorithm = iota
	SHA256RSA
)

type ManagementService struct {
	base wsman.Base
}

func NewPublicKeyManagementService(wsmanMessageCreator *wsman.WSManMessageCreator) ManagementService {
	return ManagementService{
		base: wsman.NewBase(wsmanMessageCreator, AMT_PublicKeyManagementService),
	}
}
func (PublicKeyManagementService ManagementService) Get() string {
	return PublicKeyManagementService.base.Get(nil)
}
func (PublicKeyManagementService ManagementService) Enumerate() string {
	return PublicKeyManagementService.base.Enumerate()
}
func (PublicKeyManagementService ManagementService) Pull(enumerationContext string) string {
	return PublicKeyManagementService.base.Pull(enumerationContext)
}

func (PublicKeyManagementService ManagementService) Delete(selector *wsman.Selector) string {
	return PublicKeyManagementService.base.Delete(selector)
}
func (p ManagementService) AddCertificate(certificateBlob string) string {
	header := p.base.WSManMessageCreator.CreateHeader(string(actions.AddCertificate), AMT_PublicKeyManagementService, nil, "", "")
	certificate := AddCertificate_INPUT{CertificateBlob: certificateBlob}
	body := p.base.WSManMessageCreator.CreateBody("AddCertificate_INPUT", AMT_PublicKeyManagementService, &certificate)

	return p.base.WSManMessageCreator.CreateXML(header, body)
}

func (p ManagementService) AddTrustedRootCertificate(certificateBlob string) string {
	header := p.base.WSManMessageCreator.CreateHeader(string(actions.AddTrustedRootCertificate), AMT_PublicKeyManagementService, nil, "", "")
	trustedRootCert := AddTrustedRootCertificate_INPUT{CertificateBlob: certificateBlob}
	body := p.base.WSManMessageCreator.CreateBody("AddTrustedRootCertificate_INPUT", AMT_PublicKeyManagementService, &trustedRootCert)

	return p.base.WSManMessageCreator.CreateXML(header, body)
}

func (p ManagementService) GenerateKeyPair(keyPairParameters GenerateKeyPair_INPUT) string {
	header := p.base.WSManMessageCreator.CreateHeader(string(actions.GenerateKeyPair), AMT_PublicKeyManagementService, nil, "", "")
	body := p.base.WSManMessageCreator.CreateBody("GenerateKeyPair_INPUT", AMT_PublicKeyManagementService, &keyPairParameters)

	return p.base.WSManMessageCreator.CreateXML(header, body)
}

func (p ManagementService) GeneratePKCS10RequestEx(pkcs10Request PKCS10Request) string {
	header := p.base.WSManMessageCreator.CreateHeader(string(actions.GeneratePKCS10RequestEx), AMT_PublicKeyManagementService, nil, "", "")
	body := p.base.WSManMessageCreator.CreateBody("GeneratePKCS10RequestEx_INPUT", AMT_PublicKeyManagementService, &pkcs10Request)

	return p.base.WSManMessageCreator.CreateXML(header, body)
}

type AddKeyParameters struct {
	KeyBlob []byte `xml:"h:KeyBlob"`
}

func (p ManagementService) AddKey(keyBlob []byte) string {
	header := p.base.WSManMessageCreator.CreateHeader(string(actions.AddKey), AMT_PublicKeyManagementService, nil, "", "")
	params := &AddKeyParameters{
		KeyBlob: keyBlob,
	}
	body := p.base.WSManMessageCreator.CreateBody("AddKey_INPUT", AMT_PublicKeyManagementService, params)

	return p.base.WSManMessageCreator.CreateXML(header, body)
}
