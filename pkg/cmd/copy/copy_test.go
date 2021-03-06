package copy_test

import (
	"context"
	"testing"

	"github.com/jenkins-x-plugins/jx-secret/pkg/cmd/copy"
	"github.com/jenkins-x/jx-helpers/v3/pkg/testhelpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
)

var (
	ns = "jx"
)

func TestCopyBySelector(t *testing.T) {
	kubeObjects := []runtime.Object{
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "knative-docker-user-pass",
				Namespace: ns,
				Labels: map[string]string{
					"beer": "stella",
				},
			},
			Data: map[string][]byte{
				"username": []byte("dummyDockerUsername"),
				"password": []byte("dummyDockerPassword"),
			},
		},
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "lighthouse-oauth-token",
				Namespace: ns,
				Labels: map[string]string{
					"beer": "stella",
				},
			},
			Data: map[string][]byte{
				"oauth": []byte("dummyPipelineUserToken"),
			},
		},
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "lighthouse-hmac-token",
				Namespace: ns,
			},
			Data: map[string][]byte{
				"oauth": []byte("dummyPipelineUserToken"),
			},
		},
	}
	_, o := copy.NewCmdCopy()
	o.Namespace = ns
	o.KubeClient = fake.NewSimpleClientset(kubeObjects...)

	o.ToNamespace = "my-preview-env"
	o.Selector = "beer=stella"
	o.CreateNamespace = true
	err := o.Run()
	require.NoError(t, err, "failed to run export")

	secret, message := testhelpers.RequireSecretExists(t, o.KubeClient, o.ToNamespace, "lighthouse-oauth-token")
	testhelpers.AssertSecretEntryEquals(t, secret, "oauth", "dummyPipelineUserToken", message)

	tons, err := o.KubeClient.CoreV1().Namespaces().Get(context.TODO(), o.ToNamespace, metav1.GetOptions{})
	require.NoError(t, err, "should have found to namespace %s", o.ToNamespace)
	require.NotNil(t, tons, "should have to namespace %s", o.ToNamespace)
	assert.Equal(t, o.ToNamespace, tons.Name, "toNamespace.Name")
}

func TestCopyByName(t *testing.T) {
	kubeObjects := []runtime.Object{
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "lighthouse-oauth-token",
				Namespace: ns,
				Labels: map[string]string{
					"beer": "stella",
				},
			},
			Data: map[string][]byte{
				"oauth": []byte("dummyPipelineUserToken"),
			},
		},
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "knative-docker-user-pass",
				Namespace: ns,
				Labels: map[string]string{
					"beer": "stella",
				},
			},
			Data: map[string][]byte{
				"username": []byte("dummyDockerUsername"),
				"password": []byte("dummyDockerPassword"),
			},
		},
		&corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "lighthouse-hmac-token",
				Namespace: ns,
			},
			Data: map[string][]byte{
				"oauth": []byte("dummyPipelineUserToken"),
			},
		},
	}
	_, o := copy.NewCmdCopy()
	o.Namespace = ns
	o.KubeClient = fake.NewSimpleClientset(kubeObjects...)
	o.ToNamespace = "my-preview-env"
	o.Name = "lighthouse-oauth-token"
	o.CreateNamespace = true
	err := o.Run()
	require.NoError(t, err, "failed to run export")

	secret, message := testhelpers.RequireSecretExists(t, o.KubeClient, o.ToNamespace, "lighthouse-oauth-token")
	testhelpers.AssertSecretEntryEquals(t, secret, "oauth", "dummyPipelineUserToken", message)

	ctx := context.TODO()
	tons, err := o.KubeClient.CoreV1().Namespaces().Get(ctx, o.ToNamespace, metav1.GetOptions{})
	require.NoError(t, err, "should have found to namespace %s", o.ToNamespace)
	require.NotNil(t, tons, "should have to namespace %s", o.ToNamespace)
	assert.Equal(t, o.ToNamespace, tons.Name, "toNamespace.Name")

	secretList, err := o.KubeClient.CoreV1().Secrets(o.ToNamespace).List(ctx, metav1.ListOptions{})
	require.NoError(t, err, "failed to list secrets in namespace %s", o.ToNamespace)
	require.NotNil(t, secretList, "should have SecretList from namespace %s", o.ToNamespace)
	assert.Len(t, secretList.Items, 1, "should have 1 secret in the namespace %s", o.ToNamespace)
}

func TestCopyNoSelectorOrNameFails(t *testing.T) {
	_, o := copy.NewCmdCopy()
	o.Namespace = ns
	o.KubeClient = fake.NewSimpleClientset()

	o.ToNamespace = "my-preview-env"
	o.CreateNamespace = true
	err := o.Run()
	require.Error(t, err, "should have failed")
}
