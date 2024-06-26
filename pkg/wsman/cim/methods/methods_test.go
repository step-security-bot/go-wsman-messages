/*********************************************************************
 * Copyright (c) Intel Corporation 2023
 * SPDX-License-Identifier: Apache-2.0
 **********************************************************************/

package methods

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TestMethod string = "TestMethod"
	CIMTest    string = "CIM_Test"
)

func TestMethods(t *testing.T) {
	t.Run("GenerateAction Test", func(t *testing.T) {
		expectedResult := "http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_Test/TestMethod"
		className := CIMTest
		methodName := TestMethod
		result := GenerateAction(className, methodName)
		assert.Equal(t, expectedResult, result)
	})
	t.Run("GenerateMethod Test", func(t *testing.T) {
		expectedResult := "TestMethod_INPUT"
		methodName := TestMethod
		result := GenerateInputMethod(methodName)
		assert.Equal(t, expectedResult, result)
	})
	t.Run("RequestStateChange Test", func(t *testing.T) {
		expectedResult := "http://schemas.dmtf.org/wbem/wscim/1/cim-schema/2/CIM_Test/RequestStateChange"
		className := CIMTest
		result := RequestStateChange(className)
		assert.Equal(t, expectedResult, result)
	})
}
