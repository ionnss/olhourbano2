package db

// Migration represents a databse migrations
type Migration struct {
	ID   string
	File string
	SQL  string
}

// RunMigragtions runs the migrations files in order
func RunMigrations() error {
	return nil
}
