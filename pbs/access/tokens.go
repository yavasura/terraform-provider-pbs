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
	"net/url"
	"strings"
)

// APIToken represents PBS API token metadata.
type APIToken struct {
	UserID    string   `json:"userid,omitempty"`
	TokenName string   `json:"-"`
	TokenID   string   `json:"tokenid,omitempty"`
	Comment   string   `json:"comment,omitempty"`
	Enable    *bool    `json:"enable,omitempty"`
	Expire    *int64   `json:"expire,omitempty"`
	Digest    string   `json:"digest,omitempty"`
	Delete    []string `json:"delete,omitempty"`
}

// GeneratedAPIToken represents the one-time secret returned by PBS when a token is created.
type GeneratedAPIToken struct {
	TokenID string `json:"tokenid"`
	Value   string `json:"value"`
}

// ListUserTokens lists API tokens for a specific PBS user.
func (c *Client) ListUserTokens(ctx context.Context, userID string) ([]APIToken, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}

	path := fmt.Sprintf("/access/users/%s/token", url.PathEscape(userID))
	resp, err := c.api.Get(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("failed to list tokens for user %s: %w", userID, err)
	}

	var tokens []APIToken
	if err := json.Unmarshal(resp.Data, &tokens); err != nil {
		return nil, fmt.Errorf("failed to unmarshal tokens for user %s: %w", userID, err)
	}

	for i := range tokens {
		normalizeTokenIdentity(&tokens[i], userID)
	}

	return tokens, nil
}

// GetUserToken gets a specific API token for a PBS user.
func (c *Client) GetUserToken(ctx context.Context, userID, tokenName string) (*APIToken, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}
	if tokenName == "" {
		return nil, fmt.Errorf("token name is required")
	}

	tokens, err := c.ListUserTokens(ctx, userID)
	if err != nil {
		return nil, err
	}

	for i := range tokens {
		if tokens[i].TokenName == tokenName || tokens[i].TokenID == FormatAPITokenID(userID, tokenName) {
			return &tokens[i], nil
		}
	}

	return nil, fmt.Errorf("token %s for user %s not found", tokenName, userID)
}

// CreateUserToken generates a new API token for a PBS user.
func (c *Client) CreateUserToken(ctx context.Context, userID, tokenName string, token *APIToken) (*GeneratedAPIToken, error) {
	if userID == "" {
		return nil, fmt.Errorf("user ID is required")
	}
	if tokenName == "" {
		return nil, fmt.Errorf("token name is required")
	}

	body := map[string]interface{}{}
	if token != nil {
		setTokenFields(body, token)
	}

	path := fmt.Sprintf(
		"/access/users/%s/token/%s",
		url.PathEscape(userID),
		url.PathEscape(tokenName),
	)

	resp, err := c.api.Post(ctx, path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create token %s for user %s: %w", tokenName, userID, err)
	}

	var generated GeneratedAPIToken
	if err := json.Unmarshal(resp.Data, &generated); err != nil {
		return nil, fmt.Errorf("failed to unmarshal generated token %s for user %s: %w", tokenName, userID, err)
	}

	if generated.TokenID == "" {
		generated.TokenID = FormatAPITokenID(userID, tokenName)
	}

	return &generated, nil
}

// DeleteUserToken deletes an API token for a PBS user.
func (c *Client) DeleteUserToken(ctx context.Context, userID, tokenName, digest string) error {
	if userID == "" {
		return fmt.Errorf("user ID is required")
	}
	if tokenName == "" {
		return fmt.Errorf("token name is required")
	}

	path := fmt.Sprintf(
		"/access/users/%s/token/%s",
		url.PathEscape(userID),
		url.PathEscape(tokenName),
	)
	if digest != "" {
		path = fmt.Sprintf("%s?digest=%s", path, url.QueryEscape(digest))
	}

	_, err := c.api.Delete(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to delete token %s for user %s: %w", tokenName, userID, err)
	}

	return nil
}

func setTokenFields(body map[string]interface{}, token *APIToken) {
	if token.Comment != "" {
		body["comment"] = token.Comment
	}
	if token.Enable != nil {
		body["enable"] = *token.Enable
	}
	if token.Expire != nil {
		body["expire"] = *token.Expire
	}
	if token.Digest != "" {
		body["digest"] = token.Digest
	}
	if len(token.Delete) > 0 {
		body["delete"] = token.Delete
	}
}

func normalizeTokenIdentity(token *APIToken, fallbackUserID string) {
	if token == nil {
		return
	}

	if token.TokenID == "" && token.UserID != "" && token.TokenName != "" {
		token.TokenID = FormatAPITokenID(token.UserID, token.TokenName)
	}

	if token.TokenID != "" {
		if userID, tokenName, ok := SplitAPITokenID(token.TokenID); ok {
			if token.UserID == "" {
				token.UserID = userID
			}
			token.TokenName = tokenName
		}
	}

	if token.UserID == "" {
		token.UserID = fallbackUserID
	}

	if token.TokenID == "" && token.UserID != "" && token.TokenName != "" {
		token.TokenID = FormatAPITokenID(token.UserID, token.TokenName)
	}
}

// FormatAPITokenID formats a full PBS token ID from a user ID and token name.
func FormatAPITokenID(userID, tokenName string) string {
	return fmt.Sprintf("%s!%s", userID, tokenName)
}

// SplitAPITokenID splits a PBS token ID into user ID and token name components.
func SplitAPITokenID(tokenID string) (string, string, bool) {
	parts := strings.SplitN(tokenID, "!", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", false
	}

	return parts[0], parts[1], true
}
