package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
)

func main() {

	loadEnvFile(".env")


	InitDB()


	uploadsDir := os.Getenv("UPLOADS_DIR")
	if uploadsDir == "" {
		uploadsDir = "./uploads"
	}
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}


	app := fiber.New(fiber.Config{
		BodyLimit:    20 * 1024 * 1024,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 120 * time.Second,
		ProxyHeader:  fiber.HeaderXForwardedFor,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		},
	})


	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:5173,http://localhost:3000,http://18.209.173.4.nip.io,https://18.209.173.4.nip.io"
	}

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	app.Use(limiter.New(limiter.Config{
		Max:        100,
		Expiration: 1 * time.Minute,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(429).JSON(fiber.Map{"error": "Too many requests. Please try again later."})
		},
	}))


	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"status": "healthy"})
	})

	auth := app.Group("/api/auth")
	auth.Post("/google", HandleGoogleAuth)
	auth.Get("/me", AuthMiddleware, HandleGetMe)
	auth.Post("/logout", HandleLogout)


	api := app.Group("/api", AuthMiddleware)


	api.Post("/analyze", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*JWTClaims)


		companyName := c.FormValue("companyName")
		jobTitle := c.FormValue("jobTitle")
		jobDesc := c.FormValue("jobDescription")

		if companyName == "" || jobTitle == "" {
			return c.Status(400).JSON(fiber.Map{"error": "companyName and jobTitle are required"})
		}


		fileHeader, err := c.FormFile("resume")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "No PDF file uploaded. Field name must be 'resume'"})
		}


		if filepath.Ext(fileHeader.Filename) != ".pdf" {
			return c.Status(400).JSON(fiber.Map{"error": "Only PDF files are accepted"})
		}


		file, err := fileHeader.Open()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to read uploaded file"})
		}
		defer file.Close()

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to read file content"})
		}


		reader := bytes.NewReader(fileBytes)
		resumeText, err := ExtractTextFromPDF(reader, int64(len(fileBytes)))
		if err != nil {
			return c.Status(422).JSON(fiber.Map{
				"error": fmt.Sprintf("Could not extract text from PDF: %s. Please ensure the PDF is not a scanned image.", err.Error()),
			})
		}


		resumeID := uuid.New().String()
		savePath := filepath.Join(uploadsDir, resumeID+".pdf")
		if err := os.WriteFile(savePath, fileBytes, 0644); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save uploaded file"})
		}


		feedback, err := AnalyzeResume(resumeText, jobTitle, jobDesc)
		if err != nil {

			os.Remove(savePath)
			return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("AI analysis failed: %s", err.Error())})
		}


		feedbackJSON, err := json.Marshal(feedback)
		if err != nil {
			os.Remove(savePath)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to serialize feedback"})
		}


		resume := Resume{
			ID:          resumeID,
			UserID:      user.UserID,
			CompanyName: companyName,
			JobTitle:    jobTitle,
			JobDesc:     jobDesc,
			FilePath:    savePath,
			FeedbackRaw: string(feedbackJSON),
		}
		if err := DB.Create(&resume).Error; err != nil {
			os.Remove(savePath)
			return c.Status(500).JSON(fiber.Map{"error": "Failed to save resume to database"})
		}

		resume.Feedback = feedback
		return c.Status(201).JSON(resume)
	})


	api.Get("/resumes", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*JWTClaims)

		var resumes []Resume
		if err := DB.Where("user_id = ?", user.UserID).
			Order("created_at desc").
			Find(&resumes).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch resumes"})
		}

		return c.JSON(resumes)
	})


	api.Get("/resume/:id", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*JWTClaims)
		resumeID := c.Params("id")

		var resume Resume
		if err := DB.Where("id = ? AND user_id = ?", resumeID, user.UserID).First(&resume).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Resume not found"})
		}

		return c.JSON(resume)
	})


	api.Get("/resume/:id/download", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*JWTClaims)
		resumeID := c.Params("id")

		var resume Resume
		if err := DB.Where("id = ? AND user_id = ?", resumeID, user.UserID).First(&resume).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Resume not found"})
		}

		fileInfo, err := os.Stat(resume.FilePath)
		if err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "PDF file not found on server"})
		}

		c.Set("Content-Type", "application/pdf")
		c.Set("Content-Disposition", fmt.Sprintf("inline; filename=\"resume-%s.pdf\"", resumeID))
		c.Set("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
		return c.SendFile(resume.FilePath)
	})


	api.Delete("/resume/:id", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*JWTClaims)
		resumeID := c.Params("id")

		var resume Resume
		if err := DB.Where("id = ? AND user_id = ?", resumeID, user.UserID).First(&resume).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{"error": "Resume not found"})
		}


		os.Remove(resume.FilePath)


		DB.Delete(&resume)

		return c.JSON(fiber.Map{"message": "Resume deleted successfully"})
	})


	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server running on http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}



func loadEnvFile(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return
	}

	lines := splitLines(string(data))
	for _, line := range lines {

		if len(line) == 0 || line[0] == '#' {
			continue
		}
		for i, ch := range line {
			if ch == '=' {
				key := line[:i]
				val := line[i+1:]

				if os.Getenv(key) == "" {
					os.Setenv(key, val)
				}
				break
			}
		}
	}
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i, c := range s {
		if c == '\n' {
			line := s[start:i]
			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			lines = append(lines, line)
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
