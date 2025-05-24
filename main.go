package main

import (
	"backend-facturacion/mercadopago"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.POST("/API/v1/webhook")
	r.POST("/API/v1/payment", mercadopago.Payment)
	r.Run(":3050")

	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Modelos
type Cliente struct {
	ID        uint `gorm:"primaryKey"`
	Nombre    string
	Correo    string
	Telefono  string
	Direccion string
	Facturas  []Factura
}

type Factura struct {
	ID        uint `gorm:"primaryKey"`
	ClienteID uint
	Cliente   Cliente
	Fecha     string
	Total     float64
	Detalles  []DetalleFactura
	Pagos     []Pago
}

type DetalleFactura struct {
	ID        uint `gorm:"primaryKey"`
	FacturaID uint
	Producto  string
	Cantidad  int
	Precio    float64
}

type Pago struct {
	ID        uint `gorm:"primaryKey"`
	FacturaID uint
	Monto     float64
	FechaPago string
	Metodo    string
}

var db *gorm.DB

func initDatabase() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error cargando .env")
	}

	dsn := os.Getenv("DATABASE_URL")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}

	err = db.AutoMigrate(&Cliente{}, &Factura{}, &DetalleFactura{}, &Pago{})
	if err != nil {
		log.Fatal("Error en la migración de esquemas:", err)
	}
	fmt.Println("Migraciones realizadas con éxito")
}

func main() {
	initDatabase()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.Run(":8080")
}
