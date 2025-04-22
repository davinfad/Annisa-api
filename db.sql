CREATE TABLE `cabang` ( 
  `id_cabang` INT AUTO_INCREMENT NOT NULL,
  `nama_cabang` VARCHAR(255) NOT NULL,
  `kode_cabang` VARCHAR(50) NOT NULL,
  `jam_buka` TIME NOT NULL,
  `jam_tutup` TIME NOT NULL,
  CONSTRAINT `PRIMARY` PRIMARY KEY (`id_cabang`)
)
ENGINE = InnoDB;
CREATE TABLE `item_transaksi` ( 
  `id_item_transaksi` INT AUTO_INCREMENT NOT NULL,
  `id_transaksi` INT NULL,
  `id_karyawan` INT NULL,
  `id_layanan` INT NULL,
  `catatan` TEXT NULL,
  `harga` DECIMAL(10,2) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ,
  CONSTRAINT `PRIMARY` PRIMARY KEY (`id_item_transaksi`)
)
ENGINE = InnoDB;
CREATE TABLE `karyawan` ( 
  `id_karyawan` INT AUTO_INCREMENT NOT NULL,
  `nama_karyawan` VARCHAR(255) NOT NULL,
  `id_cabang` INT NULL,
  `nomor_telepon` VARCHAR(50) NULL,
  `alamat` VARCHAR(250) NULL,
  `komisi` DECIMAL(10,2) NULL DEFAULT 0.00 ,
  `komisi_harian` DECIMAL(10,2) NULL DEFAULT 0.00 ,
  CONSTRAINT `PRIMARY` PRIMARY KEY (`id_karyawan`)
)
ENGINE = InnoDB;
CREATE TABLE `layanan` ( 
  `id_layanan` INT AUTO_INCREMENT NOT NULL,
  `nama_layanan` VARCHAR(255) NOT NULL,
  `persen_komisi` DECIMAL(5,2) NOT NULL,
  `persen_komisi_luarjam` DECIMAL(5,2) NOT NULL,
  `kategori` VARCHAR(255) NOT NULL,
  CONSTRAINT `PRIMARY` PRIMARY KEY (`id_layanan`)
)
ENGINE = InnoDB;
CREATE TABLE `member` ( 
  `id_member` INT AUTO_INCREMENT NOT NULL,
  `nomor_pelanggan` VARCHAR(50) NOT NULL,
  `nama_member` VARCHAR(255) NOT NULL,
  `nomor_telepon` VARCHAR(15) NOT NULL,
  `alamat` TEXT NOT NULL,
  `tanggal_lahir` DATE NOT NULL,
  `tanggal_daftar` DATE NOT NULL,
  `id_cabang` INT NOT NULL,
  CONSTRAINT `PRIMARY` PRIMARY KEY (`id_member`)
)
ENGINE = InnoDB;
CREATE TABLE `transaksi` ( 
  `id_transaksi` INT AUTO_INCREMENT NOT NULL,
  `id_cabang` INT NULL,
  `id_member` INT NULL,
  `nama_pelanggan` VARCHAR(255) NOT NULL,
  `nomor_telepon` VARCHAR(20) NOT NULL,
  `total_harga` DECIMAL(10,2) NOT NULL,
  `metode_pembayaran` VARCHAR(50) NOT NULL,
  `created_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ,
  `status` TINYINT NULL DEFAULT 0 ,
  CONSTRAINT `PRIMARY` PRIMARY KEY (`id_transaksi`)
)
ENGINE = InnoDB;
CREATE TABLE `users` ( 
  `username` VARCHAR(255) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `id_cabang` INT NULL,
  CONSTRAINT `PRIMARY` PRIMARY KEY (`username`)
)
ENGINE = InnoDB;
ALTER TABLE `item_transaksi` ADD CONSTRAINT `fk_id_karyawan` FOREIGN KEY (`id_karyawan`) REFERENCES `karyawan` (`id_karyawan`) ON DELETE SET NULL ON UPDATE NO ACTION;
ALTER TABLE `item_transaksi` ADD CONSTRAINT `item_transaksi_ibfk_1` FOREIGN KEY (`id_transaksi`) REFERENCES `transaksi` (`id_transaksi`) ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE `item_transaksi` ADD CONSTRAINT `item_transaksi_ibfk_2` FOREIGN KEY (`id_karyawan`) REFERENCES `karyawan` (`id_karyawan`) ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE `item_transaksi` ADD CONSTRAINT `item_transaksi_ibfk_3` FOREIGN KEY (`id_layanan`) REFERENCES `layanan` (`id_layanan`) ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE `karyawan` ADD CONSTRAINT `karyawan_ibfk_1` FOREIGN KEY (`id_cabang`) REFERENCES `cabang` (`id_cabang`) ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE `member` ADD CONSTRAINT `member_ibfk_1` FOREIGN KEY (`id_cabang`) REFERENCES `cabang` (`id_cabang`) ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE `transaksi` ADD CONSTRAINT `fk_cabang` FOREIGN KEY (`id_cabang`) REFERENCES `cabang` (`id_cabang`) ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE `transaksi` ADD CONSTRAINT `transaksi_ibfk_1` FOREIGN KEY (`id_member`) REFERENCES `member` (`id_member`) ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE `users` ADD CONSTRAINT `users_ibfk_1` FOREIGN KEY (`id_cabang`) REFERENCES `cabang` (`id_cabang`) ON DELETE NO ACTION ON UPDATE NO ACTION;
