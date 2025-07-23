package ui

import "context"

// DashboardUI abstracts a terminal dashboard implementation.
type DashboardUI interface {
	Run(ctx context.Context) error
}
