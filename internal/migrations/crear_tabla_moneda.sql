CREATE TABLE IF NOT EXISTS monedas (
    Id INTEGER PRIMARY KEY AUTOINCREMENT,
    Nombre TEXT NOT NULL
);

-- Insertar los registros en la tabla CriptoMoneda
INSERT INTO monedas (nombre) VALUES ('Bitcoin','BTC');
INSERT INTO monedas (nombre) VALUES ('Ethereum','ETH');
INSERT INTO monedas (nombre) VALUES ('Ripple RXP 589','XRP');
INSERT INTO monedas (nombre) VALUES ('Litecoin','LTC');
INSERT INTO monedas (nombre) VALUES ('Bitcoin Cash','BHC');