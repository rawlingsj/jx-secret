package secretfacade

import (
	v1 "github.com/jenkins-x/jx-secret/pkg/apis/external/v1"
	schema "github.com/jenkins-x/jx-secret/pkg/apis/schema/v1alpha1"
	"github.com/jenkins-x/jx-secret/pkg/extsecrets"
	"github.com/jenkins-x/jx-secret/pkg/schemas"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

// Options options for verifying secrets
type Options struct {
	SecretClient extsecrets.Interface
	KubeClient   kubernetes.Interface
	Namespace    string

	// ExternalSecrets the loaded secrets
	ExternalSecrets []*v1.ExternalSecret
}

// SecretError returns an error for a secret
type SecretError struct {
	// ExternalSecret the external secret which is not valid
	ExternalSecret v1.ExternalSecret

	// EntryErrors the errors for each secret entry
	EntryErrors []*EntryError
}

// EntryError represents the missing entries
type EntryError struct {
	// Key the secret key
	Key string

	// Properties property names for the key
	Properties []string
}

// SecretPair the external secret and the associated Secret an error for a secret
type SecretPair struct {
	// ExternalSecret the external secret which is not valid
	ExternalSecret v1.ExternalSecret

	// Secret the secret if there is one
	Secret *corev1.Secret

	// Error last validation error at last check
	Error *SecretError

	// schemaObject caches the schema object
	schemaObject *schema.Object
}

// IsInvalid returns true if the validation failed
func (p *SecretPair) IsInvalid() bool {
	return p.Error != nil && len(p.Error.EntryErrors) > 0
}

// IsMandatory returns true if the secret is a mandatory secret
func (p *SecretPair) IsMandatory() bool {
	obj, err := p.SchemaObject()
	if err == nil && obj != nil {
		return obj.Mandatory
	}
	return false
}

// SchemaObject returns the optional schema object from the annotation
func (p *SecretPair) SchemaObject() (*schema.Object, error) {
	if p.schemaObject != nil {
		return p.schemaObject, nil
	}
	ann := p.ExternalSecret.Annotations
	if ann == nil {
		return nil, nil
	}
	text := ann[extsecrets.SchemaObjectAnnotation]
	if text == "" {
		return nil, nil
	}
	var err error
	p.schemaObject, err = schemas.ObjectFromString(text)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load schema object from ExternalSecret annotation for %s", p.ExternalSecret.Name)
	}
	return p.schemaObject, nil
}
