package handler

import (
	"annisa-api/auth"
	"annisa-api/database"
	"annisa-api/middleware"
	"annisa-api/repository"
	"annisa-api/service"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartApp() {
	db, err := database.InitDb()
	if err != nil {
		log.Fatal("Eror Db Connection")
	}

	secretKey := os.Getenv("SECRET_KEY")
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowHeaders:    []string{"Access-Control-Allow-Headers", "Access-Control-Allow-Origin", "Origin , Accept , X-Requested-With , Content-Type, Access-Control-Request-Method, Access-Control-Request-Headers, Authorization"},
		AllowMethods:    []string{"POST, OPTIONS, GET, PUT, DELETE"},
	}))

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	authService := auth.NewUserAuthService()
	authService.SetSecretKey(secretKey)
	userHandler := NewUserHandler(userService, authService)

	router.POST("/register/", userHandler.RegisterUser)
	router.POST("/login/", userHandler.Login)

	cabangRepository := repository.NewCabangRepository(db)
	cabangService := service.NewCabangService(cabangRepository)
	cabangHandler := NewCabangHandler(cabangService)

	cabang := router.Group("api/cabang")
	cabang.POST("/", middleware.AuthMiddleware(authService, userService), cabangHandler.Create)

	karyawanRepository := repository.NewKaryawanRepository(db)
	karyawanService := service.NewKaryawanService(karyawanRepository)
	karyawanService.StartCommissionScheduler()
	karyawanHandler := NewKaryawanHandler(karyawanService)

	karyawan := router.Group("karyawan")
	karyawan.POST("/", middleware.AuthMiddleware(authService, userService), karyawanHandler.Create)
	karyawan.PUT("/:id", middleware.AuthMiddleware(authService, userService), karyawanHandler.Update)
	karyawan.GET("/id/:id", middleware.AuthMiddleware(authService, userService), karyawanHandler.GetByID)
	karyawan.GET("/:id_cabang", middleware.AuthMiddleware(authService, userService), karyawanHandler.GetByIDCabang)
	karyawan.DELETE("/:id", middleware.AuthMiddleware(authService, userService), karyawanHandler.Delete)

	memberRepository := repository.NewMemberRepository(db)
	memberService := service.NewMemberService(memberRepository)
	memberHandler := NewMemberHandler(memberService)

	member := router.Group("member")
	member.POST("/", middleware.AuthMiddleware(authService, userService), memberHandler.Create)
	member.PUT("/:id", middleware.AuthMiddleware(authService, userService), memberHandler.Update)
	member.GET("/:id", middleware.AuthMiddleware(authService, userService), memberHandler.GetByID)
	router.GET("/members", middleware.AuthMiddleware(authService, userService), memberHandler.GetAll)
	member.GET("/cabang/:id_cabang", middleware.AuthMiddleware(authService, userService), memberHandler.GetMemberByIDCabang)
	member.DELETE("/:id", middleware.AuthMiddleware(authService, userService), memberHandler.Delete)

	layananRepository := repository.NewLayananRepository(db)
	layananService := service.NewLayananService(layananRepository)
	layananHandler := NewLayananHandler(layananService)

	layanan := router.Group("layanan")
	layanan.POST("/", middleware.AuthMiddleware(authService, userService), layananHandler.Create)
	layanan.PUT("/:id", middleware.AuthMiddleware(authService, userService), layananHandler.Update)
	layanan.GET("/:id", middleware.AuthMiddleware(authService, userService), layananHandler.GetByID)
	layanan.DELETE("/:id", middleware.AuthMiddleware(authService, userService), layananHandler.Delete)
	layanan.GET("/", middleware.AuthMiddleware(authService, userService), layananHandler.GetAll)

	transaksiRepository := repository.NewTransaksiRepository(db)
	itemTransaksiRepository := repository.NewItemTransaksiRepository(db)
	transaksiService := service.NewTransaksiService(db, transaksiRepository, cabangRepository, itemTransaksiRepository, layananRepository, karyawanRepository)
	transaksiHandler := NewHandlerTransaksi(db, transaksiService)

	router.POST("/transaksi", middleware.AuthMiddleware(authService, userService), transaksiHandler.AddTransaksi)
	router.GET("/transaksi/:id", middleware.AuthMiddleware(authService, userService), transaksiHandler.GetTransaksiByID)
	router.GET("/transaksi/date/:date/cabang/:id_cabang", middleware.AuthMiddleware(authService, userService), transaksiHandler.GetTransaksiByDateAndCabang)
	router.GET("/transaksi/month/:month/year/:year/cabang/:id_cabang", middleware.AuthMiddleware(authService, userService), transaksiHandler.GetMonthlyTransaksiByCabang)
	router.GET("/transaksi/draft/cabang/:id_cabang", middleware.AuthMiddleware(authService, userService), transaksiHandler.GetDraftTransaksiByCabang)
	router.DELETE("/transaksi/:id_transaksi", middleware.AuthMiddleware(authService, userService), transaksiHandler.DeleteTransaksi)
	router.GET("/money/date/:date/cabang/:id_cabang", transaksiHandler.GetTotalMoneyByDateAndCabang)
	router.GET("/total_money/month/:month/year/:year/cabang/:id_cabang", transaksiHandler.GetTotalMoneyByMonthAndYear)

	router.Run(":8080")
}
