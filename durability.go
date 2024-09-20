/*
Copyright 2021 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"

	"vitess.io/vitess/go/vt/log"
	topodatapb "vitess.io/vitess/go/vt/proto/topodata"
	"vitess.io/vitess/go/vt/vtctl/reparentutil/promotionrule"
)

//=======================================================================

// A NewDurabler is a function that creates a new Durabler based on the
// properties specified in the input map. Every Durabler must
// register a NewDurabler function.
type NewDurabler func() Durabler

var (
	// durabilityPolicies is a map that stores the functions needed to create a new Durabler
	durabilityPolicies = make(map[string]NewDurabler)
)

// Durabler is the interface which is used to get the promotion rules for candidates and the semi sync setup
type Durabler interface {
	// PromotionRule represents the precedence in which we want to tablets to be promoted.
	// The higher the promotion rule of a tablet, the more we want it to be promoted in case of a failover
	PromotionRule(*topodatapb.Tablet) promotionrule.CandidatePromotionRule
	// SemiSyncAckers represents the number of semi-sync ackers required for a given tablet if it were to become the PRIMARY instance
	SemiSyncAckers(*topodatapb.Tablet) int
	// IsReplicaSemiSync returns whether the "replica" should send semi-sync acks if "primary" were to become the PRIMARY instance
	IsReplicaSemiSync(primary, replica *topodatapb.Tablet) bool
}

func RegisterDurability(name string, newDurablerFunc NewDurabler) {
	if durabilityPolicies[name] != nil {
		log.Fatalf("durability policy %v already registered", name)
	}
	durabilityPolicies[name] = newDurablerFunc
}

//=======================================================================

// GetDurabilityPolicy is used to get a new durability policy from the registered policies
func GetDurabilityPolicy(name string) (Durabler, error) {
	newDurabilityCreationFunc, found := durabilityPolicies[name]
	if !found {
		return nil, fmt.Errorf("durability policy %v not found", name)
	}
	return newDurabilityCreationFunc(), nil
}

// CheckDurabilityPolicyExists is used to check if the durability policy is part of the registered policies
func CheckDurabilityPolicyExists(name string) bool {
	_, found := durabilityPolicies[name]
	return found
}

// PromotionRule returns the promotion rule for the instance.
func PromotionRule(durability Durabler, tablet *topodatapb.Tablet) promotionrule.CandidatePromotionRule {
	// Prevent panics.
	if tablet == nil || tablet.Alias == nil {
		return promotionrule.MustNot
	}
	return durability.PromotionRule(tablet)
}

// SemiSyncAckers returns the primary semi-sync setting for the instance.
// 0 means none. Non-zero specifies the number of required ackers.
func SemiSyncAckers(durability Durabler, tablet *topodatapb.Tablet) int {
	return durability.SemiSyncAckers(tablet)
}

// IsReplicaSemiSync returns the replica semi-sync setting from the tablet record.
// Prefer using this function if tablet record is available.
func IsReplicaSemiSync(durability Durabler, primary, replica *topodatapb.Tablet) bool {
	// Prevent panics.
	if primary == nil || primary.Alias == nil || replica == nil || replica.Alias == nil {
		return false
	}
	return durability.IsReplicaSemiSync(primary, replica)
}
