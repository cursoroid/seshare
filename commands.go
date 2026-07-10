package main

import "errors"

// Command entrypoints. Stubs for now — implemented per the plan.
// Each will split into its own file (contacts, session, transport, importer)
// as it grows; kept together while empty to avoid scaffolding for later.

var errNotImplemented = errors.New("not implemented yet")

func cmdPair(args []string) error { return errNotImplemented }

func cmdSend(args []string) error { return errNotImplemented }

func cmdRecv(args []string) error { return errNotImplemented }
