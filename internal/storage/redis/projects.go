package redis

import (
	"context"
	"fmt"
)

func (c *Client) PushProject(ctx context.Context, projectID int, state string) error {
	key := fmt.Sprintf("project:%d", projectID)
	err := c.client.Set(ctx, key, state, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to push project to Redis: %v", err)
	}
	return nil
}

func (c *Client) DeleteProject(ctx context.Context, projectID int) error {
	key := fmt.Sprintf("project:%d", projectID)
	err := c.client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("failed to delete project from Redis: %v", err)
	}
	return nil
}
