// Copyright (c) Keeper Security, Inc.
// SPDX-License-Identifier: MPL-2.0

package enterprisepush

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

// readFileAndParseJSON reads the file at filePath and parses it as JSON.
// Returns content, fileData for the API, and an error if read or parse fails.
func readFileAndParseJSON(filePath string) (content []byte, fileData map[string]interface{}, err error) {
	content, err = os.ReadFile(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read file at %q: %w", filePath, err)
	}
	fileData = make(map[string]interface{})
	if err := json.Unmarshal(content, &fileData); err != nil {
		return nil, nil, fmt.Errorf("file %q is not valid JSON: %w", filePath, err)
	}
	return content, fileData, nil
}

// computeID returns a deterministic ID from content and config so that:
// - Same content + emails + teams + managed_company → same ID (no re-push).
// - Any change → new ID → replace → push again.
func computeID(content []byte, data *EnterprisePushResourceModel) string {
	sortedEmails := sortedSetStrings(data.Email)
	sortedTeams := sortedSetStrings(data.Team)
	managedCompany := ""
	if !data.ManagedCompany.IsNull() {
		managedCompany = data.ManagedCompany.ValueString()
	}
	h := sha256.New()
	h.Write(content)
	h.Write([]byte("\n"))
	h.Write([]byte(strings.Join(sortedEmails, "\n")))
	h.Write([]byte("\n"))
	h.Write([]byte(strings.Join(sortedTeams, "\n")))
	h.Write([]byte("\n"))
	h.Write([]byte(managedCompany))
	return hex.EncodeToString(h.Sum(nil))
}

// contentSHA256Hex returns the SHA256 hash of content as a hex string (used for file_content_sha256 and by ModifyPlan).
func contentSHA256Hex(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}

// sortedSetStrings returns a sorted slice of non-null, known string values from the set.
func sortedSetStrings(set types.Set) []string {
	var out []string
	for _, elem := range set.Elements() {
		s, ok := elem.(types.String)
		if !ok || s.IsNull() || s.IsUnknown() {
			continue
		}
		out = append(out, s.ValueString())
	}
	sort.Strings(out)
	return out
}

// buildEnterprisePushCommand builds the enterprise-push command with FILEDATA and --email/--team flags from the model.
func buildEnterprisePushCommand(data *EnterprisePushResourceModel) string {
	return buildEnterprisePushCommandWithTargets(sortedSetStrings(data.Email), sortedSetStrings(data.Team))
}

// buildEnterprisePushCommandWithTargets builds the command for given email and team slices (e.g. for Update with added-only targets).
func buildEnterprisePushCommandWithTargets(emails, teams []string) string {
	var parts []string
	parts = append(parts, "enterprise-push", "FILEDATA")
	for _, e := range emails {
		parts = append(parts, fmt.Sprintf("--email='%s'", e))
	}
	for _, t := range teams {
		parts = append(parts, fmt.Sprintf("--team='%s'", t))
	}
	return strings.Join(parts, " ")
}

// setDifference returns elements that are in a but not in b.
func setDifference(a, b []string) []string {
	bm := make(map[string]struct{}, len(b))
	for _, x := range b {
		bm[x] = struct{}{}
	}
	var out []string
	for _, x := range a {
		if _, ok := bm[x]; !ok {
			out = append(out, x)
		}
	}
	return out
}
