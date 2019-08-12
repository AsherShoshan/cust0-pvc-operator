package controller

import (
	"github.com/AsherShoshan/cust0-pvc-operator/pkg/controller/pvc"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, pvc.Add)
}
