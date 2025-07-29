package database

import (
	"log"
)

// SetupTables creates all the necessary tables in the correct order to respect foreign key constraints
func SetupTables() {
	if Client == nil {
		log.Fatal("Database client is not initialized")
	}

	// Create tables in order of dependencies (parent tables first)
	createRestaurantsTable() // No dependencies
	createUsersTable()       // References restaurants
	createMenusTable()       // References restaurants
	createTablesTable()      // References restaurants
	createFoodsTable()       // References menus and restaurants
	createOrdersTable()      // References tables and restaurants
	createOrderItemsTable()  // References orders and foods
	createNotesTable()       // References restaurants and orders
	createInvoicesTable()    // References orders and restaurants
}

func createFoodsTable() {
	query := `CREATE TABLE IF NOT EXISTS foods (
		id SERIAL PRIMARY KEY, 
		name VARCHAR(100) NOT NULL, 
		price NUMERIC(10, 2) NOT NULL, 
		description TEXT, 
		image_url VARCHAR(255), 
		menu_id INTEGER NOT NULL REFERENCES menus(id) ON DELETE CASCADE, 
		restaurant_id INTEGER NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
		ingredients TEXT,
		prep_time INTEGER DEFAULT NULL,
		calories INTEGER DEFAULT NULL,
		spicy_level INTEGER DEFAULT NULL,
		vegetarian BOOLEAN DEFAULT TRUE,
		available BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, 
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`

	if _, err := Client.Exec(query); err != nil {
		log.Printf("Error creating foods table: %v", err)
	}
}

func createMenusTable() {
	query := `CREATE TABLE IF NOT EXISTS menus (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) NOT NULL CHECK (name IN ('appetizer', 'main_course', 'dessert', 'beverage')),
		restaurant_id INTEGER NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
		category VARCHAR(50),
		description TEXT,
		active BOOLEAN DEFAULT TRUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`

	if _, err := Client.Exec(query); err != nil {
		log.Printf("Error creating menus table: %v", err)
	}
}

func createOrdersTable() {
	query := `CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		table_id INTEGER REFERENCES tables(id) ON DELETE CASCADE,
		restaurant_id INTEGER NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
		order_date TIMESTAMP,
		total_price NUMERIC(10, 2),
		status VARCHAR(10) CHECK (status IN ('pending', 'preparing', 'ready', 'served', 'paid', 'cancelled')) DEFAULT 'pending',
		notes TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`

	if _, err := Client.Exec(query); err != nil {
		log.Printf("Error creating orders table: %v", err)
	}
}

func createOrderItemsTable() {
	query := `CREATE TABLE IF NOT EXISTS orderitems (
		id SERIAL PRIMARY KEY,
		order_id INTEGER NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
		food_id INTEGER NOT NULL REFERENCES foods(id) ON DELETE CASCADE,
		quantity INTEGER DEFAULT 1,
		unit_price NUMERIC(5, 2),
		subtotal NUMERIC(6, 2),
		notes TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`

	if _, err := Client.Exec(query); err != nil {
		log.Printf("Error creating order_items table: %v", err)
	}
}

func createTablesTable() {
	query := `CREATE TABLE IF NOT EXISTS tables (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		table_number INTEGER NOT NULL,
		capacity INTEGER NOT NULL,
		restaurant_id INTEGER NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
		location VARCHAR(255),
		status VARCHAR(10) NOT NULL CHECK (status IN ('available', 'occupied', 'reserved')) DEFAULT 'available',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`

	if _, err := Client.Exec(query); err != nil {
		log.Printf("Error creating tables table: %v", err)
	}
}

func createUsersTable() {
	query := `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		email VARCHAR(50) NOT NULL UNIQUE,
		password VARCHAR(60) NOT NULL,
		phone VARCHAR(20) NOT NULL,
		role VARCHAR(10) NOT NULL,
		token TEXT,
		avatar_url TEXT,
		restaurant_id INTEGER REFERENCES restaurants(id) ON DELETE SET NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`

	if _, err := Client.Exec(query); err != nil {
		log.Printf("Error creating users table: %v", err)
	}
}

func createRestaurantsTable() {
	query := `CREATE TABLE IF NOT EXISTS restaurants (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		owner_id INTEGER NOT NULL,
		logo VARCHAR(255),
		address VARCHAR(500) NOT NULL,
		phone VARCHAR(20),
		email VARCHAR(100),
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`

	if _, err := Client.Exec(query); err != nil {
		log.Printf("Error creating restaurants table: %v", err)
	}
}

func createNotesTable() {
	query := `CREATE TABLE IF NOT EXISTS notes (
		id SERIAL PRIMARY KEY,
		title VARCHAR(100) NOT NULL,
		content TEXT NOT NULL,
		priority VARCHAR(10) NOT NULL CHECK (priority IN ('low', 'medium', 'high')),
		restaurant_id INTEGER NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
		order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`

	if _, err := Client.Exec(query); err != nil {
		log.Printf("Error creating notes table: %v", err)
	}
}

func createInvoicesTable() {
	query := `CREATE TABLE IF NOT EXISTS invoices (
		id SERIAL PRIMARY KEY,
		order_id INTEGER UNIQUE REFERENCES orders(id) ON DELETE CASCADE,
		restaurant_id INTEGER REFERENCES restaurants(id) ON DELETE CASCADE,
		amount NUMERIC(10, 2),
		tax NUMERIC(10, 2),
		total NUMERIC(10, 2),
		tax_amount NUMERIC(10, 2) DEFAULT 0,
		discount_amount NUMERIC(10, 2) DEFAULT 0,
		status VARCHAR(10) CHECK (status IN ('pending', 'paid')) DEFAULT 'pending',
		payment_method VARCHAR(20) CHECK (payment_method IN ('cash', 'credit_card', 'debit_card', 'online')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP)`

	if _, err := Client.Exec(query); err != nil {
		log.Printf("Error creating invoices table: %v", err)
	}
}
