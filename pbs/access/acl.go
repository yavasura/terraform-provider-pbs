/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package access

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// ACL represents a PBS access control entry.
type ACL struct {
	Path      string `json:"path"`
	UGID      string `json:"ugid"`
	RoleID    string `json:"roleid"`
	Propagate *bool  `json:"propagate,omitempty"`
}

// ListACLs lists all PBS ACL entries.
func (c *Client) ListACLs(ctx context.Context) ([]ACL, error) {
	resp, err := c.api.Get(ctx, "/access/acl")
	if err != nil {
		return nil, fmt.Errorf("failed to list ACLs: %w", err)
	}

	var acls []ACL
	if err := json.Unmarshal(resp.Data, &acls); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ACLs: %w", err)
	}

	return acls, nil
}

// GetACL gets a specific ACL entry by path and subject.
func (c *Client) GetACL(ctx context.Context, aclPath, ugid string) (*ACL, error) {
	if aclPath == "" {
		return nil, fmt.Errorf("path is required")
	}
	if ugid == "" {
		return nil, fmt.Errorf("ugid is required")
	}

	acls, err := c.ListACLs(ctx)
	if err != nil {
		return nil, err
	}

	for i := range acls {
		if acls[i].Path == aclPath && acls[i].UGID == ugid {
			return &acls[i], nil
		}
	}

	return nil, fmt.Errorf("ACL for %s on %s not found", ugid, aclPath)
}

// GetACLForRole gets a specific ACL entry by path, subject, and role.
func (c *Client) GetACLForRole(ctx context.Context, aclPath, ugid, roleID string) (*ACL, error) {
	if roleID == "" {
		return c.GetACL(ctx, aclPath, ugid)
	}

	acls, err := c.ListACLs(ctx)
	if err != nil {
		return nil, err
	}

	for i := range acls {
		if acls[i].Path == aclPath && acls[i].UGID == ugid && acls[i].RoleID == roleID {
			return &acls[i], nil
		}
	}

	return nil, fmt.Errorf("ACL for %s on %s with role %s not found", ugid, aclPath, roleID)
}

// SetACL creates or updates a PBS ACL entry.
func (c *Client) SetACL(ctx context.Context, acl *ACL) error {
	if acl == nil {
		return fmt.Errorf("acl is required")
	}
	if acl.Path == "" {
		return fmt.Errorf("path is required")
	}
	if acl.UGID == "" {
		return fmt.Errorf("ugid is required")
	}
	if acl.RoleID == "" {
		return fmt.Errorf("roleid is required")
	}

	body := map[string]interface{}{
		"path": acl.Path,
		"role": acl.RoleID,
	}
	if strings.Contains(acl.UGID, "@") {
		body["auth-id"] = acl.UGID
	} else {
		body["group"] = acl.UGID
	}
	if acl.Propagate != nil {
		body["propagate"] = *acl.Propagate
	}

	_, err := c.api.Put(ctx, "/access/acl", body)
	if err != nil {
		return fmt.Errorf("failed to set ACL for %s on %s: %w", acl.UGID, acl.Path, err)
	}

	return nil
}

// DeleteACL deletes a PBS ACL entry.
func (c *Client) DeleteACL(ctx context.Context, aclPath, ugid, roleID string) error {
	if aclPath == "" {
		return fmt.Errorf("path is required")
	}
	if ugid == "" {
		return fmt.Errorf("ugid is required")
	}

	body := map[string]interface{}{
		"path":   aclPath,
		"delete": true,
	}
	if strings.Contains(ugid, "@") {
		body["auth-id"] = ugid
	} else {
		body["group"] = ugid
	}
	if roleID != "" {
		body["role"] = roleID
	}

	_, err := c.api.Put(ctx, "/access/acl", body)
	if err != nil {
		return fmt.Errorf("failed to delete ACL for %s on %s: %w", ugid, aclPath, err)
	}

	return nil
}
