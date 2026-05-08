package server

// Site-controller consumes backend events in this process. The current
// implementation keeps the event surface intentionally small: node online and
// health events trigger future switch-policy refresh handling, while the gRPC
// API remains the authoritative control path for site operations.
