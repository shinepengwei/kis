package core

import "kis/apiserver/runtime"

// Pod is a collection of containers, used as either input (create, update) or as output (list, get).
type Pod struct {
	Name string
}


func (in *Pod) GetName() string{
	return in.Name
}
// PodList is a list of Pods.
type PodList struct {
	Items []Pod
}



func (in *PodList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(PodList)
	in.DeepCopyInto(out)
	return out
}
func (in *PodList) DeepCopyInto(out *PodList) {
	*out = *in
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Pod, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

func (in *Pod) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(Pod)
	in.DeepCopyInto(out)
	return out
}

func (in *Pod) DeepCopyInto(out *Pod) {
	*out = *in
	out.Name = in.Name
	return
}