/*
Daytona Server API

Daytona Server API

API version: v0.0.0-dev
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// checks if the Session type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Session{}

// Session struct for Session
type Session struct {
	Commands  []Command `json:"commands"`
	SessionId string    `json:"sessionId"`
}

type _Session Session

// NewSession instantiates a new Session object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSession(commands []Command, sessionId string) *Session {
	this := Session{}
	this.Commands = commands
	this.SessionId = sessionId
	return &this
}

// NewSessionWithDefaults instantiates a new Session object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSessionWithDefaults() *Session {
	this := Session{}
	return &this
}

// GetCommands returns the Commands field value
func (o *Session) GetCommands() []Command {
	if o == nil {
		var ret []Command
		return ret
	}

	return o.Commands
}

// GetCommandsOk returns a tuple with the Commands field value
// and a boolean to check if the value has been set.
func (o *Session) GetCommandsOk() ([]Command, bool) {
	if o == nil {
		return nil, false
	}
	return o.Commands, true
}

// SetCommands sets field value
func (o *Session) SetCommands(v []Command) {
	o.Commands = v
}

// GetSessionId returns the SessionId field value
func (o *Session) GetSessionId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.SessionId
}

// GetSessionIdOk returns a tuple with the SessionId field value
// and a boolean to check if the value has been set.
func (o *Session) GetSessionIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SessionId, true
}

// SetSessionId sets field value
func (o *Session) SetSessionId(v string) {
	o.SessionId = v
}

func (o Session) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Session) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["commands"] = o.Commands
	toSerialize["sessionId"] = o.SessionId
	return toSerialize, nil
}

func (o *Session) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"commands",
		"sessionId",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err
	}

	for _, requiredProperty := range requiredProperties {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varSession := _Session{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varSession)

	if err != nil {
		return err
	}

	*o = Session(varSession)

	return err
}

type NullableSession struct {
	value *Session
	isSet bool
}

func (v NullableSession) Get() *Session {
	return v.value
}

func (v *NullableSession) Set(val *Session) {
	v.value = val
	v.isSet = true
}

func (v NullableSession) IsSet() bool {
	return v.isSet
}

func (v *NullableSession) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSession(val *Session) *NullableSession {
	return &NullableSession{value: val, isSet: true}
}

func (v NullableSession) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSession) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
