package main

import (
	"database/sql"
	"log"

	_ "primerProjecto/docs"

	swaggerFiles "github.com/swaggo/files"

	controllers "primerProjecto/internal/adapters/controllers"
	"primerProjecto/internal/adapters/cotizadores"
	repositories "primerProjecto/internal/adapters/repositories"
	"primerProjecto/internal/services"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title     Cripto Api
//@version 1.0
//@description app para cotizaciones de criptos

// host localhost:8080
func main() {

	router := gin.Default()

	// Configurar ruta para la documentación JSON de Swagger
	router.GET("/swagger.json", func(c *gin.Context) {
		c.File("primerProjecto/docs/swagger.json")
	})

	// Configurar la conexión a la base de datos MySQL
	db, err := sql.Open("mysql", "root:1234@tcp(172.18.224.1:3306)/?parseTime=true")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Intentar crear la base de datos (opcional si ya existe)
	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS proyecto_cripto")
	if err != nil {
		log.Fatal(err)
	}

	// Seleccionar la base de datos recién creada o existente
	_, err = db.Exec("USE proyecto_cripto")
	if err != nil {
		log.Fatal(err)
	}

	// Crear tabla monedas
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS monedas (
        id INT AUTO_INCREMENT PRIMARY KEY,
        nombre VARCHAR(100) NOT NULL
    )
`)
	if err != nil {
		log.Fatal(err)
	}

	// Crear tabla cotizaciones
	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS cotizaciones (
    id INT AUTO_INCREMENT PRIMARY KEY,
    cripto_id INT NOT NULL,
    cotizacion DECIMAL(10, 2) NOT NULL,
    fecha DATETIME NOT NULL,
    manual BOOLEAN NOT NULL DEFAULT FALSE,
    usuario_id INT DEFAULT NULL,
    FOREIGN KEY (cripto_id) REFERENCES monedas(id),
    FOREIGN KEY (usuario_id) REFERENCES usuarios(id)
)
`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS usuarios (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL,
    apellidos VARCHAR(100) NOT NULL,
    fecha_nacimiento DATE NOT NULL,
    codigo_usuario VARCHAR(100) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    tipo_documento ENUM('DNI', 'pasaporte', 'cedula') NOT NULL,
    fecha_registro DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    esta_activo BOOLEAN NOT NULL DEFAULT TRUE
    )
`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
    CREATE TABLE IF NOT EXISTS usuario_moneda (
            usuario_id INT NOT NULL,
            moneda_id INT NOT NULL,
            PRIMARY KEY (usuario_id, moneda_id),
            FOREIGN KEY (usuario_id) REFERENCES usuarios(id),
            FOREIGN KEY (moneda_id) REFERENCES monedas(id)
    )
`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS auditoria_cotizacion (
		id BIGINT AUTO_INCREMENT PRIMARY KEY,
		usuario_id INT,
		cotizacion_id INT,
		log TEXT NOT NULL,
		created_at DATETIME DEFAULT NOW(),
		FOREIGN KEY (usuario_id) REFERENCES usuarios(id) ON DELETE SET NULL,
		FOREIGN KEY (cotizacion_id) REFERENCES cotizaciones(id) ON DELETE SET NULL
	)
`)
	if err != nil {
		log.Fatal(err)
	}

	// Crear las instancias de los repositorios
	repoUsuario := repositories.NewMySQLUsuarioRepository(db)
	repoCripto := repositories.NewMySQLCryptoRepository(db)

	// Crear las instancias de los servicios usando las interfaces
	serviceUsuario := services.NewUsuarioService(repoUsuario, repoCripto)
	serviceCripto := services.NewCryptoService(repoCripto, cotizadores.GetCotizador)

	//handlers/controllers
	criptoHandler := controllers.NewCryptoController(serviceCripto)
	usuarioHandler := controllers.NewUsuarioHandler(serviceUsuario)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// Configurar tus rutas y controladores
	/*router.POST("/cryptocurrencies", controller.RegistrarCriptoMoneda)
	router.POST("/cotization", controller.RegistrarCotizacion)
	router.POST("/cryptocurrencies/externa", controller.SaveMonedaConCotizacion)
	router.POST("/cotization/externa", controller.SaveCotizacionExterna)*/

	//cotizaciones manuales
	router.POST("cotization/manual", usuarioHandler.RegistrarCotizacionManual)
	router.DELETE("cotizacion/manual/:id", usuarioHandler.BorrarCotizacionManual)
	router.PUT("cotizacion/manual/:usuarioId/:cotizacionId", usuarioHandler.ActualizarCotizacionManual)

	router.POST("/usuarios", usuarioHandler.CreateUsuario)
	router.PUT("/usuarios/:id", usuarioHandler.UpdateUsuarioByID)
	router.PUT("/usuarios/:id/monedasFavoritas", usuarioHandler.GuardarMonedaFavorita)
	router.PATCH("/usuarios/:id", usuarioHandler.PatchUsuarioByID)
	router.GET("/usuarios/:id", usuarioHandler.FindUsuarioByID)
	router.GET("/usuarios/:id/monedas", usuarioHandler.FindMonedasByUsuarioID)
	router.GET("/usuarios/:id/cotizaciones", criptoHandler.FindAllByFilterUsuario)
	router.POST("/upsert-usuario", usuarioHandler.UpsertUsuario)

	router.POST("/cryptocurrencies", services.AuthMiddleware(), criptoHandler.RegistrarCriptoMoneda)
	router.POST("/cotization", services.AuthMiddleware(), criptoHandler.RegistrarCotizacion)
	router.POST("/cryptocurrencies/externa", services.AuthMiddleware(), criptoHandler.SaveMonedaConCotizacion)
	router.POST("/cotization/externa", services.AuthMiddleware(), criptoHandler.SaveCotizacionExterna)

	router.GET("/cryptocurrencies/All", criptoHandler.FindAll)
	router.GET("/cryptocurrencies/cryptocurrency/:id", criptoHandler.FindMonedaByID)
	router.GET("/cryptocurrencies/:nombre/cryptocurrency", criptoHandler.FindMondaByNombre)
	router.GET("/cryptocurrencies", criptoHandler.FindAllByFilter)
	router.GET("/cryptocurrencies/lastcotization/:nombre", criptoHandler.FindUltimaCotizacion)
	router.PUT("/cryptocurrency/:id", criptoHandler.HandleUpdateCryptoByID)

	// Iniciar el servidor HTTP
	router.Run(":8080")
}
