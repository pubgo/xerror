package logutil

import (
	"k8s.io/klog/v2"
	_ "k8s.io/klog/v2"
)

func init() {
	klog.Background()
}
