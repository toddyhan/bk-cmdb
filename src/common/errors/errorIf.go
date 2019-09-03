/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package errors

// CCErrorCoder get Error Code
type CCErrorCoder interface {
	error
	// GetCode return the error code
	GetCode() int
}

// DefaultCCErrorIf defines default error code interface
type DefaultCCErrorIf interface {
	// Error returns an error with error code
	Error(errCode int) error
	// Errorf returns an error with error code
	Errorf(errCode int, args ...interface{}) error

	// CCError returns an error with error code
	CCError(errCode int) CCErrorCoder
	// CCErrorf returns an error with error code
	CCErrorf(errCode int, args ...interface{}) CCErrorCoder

	// New create a new error with error code and message
	New(errorCode int, msg string) error
}

// CCErrorIf defines error information conversion
type CCErrorIf interface {
	// CreateDefaultCCErrorIf create new language error interface instance
	CreateDefaultCCErrorIf(language string) DefaultCCErrorIf
	// Error returns an error for specific language
	Error(language string, errCode int) error
	// Errorf Errorf returns an error with args for specific language
	Errorf(language string, errCode int, args ...interface{}) error

	Load(res map[string]ErrorCode)
}

func New(errCode int, msg string) CCErrorCoder {
	return &ccError{
		code: errCode,
		callback: func() string {
			return msg
		},
	}
}
