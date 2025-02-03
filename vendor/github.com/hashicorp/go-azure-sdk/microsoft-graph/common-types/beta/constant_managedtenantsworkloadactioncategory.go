package beta

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Copyright (c) HashiCorp Inc. All rights reserved.
// Licensed under the MIT License. See NOTICE.txt in the project root for license information.

type ManagedTenantsWorkloadActionCategory string

const (
	ManagedTenantsWorkloadActionCategory_Automated ManagedTenantsWorkloadActionCategory = "automated"
	ManagedTenantsWorkloadActionCategory_Manual    ManagedTenantsWorkloadActionCategory = "manual"
)

func PossibleValuesForManagedTenantsWorkloadActionCategory() []string {
	return []string{
		string(ManagedTenantsWorkloadActionCategory_Automated),
		string(ManagedTenantsWorkloadActionCategory_Manual),
	}
}

func (s *ManagedTenantsWorkloadActionCategory) UnmarshalJSON(bytes []byte) error {
	var decoded string
	if err := json.Unmarshal(bytes, &decoded); err != nil {
		return fmt.Errorf("unmarshaling: %+v", err)
	}
	out, err := parseManagedTenantsWorkloadActionCategory(decoded)
	if err != nil {
		return fmt.Errorf("parsing %q: %+v", decoded, err)
	}
	*s = *out
	return nil
}

func parseManagedTenantsWorkloadActionCategory(input string) (*ManagedTenantsWorkloadActionCategory, error) {
	vals := map[string]ManagedTenantsWorkloadActionCategory{
		"automated": ManagedTenantsWorkloadActionCategory_Automated,
		"manual":    ManagedTenantsWorkloadActionCategory_Manual,
	}
	if v, ok := vals[strings.ToLower(input)]; ok {
		return &v, nil
	}

	// otherwise presume it's an undefined value and best-effort it
	out := ManagedTenantsWorkloadActionCategory(input)
	return &out, nil
}
