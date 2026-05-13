package server

// Event handling will consume node/state/health and registry events in the
// same process. For this first production slice, site actions are driven via
// API-GW -> site-controller gRPC, while the service remains event-ready.
