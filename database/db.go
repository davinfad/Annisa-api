package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func InitDb() (*sql.DB, error) {
	if _, exists := os.LookupEnv("RAILWAY_ENVIRONMENT"); exists == false {
		if err := godotenv.Load(); err != nil {
			log.Fatal("error loading .env file:", err)
		}
	}

	dbUsername := os.Getenv("MYSQLUSER")
	dbPassword := os.Getenv("MYSQLPASSWORD")
	dbHost := os.Getenv("MYSQLHOST")
	dbPort := os.Getenv("MYSQLPORT")
	dbName := os.Getenv("MYSQLDATABASE")

	// dsn := dbUsername + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Asia%2FJakarta&multiStatements=true"

	dsn := dbUsername + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true"
	// dsn := "root:@tcp(127.0.0.1:3306)/annisa?parseTime=true&loc=loc=%2B07%3A00&multiStatements=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB Ping Error:", err)
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	query := `
	CREATE TABLE IF NOT EXISTS cabang (
		id_cabang INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		nama_cabang VARCHAR(255) NOT NULL,
		kode_cabang VARCHAR(50) NOT NULL,
		jam_buka TIME NOT NULL,
		jam_tutup TIME NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS karyawan (
		id_karyawan INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		nama_karyawan VARCHAR(255) NOT NULL,
		id_cabang INT,
		nomor_telepon VARCHAR(50),
		alamat VARCHAR(250),
		komisi DECIMAL(10,2) DEFAULT 0.00,
		komisi_harian DECIMAL(10,2) DEFAULT 0.00,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (id_cabang) REFERENCES cabang(id_cabang)
	);

	CREATE TABLE IF NOT EXISTS layanan (
		id_layanan INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		nama_layanan VARCHAR(255) NOT NULL,
		persen_komisi DECIMAL(5,2) NOT NULL,
		persen_komisi_luarjam DECIMAL(5,2) NOT NULL,
		kategori VARCHAR(255) NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS member (
		id_member INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		nomor_pelanggan VARCHAR(50) NOT NULL,
		nama_member VARCHAR(255) NOT NULL,
		nomor_telepon VARCHAR(15) NOT NULL,
		alamat TEXT NOT NULL,
		tanggal_lahir DATE NOT NULL,
		tanggal_daftar DATE NOT NULL,
		id_cabang INT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (id_cabang) REFERENCES cabang(id_cabang)
	);

	CREATE TABLE IF NOT EXISTS transaksi (
		id_transaksi INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		id_cabang INT,
		id_member INT,
		nama_pelanggan VARCHAR(255) NOT NULL,
		nomor_telepon VARCHAR(20) NOT NULL,
		total_harga DECIMAL(10,2) NOT NULL,
		diskon DECIMAL(5,2),
		metode_pembayaran VARCHAR(50) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		status TINYINT DEFAULT 0,
		FOREIGN KEY (id_cabang) REFERENCES cabang(id_cabang),
		FOREIGN KEY (id_member) REFERENCES member(id_member)
	);

	CREATE TABLE IF NOT EXISTS item_transaksi (
		id_item_transaksi INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		id_transaksi INT,
		id_karyawan INT,
		id_layanan INT,
		catatan TEXT,
		harga DECIMAL(10,2) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (id_transaksi) REFERENCES transaksi(id_transaksi),
		FOREIGN KEY (id_karyawan) REFERENCES karyawan(id_karyawan),
		FOREIGN KEY (id_layanan) REFERENCES layanan(id_layanan)
	);

	CREATE TABLE IF NOT EXISTS users (
		username VARCHAR(255) NOT NULL PRIMARY KEY,
		password VARCHAR(255) NOT NULL,
		access_code VARCHAR(15),
		id_cabang INT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (id_cabang) REFERENCES cabang(id_cabang)
	);

	CREATE TABLE IF NOT EXISTS inventory (
		id_inventory INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
		id_cabang INT NOT NULL,
		nama_item VARCHAR(255) NOT NULL,
		batas_bawah INT NOT NULL DEFAULT 0,
		batas_atas INT NOT NULL DEFAULT 0,
		stok INT NOT NULL DEFAULT 0,
		satuan VARCHAR(50) NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (id_cabang) REFERENCES cabang(id_cabang)
	);
	`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating tables: %v", err)
	}
	log.Println("Tables created or already exist.")

	migrations := []string{
		`ALTER TABLE users ADD COLUMN access_code VARCHAR(15)`,
		`ALTER TABLE users ADD COLUMN created_at DATETIME NOT NULL DEFAULT NOW()`,
		`ALTER TABLE users ADD COLUMN updated_at DATETIME NOT NULL DEFAULT NOW()`,
		`ALTER TABLE transaksi ADD COLUMN diskon DECIMAL(5,2)`,
		`ALTER TABLE cabang ADD COLUMN created_at DATETIME NOT NULL DEFAULT NOW()`,
		`ALTER TABLE cabang ADD COLUMN updated_at DATETIME NOT NULL DEFAULT NOW()`,
		`ALTER TABLE karyawan ADD COLUMN created_at DATETIME NOT NULL DEFAULT NOW()`,
		`ALTER TABLE karyawan ADD COLUMN updated_at DATETIME NOT NULL DEFAULT NOW()`,
		`ALTER TABLE layanan ADD COLUMN created_at DATETIME NOT NULL DEFAULT NOW()`,
		`ALTER TABLE layanan ADD COLUMN updated_at DATETIME NOT NULL DEFAULT NOW()`,
		`ALTER TABLE member ADD COLUMN created_at DATETIME NOT NULL DEFAULT NOW()`,
		`ALTER TABLE member ADD COLUMN updated_at DATETIME NOT NULL DEFAULT NOW()`,
	}
	for _, m := range migrations {
		if _, err := db.Exec(m); err != nil {
			log.Printf("Migration skipped (column likely exists): %v", err)
		}
	}

	return db, nil
}
