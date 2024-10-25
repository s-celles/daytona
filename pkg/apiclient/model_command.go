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

// checks if the Command type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Command{}

// Command struct for Command
type Command struct {
	Command  string `json:"command"`
	ExitCode *int32 `json:"exitCode,omitempty"`
	Id       string `json:"id"`
}

type _Command Command

// NewCommand instantiates a new Command object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCommand(command string, id string) *Command {
	this := Command{}
	this.Command = command
	this.Id = id
	return &this
}

// NewCommandWithDefaults instantiates a new Command object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCommandWithDefaults() *Command {
	this := Command{}
	return &this
}

// GetCommand returns the Command field value
func (o *Command) GetCommand() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Command
}

// GetCommandOk returns a tuple with the Command field value
// and a boolean to check if the value has been set.
func (o *Command) GetCommandOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Command, true
}

// SetCommand sets field value
func (o *Command) SetCommand(v string) {
	o.Command = v
}

// GetExitCode returns the ExitCode field value if set, zero value otherwise.
func (o *Command) GetExitCode() int32 {
	if o == nil || IsNil(o.ExitCode) {
		var ret int32
		return ret
	}
	return *o.ExitCode
}

// GetExitCodeOk returns a tuple with the ExitCode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Command) GetExitCodeOk() (*int32, bool) {
	if o == nil || IsNil(o.ExitCode) {
		return nil, false
	}
	return o.ExitCode, true
}

// HasExitCode returns a boolean if a field has been set.
func (o *Command) HasExitCode() bool {
	if o != nil && !IsNil(o.ExitCode) {
		return true
	}

	return false
}

// SetExitCode gets a reference to the given int32 and assigns it to the ExitCode field.
func (o *Command) SetExitCode(v int32) {
	o.ExitCode = &v
}

// GetId returns the Id field value
func (o *Command) GetId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Id
}

// GetIdOk returns a tuple with the Id field value
// and a boolean to check if the value has been set.
func (o *Command) GetIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Id, true
}

// SetId sets field value
func (o *Command) SetId(v string) {
	o.Id = v
}

func (o Command) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Command) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["command"] = o.Command
	if !IsNil(o.ExitCode) {
		toSerialize["exitCode"] = o.ExitCode
	}
	toSerialize["id"] = o.Id
	return toSerialize, nil
}

func (o *Command) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"command",
		"id",
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

	varCommand := _Command{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varCommand)

	if err != nil {
		return err
	}

	*o = Command(varCommand)

	return err
}

type NullableCommand struct {
	value *Command
	isSet bool
}

func (v NullableCommand) Get() *Command {
	return v.value
}

func (v *NullableCommand) Set(val *Command) {
	v.value = val
	v.isSet = true
}

func (v NullableCommand) IsSet() bool {
	return v.isSet
}

func (v *NullableCommand) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCommand(val *Command) *NullableCommand {
	return &NullableCommand{value: val, isSet: true}
}

func (v NullableCommand) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCommand) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}