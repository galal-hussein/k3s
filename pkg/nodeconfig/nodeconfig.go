package nodeconfig

import (
	"encoding/json"
	"github.com/pkg/errors"
	"os"
	"strings"
	corev1 "k8s.io/api/core/v1"
)

const (
	NodeArgsAnnotation = "k3s.io/node-args"
	NodeEnvAnnotation = "k3s.io/node-env"
)

func getNodeArgs(nodeName string) (string,error) {
	nodeArgs, err := json.Marshal(os.Args[1:])
	if err != nil {
		return "", errors.Wrapf(err, "Failed to retrieve argument list for node %s", nodeName)
	}
	return string(nodeArgs), nil
}

func getNodeEnv(nodeName string) (string, error) {
	k3sEnv := make(map[string]string)
	for _, v := range os.Environ() {
		keyValue := strings.SplitN(v, "=", 2)
		if strings.HasPrefix(keyValue[0], "K3S_") {
			k3sEnv[keyValue[0]] = keyValue[1]
		}
	}
	k3sEnvJson, err := json.Marshal(k3sEnv)
	if err != nil {
		return "", errors.Wrapf(err, "Failed to retrieve environment map for node %s", nodeName)
	}
	return string(k3sEnvJson), nil
}

func setNodeAnnotation(node *corev1.Node, annotationKey, annotationValue string) {
	if node.Annotations == nil {
		node.Annotations = make(map[string]string)
	}
	node.Annotations[annotationKey] = annotationValue
}

func SetNodeConfigAnnotations(node *corev1.Node) error {
	nodeArgs, err := getNodeArgs(node.Name)
	if err != nil {
		return err
	}
	nodeEnv, err := getNodeEnv(node.Name)
	if err != nil {
		return err
	}
	if node.Annotations == nil {
		node.Annotations = make(map[string]string)
	}
	node.Annotations[NodeEnvAnnotation] = nodeEnv
	node.Annotations[NodeArgsAnnotation] = nodeArgs
	return nil
}