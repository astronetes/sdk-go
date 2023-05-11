package controller

const (
	ReachMaxAllowedAttemptsEvent = "reachMaxAllowedAttempts"
	NewRequestEvent              = "newRequest"
	TerminatingEvent             = "terminatingResource"
	ErrorEvent                   = "error"
	FinalizerExists              = "finalizerExists"
	MissingFinalizer             = "missingFinalizer"
	ResumeCreationEvent          = "resumeCreation"
	ResumeDeletionEvent          = "resumeDeletion"
)
