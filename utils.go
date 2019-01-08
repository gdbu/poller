package poller

import "os"

func areBothEmpty(last, info os.FileInfo) (empty bool) {
	if last != nil {
		return
	}

	if info != nil {
		return
	}

	return true
}

func wasCreated(last, info os.FileInfo) (created bool) {
	if last != nil {
		return
	}

	if info == nil {
		return
	}

	return true
}

func wasRemoved(last, info os.FileInfo) (removed bool) {
	if last == nil {
		return
	}

	if info != nil {
		return
	}

	return true
}

func wasUpdated(last, info os.FileInfo) (updated bool) {
	return !last.ModTime().Equal(info.ModTime())
}

func hadPermissionsChange(last, info os.FileInfo) (changed bool) {
	return last.Mode() != info.Mode()
}

// OnEvent is the on event callback handler
type OnEvent func(Event)
