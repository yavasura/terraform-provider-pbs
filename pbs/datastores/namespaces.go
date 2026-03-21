/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

package datastores

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var namespaceComponentRegex = regexp.MustCompile(`^[A-Za-z0-9_][A-Za-z0-9._-]*$`)

// Namespace represents a PBS datastore namespace.
type Namespace struct {
	Store   string `json:"-"`
	Path    string `json:"ns"`
	Comment string `json:"comment,omitempty"`
}

// NamespaceDepth returns the number of path segments in a namespace.
func NamespaceDepth(namespace string) int {
	if namespace == "" {
		return 0
	}

	return len(strings.Split(namespace, "/"))
}

// IsValidNamespacePath validates a slash-separated PBS namespace path.
func IsValidNamespacePath(namespace string) bool {
	if namespace == "" {
		return false
	}

	parts := strings.Split(namespace, "/")
	if len(parts) == 0 || len(parts) > 7 {
		return false
	}

	for _, part := range parts {
		if !namespaceComponentRegex.MatchString(part) {
			return false
		}
	}

	return true
}

// ListNamespaces lists namespaces in a datastore.
func (c *Client) ListNamespaces(ctx context.Context, store string, maxDepth *int) ([]Namespace, error) {
	if store == "" {
		return nil, fmt.Errorf("store is required")
	}

	path := fmt.Sprintf("/admin/datastore/%s/namespace", url.PathEscape(store))
	if maxDepth != nil {
		path = fmt.Sprintf("%s?max-depth=%d", path, *maxDepth)
	}

	resp, err := c.api.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces for datastore %s: %w", store, err)
	}

	var namespaces []Namespace
	if err := json.Unmarshal(resp.Data, &namespaces); err != nil {
		return nil, fmt.Errorf("failed to unmarshal namespaces for datastore %s: %w", store, err)
	}

	for i := range namespaces {
		namespaces[i].Store = store
	}

	return namespaces, nil
}

// GetNamespace gets a specific namespace by scanning the datastore namespace list.
func (c *Client) GetNamespace(ctx context.Context, store, namespace string) (*Namespace, error) {
	if store == "" {
		return nil, fmt.Errorf("store is required")
	}
	if namespace == "" {
		return nil, fmt.Errorf("namespace is required")
	}

	namespaces, err := c.ListNamespaces(ctx, store, nil)
	if err != nil {
		return nil, err
	}

	for i := range namespaces {
		if namespaces[i].Path == namespace {
			return &namespaces[i], nil
		}
	}

	return nil, fmt.Errorf("namespace %s not found in datastore %s", namespace, store)
}

// CreateNamespace creates a namespace in a datastore.
func (c *Client) CreateNamespace(ctx context.Context, store, namespace, comment string) error {
	if store == "" {
		return fmt.Errorf("store is required")
	}
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	parent, name := SplitNamespacePath(namespace)
	body := map[string]interface{}{
		"name": name,
	}
	if parent != "" {
		body["parent"] = parent
	}
	if comment != "" {
		body["comment"] = comment
	}

	path := fmt.Sprintf("/admin/datastore/%s/namespace", url.PathEscape(store))
	_, err := c.api.Post(ctx, path, body)
	if err != nil {
		return fmt.Errorf("failed to create namespace %s in datastore %s: %w", namespace, store, err)
	}

	return nil
}

// SplitNamespacePath splits a namespace path into parent path and leaf name.
func SplitNamespacePath(namespace string) (string, string) {
	idx := strings.LastIndex(namespace, "/")
	if idx == -1 {
		return "", namespace
	}

	return namespace[:idx], namespace[idx+1:]
}

// DeleteNamespace deletes a namespace from a datastore.
func (c *Client) DeleteNamespace(ctx context.Context, store, namespace string) error {
	if store == "" {
		return fmt.Errorf("store is required")
	}
	if namespace == "" {
		return fmt.Errorf("namespace is required")
	}

	path := fmt.Sprintf(
		"/admin/datastore/%s/namespace?ns=%s",
		url.PathEscape(store),
		url.QueryEscape(namespace),
	)
	_, err := c.api.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to delete namespace %s in datastore %s: %w", namespace, store, err)
	}

	return nil
}
