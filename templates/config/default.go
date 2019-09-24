package config

import "time"

const (
	// Root URLs
	// ProdKumparanDesktopWebRootURL :nodoc:
	ProdKumparanDesktopWebRootURL = "https://kumparan.com/"
	// ProdKumparanMobileWebRootURL :nodoc:
	ProdKumparanMobileWebRootURL = "https://m.kumparan.com/"
	// StagingKumparanDesktopWebRootURL :nodoc:
	StagingKumparanDesktopWebRootURL = "https://mpreview.kumparan.com/"
	// StagingKumparanMobileWebRootURL :nodoc:
	StagingKumparanMobileWebRootURL = "https://mpreview.kumparan.com/"
	// DevKumparanDesktopWebRootURL :nodoc:
	DevKumparanDesktopWebRootURL = "https://dev.kumparan.com/"
	// DevKumparanMobileWebRootURL :nodoc:
	DevKumparanMobileWebRootURL = "https://dev.kumparan.com/"

	// DefaultRPCClientTimeout :nodoc:
	DefaultRPCClientTimeout = 1100
	// DefaultRPCServerTimeout :nodoc:
	DefaultRPCServerTimeout = 1000
	// DefaultHTTPTimeout :nodoc:
	DefaultHTTPTimeout = 10000 * time.Millisecond

	// DefaultCacheTTL in milliseconds
	DefaultCacheTTL = 900000

	// Worker namespace
	GocraftWorkerNamespace = "mail-service"
	// GocraftWorkerPoolConcurrency Worker default concurrency
	GocraftWorkerPoolConcurrency = 10
)
