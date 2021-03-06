package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	LabelKind = "kind"
)

// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Schema defines a schema of objects with properties
//
// +k8s:openapi-gen=true
type Schema struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata"`

	// Spec the schema specification
	Spec SchemaSpec `yaml:"spec"`
}

// SchemaList contains a list of Schema objects
//
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SchemaList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Schema `json:"items"`
}

// SchemaSpec defines the objects and their properties
type SchemaSpec struct {
	// Objects the list of objects (or kinds) in the schema
	Objects []Object `yaml:"objects"`
}

// FindObject returns the object for the given name or nil
func (s *SchemaSpec) FindObject(name string) *Object {
	for i := range s.Objects {
		o := &s.Objects[i]
		if o.Name == name {
			return o
		}
	}
	return nil
}

// Object defines a type of object with some properties
type Object struct {
	// Name the name of the object kind
	Name string `json:"name" yaml:"name" validate:"nonzero"`

	// Properties the property definitions
	Properties []Property `json:"properties" yaml:"properties"`

	// Mandatory marks this secret as being mandatory to be setup before we can install a cluster
	Mandatory bool `json:"mandatory,omitempty" yaml:"mandatory,omitempty"`
}

// FindProperty returns the property for the given name or nil
func (s *Object) FindProperty(name string) *Property {
	if s == nil {
		return nil
	}
	for i := range s.Properties {
		p := &s.Properties[i]
		if p.Name == name {
			return p
		}
	}
	return nil
}

// Property defines a property in an object
type Property struct {
	// Name the name of the property
	Name string `json:"name" yaml:"name" validate:"nonzero"`

	// Question the main prompt generated in a user interface when asking to populate the property
	Question string `json:"question" yaml:"question" validate:"nonzero"`

	// Help the tooltip or help text for this property
	Help string `json:"help,omitempty" yaml:"help"`

	// DefaultValue is used to specify default values populated on startup
	DefaultValue string `json:"defaultValue,omitempty" yaml:"defaultValue,omitempty"`

	// Pattern is a regular expression pattern used for validation
	Pattern string `json:"pattern,omitempty" yaml:"pattern,omitempty"`

	// Requires specifies a requirements expression
	Requires string `json:"requires,omitempty" yaml:"requires,omitempty"`

	// Format the format of the value
	Format string `json:"format,omitempty" yaml:"format,omitempty"`

	// Generator the name of the generator to use to create values
	// if this value is non zero we assume Generate is effectively true
	Generator string `json:"generator,omitempty" yaml:"generator,omitempty"`

	// Template the go template used to generate the value of this secret
	// if we need to combine multiple secret values together into a composite secret value.
	//
	// For example if we want to create a maven-settings.xml file or a docker config JSON
	// document made up of lots of static text but some real secrets embedded we can
	// define the template in the schema
	Template string `json:"template,omitempty" yaml:"template,omitempty"`

	// OnlyTemplateIfBlank if this is true then lets only regenerate a template value if the current value is empty
	OnlyTemplateIfBlank bool `json:"onlyTemplateIfBlank,omitempty" yaml:"onlyTemplateIfBlank,omitempty"`

	// Retry enable a retry loop if a template does not evaluate correctly first time
	Retry bool `json:"retry,omitempty" yaml:"retry,omitempty"`

	// Labels allows arbitrary metadata labels to be associated with the property
	Labels map[string]string `json:"labels,omitempty" yaml:"labels,omitempty"`

	// MinLength the minimum number of characters in the value
	MinLength int `json:"minLength,omitempty" yaml:"minLength,omitempty"`

	// MaxLength the maximum number of characters in the value
	MaxLength int `json:"maxLength,omitempty" yaml:"maxLength,omitempty"`

	// NoMask whether to exclude from Secret masking in logs
	NoMask bool `json:"noMask,omitempty" yaml:"mask,omitempty"`
}
