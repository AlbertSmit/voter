package main

// Helper to add subscription to Update struct.
func withUpdate(update Update, sub *Subscription) Update {
	update.Sub = sub
	return update
}