package database

// func CreateTables(db *sql.DB, dbConfig config.Config) error {
// 	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbConfig.DBName))
// 	if err != nil {
// 		return fmt.Errorf("failed to create database: %v", err)
// 	}

// 	_, err = db.Exec(fmt.Sprintf("USE %s", dbConfig.DBName))
// 	if err != nil {
// 		return fmt.Errorf("failed to switch to database: %v", err)
// 	}

// 	log.Printf("✅ Database '%s' is ready!", dbConfig.DBName)

// 	tables := []string{
// 		`CREATE TABLE IF NOT EXISTS users (
// 			id INT AUTO_INCREMENT PRIMARY KEY,
// 			username VARCHAR(50) UNIQUE NOT NULL,
// 			email VARCHAR(100) UNIQUE NOT NULL,
// 			password VARCHAR(255) NOT NULL,
// 			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
// 		);`,
// 		`CREATE TABLE IF NOT EXISTS posts (
// 			id INT AUTO_INCREMENT PRIMARY KEY,
// 			user_id INT NOT NULL,
// 			title VARCHAR(255) NOT NULL,
// 			content TEXT NOT NULL,
// 			views INT DEFAULT 0,
// 			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// 			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
// 		);`,
// 		`CREATE TABLE IF NOT EXISTS comments (
// 			id INT AUTO_INCREMENT PRIMARY KEY,
// 			post_id INT NOT NULL,
// 			user_id INT NOT NULL,
// 			content TEXT NOT NULL,
// 			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// 			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
// 			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
// 		);`,
// 		`CREATE TABLE IF NOT EXISTS likes (
// 			id INT AUTO_INCREMENT PRIMARY KEY,
// 			post_id INT NOT NULL,
// 			user_id INT NOT NULL,
// 			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
// 			FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
// 			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
// 		);`,
// 	}

// 	for _, table := range tables {
// 		_, err := db.Exec(table)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
