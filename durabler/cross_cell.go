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
	topodatapb "vitess.io/vitess/go/vt/proto/topodata"
	"vitess.io/vitess/go/vt/vtctl/reparentutil/promotionrule"
)

// durabilityCrossCell has 1 semi-sync setup. It only allows Primary and Replica type servers from a different cell to acknowledge semi sync.
// This means that a transaction must be in two cells for it to be acknowledged
// It returns NeutralPromoteRule for Primary and Replica tablet types, MustNotPromoteRule for everything else
type durabilityCrossCell struct {
	rdonlySemiSync bool
}

// PromotionRule implements the Durabler interface
func (d *durabilityCrossCell) PromotionRule(tablet *topodatapb.Tablet) promotionrule.CandidatePromotionRule {
	switch tablet.Type {
	case topodatapb.TabletType_PRIMARY, topodatapb.TabletType_REPLICA:
		return promotionrule.Neutral
	}
	return promotionrule.MustNot
}

// SemiSyncAckers implements the Durabler interface
func (d *durabilityCrossCell) SemiSyncAckers(tablet *topodatapb.Tablet) int {
	return 1
}

// IsReplicaSemiSync implements the Durabler interface
func (d *durabilityCrossCell) IsReplicaSemiSync(primary, replica *topodatapb.Tablet) bool {
	switch replica.Type {
	case topodatapb.TabletType_PRIMARY, topodatapb.TabletType_REPLICA:
		return primary.Alias.Cell != replica.Alias.Cell
	case topodatapb.TabletType_RDONLY:
		return d.rdonlySemiSync && primary.Alias.Cell != replica.Alias.Cell
	}
	return false
}

var DurabilityCrossCell durabilityCrossCell
