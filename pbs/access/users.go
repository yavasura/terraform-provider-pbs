/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 */

// Package access provides API client functionality for PBS access control objects.
package access

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/micah/terraform-provider-pbs/pbs/api"
)

// Client represents the access API client.
type Client struct {
	api *api.Client
}

// NewClient creates a new access API client.
func NewClient(apiClient *api.Client) *Client {
	return &Client{api: apiClient}
}

// User represents a PBS user account.
type User struct {
	UserID    string   `json:"userid"`
	Comment   string   `json:"comment,omitempty"`
	Enable    *bool    `json:"enable,omitempty"`
	Expire    *int64   `json:"expire,omitempty"`
	FirstName string   `json:"firstname,omitempty"`
	LastName  string   `json:"lastname,omitempty"`
	Email     string   `json:"email,omitempty"`
	Digest    string   `json:"digest,omitempty"`
	Delete    []string `json:"delete,omitempty"`
}

// ListUsers lists all PBS users.
func (c *Client) ListUsers(ctx context.Context) ([]User, error) {
	resp, err := c.api.Get(ctx, "/access/users")
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	var users []User
	if err := json.Unmarshal(resp.Data, &users); err != nil {
		return nil, fmt.Errorf("failed to unmarshal users: %w", err)
	}

	return users, nil
}

// GetUser gets a specific PBS user by user ID.
func (c *Client) GetUser(ctx context.Context, userID string) (*User, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	path := fmt.Sprintf("/access/users/%s", url.PathEscape(userID))
	resp, err := c.api.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to get user %s: %w", userID, err)
	}

	var user User
	if err := json.Unmarshal(resp.Data, &user); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user %s: %w", userID, err)
	}

	return &user, nil
}

// CreateUser creates a new PBS user account.
func (c *Client) CreateUser(ctx context.Context, user *User) error {
	if user == nil {
		return fmt.Errorf("user is required")
	}
	if user.UserID == "" {
		return fmt.Errorf("user ID is required")
	}

	body := map[string]interface{}{
		"userid": user.UserID,
	}
	setUserFields(body, user)

	_, err := c.api.Post(ctx, "/access/users", body)
	if err != nil {
		return fmt.Errorf("failed to create user %s: %w", user.UserID, err)
	}

	return nil
}

// UpdateUser updates an existing PBS user account.
func (c *Client) UpdateUser(ctx context.Context, userID string, user *User) error {
	if user == nil {
		return fmt.Errorf("user is required")
	}
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}

	body := map[string]interface{}{}
	setUserFields(body, user)

	if user.Digest != "" {
		body["digest"] = user.Digest
	}
	if len(user.Delete) > 0 {
		body["delete"] = user.Delete
	}

	path := fmt.Sprintf("/access/users/%s", url.PathEscape(userID))
	_, err := c.api.Put(ctx, path, body)
	if err != nil {
		return fmt.Errorf("failed to update user %s: %w", userID, err)
	}

	return nil
}

// DeleteUser deletes a PBS user account.
func (c *Client) DeleteUser(ctx context.Context, userID, digest string) error {
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}

	path := fmt.Sprintf("/access/users/%s", url.PathEscape(userID))
	if digest != "" {
		path = fmt.Sprintf("%s?digest=%s", path, url.QueryEscape(digest))
	}

	_, err := c.api.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to delete user %s: %w", userID, err)
	}

	return nil
}

func setUserFields(body map[string]interface{}, user *User) {
	if user.Comment != "" {
		body["comment"] = user.Comment
	}
	if user.Enable != nil {
		body["enable"] = *user.Enable
	}
	if user.Expire != nil {
		body["expire"] = *user.Expire
	}
	if user.FirstName != "" {
		body["firstname"] = user.FirstName
	}
	if user.LastName != "" {
		body["lastname"] = user.LastName
	}
	if user.Email != "" {
		body["email"] = user.Email
	}
}
