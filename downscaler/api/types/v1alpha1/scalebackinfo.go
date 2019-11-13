package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ScaleBackInfoSpec struct {
	Domain            string `json:"domain"`
	OriginalService   string `json:"originalService"`
	ParentIngressName string `json:"parentIngressName"`
}

type ScaleBackInfo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec ScaleBackInfoSpec `json:"spec"`
}

type ScaleBackInfoList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []ScaleBackInfo `json:"items"`
}

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *ScaleBackInfo) DeepCopyInto(out *ScaleBackInfo) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	out.Spec = ScaleBackInfoSpec{
		Replicas: in.Spec.Replicas,
	}
}

// DeepCopyObject returns a generically typed copy of an object
func (in *ScaleBackInfo) DeepCopyObject() runtime.Object {
	out := ScaleBackInfo{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *ScaleBackInfoList) DeepCopyObject() runtime.Object {
	out := ScaleBackInfoList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]ScaleBackInfo, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}
