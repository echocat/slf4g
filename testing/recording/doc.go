// Package recording provides the possibility to record instances of events.
//
// This can either be done by using a direct instance of a log.Logger using
// NewLogger() to create it.
//
// ... or hook it fully into the full application by calling NewProvider() and
// using then afterwards Provider.HookGlobally() to make it available for every
// piece of code that tries to log something. See the example for more details.
package recording
