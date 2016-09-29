// Copyright 2016 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may not
// use this file except in compliance with the License. A copy of the
// License is located at
//
// http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
// either express or implied. See the License for the specific language governing
// permissions and limitations under the License.

// Package gatherers contains routines for different types of inventory gatherers

package gatherers

import (
	"github.com/aws/amazon-ssm-agent/agent/context"
	"github.com/aws/amazon-ssm-agent/agent/contracts"
	"github.com/aws/amazon-ssm-agent/agent/inventory/gatherers/application"
	"github.com/aws/amazon-ssm-agent/agent/inventory/gatherers/custom"
	"github.com/aws/amazon-ssm-agent/agent/inventory/gatherers/windowsUpdate"
	"github.com/aws/amazon-ssm-agent/agent/inventory/model"
)

// T defines operations that all inventory gatherers support
type T interface {
	//returns the Name of the gatherer
	Name() string
	//runs the gatherer with a given configuration
	//returns array of inventory.Item as custom gatherer collects multiple
	//inventory items at a time
	Run(context context.T, configuration inventory.Config) ([]inventory.Item, error)
	//stops the execution of a gatherer
	RequestStop(stopType contracts.StopType) error
}

// Registry stores all supported types of inventory gatherers
type Registry map[string]T

// LoadGatherers loads supported inventory gatherers in memory
func LoadGatherers(context context.T) Registry {
	var m Registry
	var names []string
	m = make(map[string]T)
	log := context.Log()

	log.Infof("Loading available inventory gatherers")

	// Load application inventory item gather
	if a, err := application.Gatherer(context); err != nil {
		log.Errorf("Fake application gatherer isn't properly configured - %v", err.Error())
	} else {
		m[a.Name()] = a
		names = append(names, a.Name())
	}

	// Load windowsUpdate inventory item gather
	if a, err := windowsUpdate.Gatherer(context); err != nil {
		log.Errorf("Windows update gatherer isn't properly configured - %v", err.Error())
	} else {
		m[a.Name()] = a
		names = append(names, a.Name())
	}

	// Load custom inventory items gather
	if cg, err := custom.Gatherer(context); err != nil {
		log.Errorf("Custom inventory gatherer isn't properly configured - %v", err.Error())
	} else {
		m[cg.Name()] = cg
		names = append(names, cg.Name())
	}

	log.Infof("Supported inventory gatherers : %v", names)

	return m
}
