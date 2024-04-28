package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
)

func main() {
    // Initialize the Gin router
    router := gin.Default()

    // Load environment variables
    dbURL := os.Getenv("DB_URL")
    if dbURL == "" {
        log.Fatal("DB_URL environment variable not set")
    }

    // Connect to the database
    db, err := sql.Open("postgres", dbURL)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Define routes
    router.LoadHTMLGlob("templates/*")
    router.GET("/", func(c *gin.Context) {
        c.HTML(http.StatusOK, "index.html", gin.H{})
    })

    router.POST("/register", func(c *gin.Context) {
        // Parse form data
        fullName := c.PostForm("fullName")
        cpf := c.PostForm("cpf")
        phoneNumber := c.PostForm("phoneNumber")
        anamnesis := c.PostForm("anamnesis")

        // Insert data into the database
        _, err := db.Exec("INSERT INTO patients (full_name, cpf, phone_number, anamnesis) VALUES ($1, $2, $3, $4)", fullName, cpf, phoneNumber, anamnesis)
        if err != nil {
            log.Fatal(err)
        }

        // Redirect to the home page after registration
        c.Redirect(http.StatusSeeOther, "/")
    })

    router.GET("/patients", func(c *gin.Context) {
        // Fetch the list of patients from the database
        rows, err := db.Query("SELECT id, full_name FROM patients")
        if err != nil {
            log.Fatal(err)
        }
        defer rows.Close()

        // Store the list of patients in a slice
        var patients []struct {
            ID       int
            FullName string
        }
        for rows.Next() {
            var p struct {
                ID       int
                FullName string
            }
            err := rows.Scan(&p.ID, &p.FullName)
            if err != nil {
                log.Fatal(err)
            }
            patients = append(patients, p)
        }

        // Render the patient list template
        c.HTML(http.StatusOK, "patients.html", gin.H{
            "patients": patients,
        })
    })

    router.GET("/patient/:id", func(c *gin.Context) {
        // Get the patient ID from the URL parameter
        id := c.Param("id")

        // Fetch the patient's details from the database based on the ID
        var fullName, cpf, phoneNumber, anamnesis string
        err := db.QueryRow("SELECT full_name, cpf, phone_number, anamnesis FROM patients WHERE id = $1", id).Scan(&fullName, &cpf, &phoneNumber, &anamnesis)
        if err != nil {
            log.Fatal(err)
        }

        // Render the patient details template
        c.HTML(http.StatusOK, "patient_details.html", gin.H{
            "fullName":    fullName,
            "cpf":         cpf,
            "phoneNumber": phoneNumber,
            "anamnesis":   anamnesis,
        })
    })

    // Start the web server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    router.Run(":" + port)
}
