package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	_ "os"
	"time"

	_ "github.com/lib/pq"
)

const (
	dbHost     = "localhost"
	dbPort     = 5432
	dbUser     = "postgres"
	dbPassword = "aldizix6ZWY"
	dbName     = "assignmentsTiga"
)

func main() {
	// Membuat koneksi ke database PostgreSQL
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	// Memastikan koneksi ke database berhasil
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping the database:", err)
	}

	// Memulai proses update data setiap 15 detik
	ticker := time.NewTicker(15 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				updateData(db)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	// Mengakhiri program dengan menunggu pengguna menekan tombol 'Enter'
	fmt.Println("Press Enter to stop the program...")
	fmt.Scanln()
	quit <- struct{}{}
}

func updateData(db *sql.DB) {
	// Menghasilkan angka acak antara 1 hingga 100 untuk water dan wind
	water := rand.Intn(100) + 1
	wind := rand.Intn(100) + 1

	// Mengubah data di database menggunakan API
	_, err := db.Exec("UPDATE update SET water = $1, wind = $2 WHERE id = $3", water, wind, 1)
	if err != nil {
		log.Println("Failed to update data:", err)
		return
	}

	// Menyimpan data yang diperbarui ke tabel di PostgreSQL
	_, err = db.Exec("INSERT INTO update (water, wind) VALUES ($1, $2)", water, wind)
	if err != nil {
		log.Println("Failed to insert data:", err)
		return
	}

	// Menampilkan status dari hasil update data
	log.Println("Water:", water, "- Status:", getWaterStatus(water))
	log.Println("Wind:", wind, "- Status:", getWindStatus(wind))
}

func getWaterStatus(water int) string {
	if water < 5 {
		return "Aman"
	} else if water >= 6 && water <= 8 {
		return "Siaga"
	} else {
		return "Bahaya"
	}
}

func getWindStatus(wind int) string {
	if wind < 6 {
		return "Aman"
	} else if wind >= 7 && wind <= 15 {
		return "Siaga"
	} else {
		return "Bahaya"
	}
}
