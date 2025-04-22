package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitDb() (*sql.DB, error) {
	if _, exists := os.LookupEnv("RAILWAY_ENVIRONMENT"); exists == false {
		if err := godotenv.Load(); err != nil {
			log.Fatal("error loading .env file:", err)
		}
	}

	user := os.Getenv("MYSQLUSER")
	pass := os.Getenv("MYSQLPASSWORD")
	host := os.Getenv("MYSQLHOST")
	port := os.Getenv("MYSQLPORT")
	dbname := os.Getenv("MYSQLDATABASE")

	dsn := user + ":" + pass + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("DB Ping Error:", err)
		return nil, err
	}
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
		id_cabang INT,
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

	return db, nil
}
