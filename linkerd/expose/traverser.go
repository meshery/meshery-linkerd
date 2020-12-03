package expose

import (
	"context"
	"errors"
	"fmt"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
)

// Traverser can be used to traverse resources in the cluster
type Traverser struct {
	// Resources should be in format
	// ["<namespace> <type> <name>", "<namespace2> <type2> <name2>"...]
	Resources []string

	Client *kubernetes.Clientset

	Logger
}

// Visit function traverses each of the resource mentioned in the Traverser struct
func (traverser *Traverser) Visit(f func(runtime.Object, string, error) error, continueOnError bool) error {
	var errs []error
	for _, res := range traverser.Resources {
		md := strings.Split(res, " ")
		if len(md) != 3 {
			return fmt.Errorf(`invalid resource definition, valid format is: "<namespace> <type> <name>"`)
		}

		ns := md[0]   // Namespace in which the resoruce should be searched
		typ := md[1]  // Type of the resource
		name := md[2] // Name of the resource

		switch typ {
		case "service":
			svc, err := traverser.Client.CoreV1().Services(ns).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				traverser.Logger.Error(err)
				errs = append(errs, err)
				if !continueOnError {
					return ErrGettingResource(err)
				}
			}
			// Placing Kind and APIVersion manually
			// because client-go omits them
			// Please do no remove
			svc.Kind = "Service"
			svc.APIVersion = "v1"
			if err := f(svc, svc.Name, err); err != nil {
				traverser.Logger.Error(err)
				errs = append(errs, err)
				if !continueOnError {
					return ErrGettingResource(err)
				}
			}
		case "deployment":
			dep, err := traverser.Client.AppsV1().Deployments(ns).Get(context.TODO(), name, metav1.GetOptions{})
			if err != nil {
				traverser.Logger.Error(err)
				errs = append(errs, err)
				if !continueOnError {
					return ErrGettingResource(err)
				}
			}
			// Placing Kind and APIVersion manually
			// because client-go omits them
			// Please do no remove
			dep.Kind = "Deployment"
			dep.APIVersion = "apps/v1"
			if err := f(dep, dep.Name, err); err != nil {
				traverser.Logger.Error(err)
				errs = append(errs, err)
				if !continueOnError {
					return ErrGettingResource(err)
				}
			}
		default:
			// Don't do anything
			traverser.Logger.Warn(fmt.Errorf("invalid resource type"))
		}
	}

	err := combineErrors(errs, "\n")
	if err != nil {
		ErrTraverser(err)
	}

	return nil
}

// combineErrors merges a slice of error
// into one error seperated by the given seperator
func combineErrors(errs []error, sep string) error {
	if len(errs) == 0 {
		return nil
	}

	var errString []string
	for _, err := range errs {
		errString = append(errString, err.Error())
	}

	return errors.New(strings.Join(errString, sep))
}
