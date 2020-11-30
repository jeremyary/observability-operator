package utils

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	v13 "github.com/openshift/api/config/v1"
	routev1 "github.com/openshift/api/route/v1"
	v1 "k8s.io/api/core/v1"
	k8sclient "sigs.k8s.io/controller-runtime/pkg/client"
)

// Returns the cluster id by querying the ClusterVersion resource
func GetClusterId(ctx context.Context, client k8sclient.Client) (string, error) {
	v := &v13.ClusterVersion{}
	selector := k8sclient.ObjectKey{
		Name: "version",
	}

	err := client.Get(ctx, selector, v)
	if err != nil {
		return "", err
	}

	return string(v.Spec.ClusterID), nil
}

func IsRouteReads(route *routev1.Route) bool {
	if route == nil {
		return false
	}
	// A route has a an array of Ingress where each have an array of conditions
	for _, ingress := range route.Status.Ingress {
		for _, condition := range ingress.Conditions {
			// A successful route will have the admitted condition type as true
			if condition.Type == routev1.RouteAdmitted && condition.Status != v1.ConditionTrue {
				return false
			}
		}
	}
	return true
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomString(s int) string {
	b := GenerateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b)
}
