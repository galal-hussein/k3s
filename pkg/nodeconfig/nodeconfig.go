package nodeconfig

import (
	"encoding/json"
	"expvar"
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
	nodeArgsList := []string{}
	for i, arg := range os.Args[1:] {
		if isSecret(arg) {
			nodeArgsList = append(nodeArgsList, "")
		} else {
			nodeArgsList = append(nodeArgsList, arg)
		}
	}
	nodeArgs, err := json.Marshal(nodeArgsList)
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
	for key, value := range keyValue {
		if isSecret(key) {
			k3sEnv[key] = ""
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

func isSecret(key string) bool {
	secretData := []string{
		"K3S_TOKEN",
		"K3S_DATASTORE_",
		"K3S_AGENT_TOKEN",
		"K3S_CLUSTER_SECRET",
		"token",
		"agent-token",
		"datastore-",
		"cluster-secret",
	}
	for _, secret := range secretData {
		if strings.contains(key, secret) {
			return true
		}
	}
	return false
}