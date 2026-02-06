package preview

import (
	"context"
	"time"

	"github.com/argoproj/argo-cd/v3/applicationset/generators"
	"github.com/argoproj/argo-cd/v3/applicationset/services"
	argoprojiov1alpha1 "github.com/argoproj/argo-cd/v3/pkg/apis/application/v1alpha1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ generators.Generator = (*OfflineGitGenerator)(nil)

// OfflineGitGenerator wraps the ArgoCD Git generator with a fake client
// that returns empty AppProjects to skip GPG verification in offline mode.
type OfflineGitGenerator struct {
	inner generators.Generator
}

// NewOfflineGitGenerator creates a Git generator suitable for offline use.
func NewOfflineGitGenerator(repos services.Repos) generators.Generator {
	inner := generators.NewGitGenerator(repos, "")
	return &OfflineGitGenerator{inner: inner}
}

func (g *OfflineGitGenerator) GetTemplate(appSetGenerator *argoprojiov1alpha1.ApplicationSetGenerator) *argoprojiov1alpha1.ApplicationSetTemplate {
	return g.inner.GetTemplate(appSetGenerator)
}

func (g *OfflineGitGenerator) GetRequeueAfter(appSetGenerator *argoprojiov1alpha1.ApplicationSetGenerator) time.Duration {
	return g.inner.GetRequeueAfter(appSetGenerator)
}

func (g *OfflineGitGenerator) GenerateParams(appSetGenerator *argoprojiov1alpha1.ApplicationSetGenerator, appSet *argoprojiov1alpha1.ApplicationSet, _ client.Client) ([]map[string]any, error) {
	// Use our fake client that returns empty AppProjects (no GPG keys = no verification)
	fakeClient := &fakeK8sClient{}
	return g.inner.GenerateParams(appSetGenerator, appSet, fakeClient)
}

// fakeK8sClient implements controller-runtime client.Client interface
// but only supports Get for AppProject, returning an empty project.
type fakeK8sClient struct{}

func (f *fakeK8sClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	// Return an empty AppProject with no signature keys - this disables GPG verification
	if appProject, ok := obj.(*argoprojiov1alpha1.AppProject); ok {
		appProject.Name = key.Name
		appProject.Namespace = key.Namespace
		appProject.Spec.SignatureKeys = nil // No GPG verification
	}
	return nil
}

// The following methods are required by the client.Client interface but not used
func (f *fakeK8sClient) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	return nil
}

func (f *fakeK8sClient) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	return nil
}

func (f *fakeK8sClient) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	return nil
}

func (f *fakeK8sClient) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	return nil
}

func (f *fakeK8sClient) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	return nil
}

func (f *fakeK8sClient) DeleteAllOf(ctx context.Context, obj client.Object, opts ...client.DeleteAllOfOption) error {
	return nil
}

func (f *fakeK8sClient) Status() client.SubResourceWriter {
	return nil
}

func (f *fakeK8sClient) SubResource(subResource string) client.SubResourceClient {
	return nil
}

func (f *fakeK8sClient) Scheme() *runtime.Scheme {
	return nil
}

func (f *fakeK8sClient) RESTMapper() meta.RESTMapper {
	return nil
}

func (f *fakeK8sClient) GroupVersionKindFor(obj runtime.Object) (schema.GroupVersionKind, error) {
	return schema.GroupVersionKind{}, nil
}

func (f *fakeK8sClient) IsObjectNamespaced(obj runtime.Object) (bool, error) {
	return true, nil
}
