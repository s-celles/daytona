# Target

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Default** | **bool** |  | 
**Id** | **string** |  | 
**Name** | **string** |  | 
**TargetConfig** | [**TargetConfig**](TargetConfig.md) |  | 
**TargetConfigName** | **string** |  | 
**Workspaces** | Pointer to [**[]Workspace**](Workspace.md) |  | [optional] 

## Methods

### NewTarget

`func NewTarget(default_ bool, id string, name string, targetConfig TargetConfig, targetConfigName string, ) *Target`

NewTarget instantiates a new Target object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTargetWithDefaults

`func NewTargetWithDefaults() *Target`

NewTargetWithDefaults instantiates a new Target object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDefault

`func (o *Target) GetDefault() bool`

GetDefault returns the Default field if non-nil, zero value otherwise.

### GetDefaultOk

`func (o *Target) GetDefaultOk() (*bool, bool)`

GetDefaultOk returns a tuple with the Default field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefault

`func (o *Target) SetDefault(v bool)`

SetDefault sets Default field to given value.


### GetId

`func (o *Target) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *Target) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *Target) SetId(v string)`

SetId sets Id field to given value.


### GetName

`func (o *Target) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Target) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Target) SetName(v string)`

SetName sets Name field to given value.


### GetTargetConfig

`func (o *Target) GetTargetConfig() TargetConfig`

GetTargetConfig returns the TargetConfig field if non-nil, zero value otherwise.

### GetTargetConfigOk

`func (o *Target) GetTargetConfigOk() (*TargetConfig, bool)`

GetTargetConfigOk returns a tuple with the TargetConfig field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTargetConfig

`func (o *Target) SetTargetConfig(v TargetConfig)`

SetTargetConfig sets TargetConfig field to given value.


### GetTargetConfigName

`func (o *Target) GetTargetConfigName() string`

GetTargetConfigName returns the TargetConfigName field if non-nil, zero value otherwise.

### GetTargetConfigNameOk

`func (o *Target) GetTargetConfigNameOk() (*string, bool)`

GetTargetConfigNameOk returns a tuple with the TargetConfigName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTargetConfigName

`func (o *Target) SetTargetConfigName(v string)`

SetTargetConfigName sets TargetConfigName field to given value.


### GetWorkspaces

`func (o *Target) GetWorkspaces() []Workspace`

GetWorkspaces returns the Workspaces field if non-nil, zero value otherwise.

### GetWorkspacesOk

`func (o *Target) GetWorkspacesOk() (*[]Workspace, bool)`

GetWorkspacesOk returns a tuple with the Workspaces field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWorkspaces

`func (o *Target) SetWorkspaces(v []Workspace)`

SetWorkspaces sets Workspaces field to given value.

### HasWorkspaces

`func (o *Target) HasWorkspaces() bool`

HasWorkspaces returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


