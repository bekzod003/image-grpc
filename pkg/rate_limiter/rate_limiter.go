package ratelimiter

import "sync"

var (
	listLimit             sync.Map
	downloadAndStoreLimit sync.Map
)

const (
	maxDownloadAndStoreLimit = 10
	maxListLimit             = 100
)

func LimitDownloadAndStore(userID string) bool {
	amount, loaded := listLimit.LoadOrStore(userID, 1)
	if !loaded {
		return false
	}

	if amount.(int) == maxDownloadAndStoreLimit {
		// limit reached
		return true
	}

	listLimit.Store(userID, amount.(int)+1)
	return false
}

// call this function when download or store is done
func LimitDownloadAndStoreDone(userID string) {
	amount, loaded := listLimit.Load(userID)
	if !loaded {
		return
	}

	if amount.(int) == 1 {
		listLimit.Delete(userID)
		return
	}

	listLimit.Store(userID, amount.(int)-1)
}

// true if limit is reached
func LimitList(userID string) bool {
	amount, loaded := downloadAndStoreLimit.LoadOrStore(userID, 1)
	if !loaded {
		return false
	}

	if amount.(int) == maxListLimit {
		// limit reached
		return true
	}

	downloadAndStoreLimit.Store(userID, amount.(int)+1)
	return false
}

// call this function when list is done
func LimitListDone(userID string) {
	amount, loaded := downloadAndStoreLimit.Load(userID)
	if !loaded {
		return
	}

	if amount.(int) == 1 {
		downloadAndStoreLimit.Delete(userID)
		return
	}

	downloadAndStoreLimit.Store(userID, amount.(int)-1)
}
