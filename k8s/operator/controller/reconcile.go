package controller

import "time"

type ReconcileOpts struct {
	Timeout            time.Duration
	DefaultRequeueTime time.Duration
}
