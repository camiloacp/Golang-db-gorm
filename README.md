# go-db-gorm

Proyecto de aprendizaje que reimplementa [`golang_db`](../golang_db) usando [GORM](https://gorm.io) como ORM, con soporte para **PostgreSQL** y **MySQL**.

---

## Estructura del proyecto

```
go-db-gorm/
├── main.go
├── model/
│   └── model.go          ← structs compartidos con tags GORM
├── storage/
│   ├── storage.go        ← conexión (PostgreSQL / MySQL)
│   ├── gorm_product.go
│   ├── gorm_invoiceheader.go
│   ├── gorm_invoiceitem.go
│   └── gorm_invoice.go
└── pkg/
    ├── product/
    ├── invoice/
    ├── invoiceheader/
    └── invoiceitem/
```

---

## Conexión a la base de datos

`storage/storage.go` expone un singleton configurable por motor:

```go
storage.New(storage.PostgreSQL) // o storage.MySQL
db := storage.DB()              // retorna *gorm.DB
```

### DSNs configurados

| Motor | DSN |
|---|---|
| PostgreSQL | `host=localhost user=golang_db_user password=golang_db_password dbname=godb port=7530 sslmode=disable TimeZone=UTC` |
| MySQL | `root:root@tcp(127.0.0.1:3306)/mysql-go?charset=utf8mb4&parseTime=True&loc=Local` |

> Referencia oficial: https://gorm.io/docs/connecting_to_the_database.html

---

## AutoMigrate

### Código

```go
storage.DB().AutoMigrate(
    &model.Product{},
    &model.InvoiceHeader{},
    &model.InvoiceItem{},
)
```

### Qué hace

`AutoMigrate` lee los structs de Go y sincroniza el esquema de la base de datos **sin destruir datos existentes**:

- **Si la tabla no existe** → la crea con todas sus columnas, índices y constraints.
- **Si la tabla ya existe** → agrega columnas o índices que falten. Nunca elimina columnas ni modifica tipos que ya existan.

Es el equivalente GORM de ejecutar manualmente:

```sql
-- PostgreSQL
CREATE TABLE IF NOT EXISTS products (
    id           SERIAL PRIMARY KEY,
    created_at   TIMESTAMP,
    updated_at   TIMESTAMP,
    deleted_at   TIMESTAMP,       -- soft delete
    name         VARCHAR(100) NOT NULL,
    observations VARCHAR(100),
    price        INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS invoice_headers (
    id         SERIAL PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    client     VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS invoice_items (
    id                 SERIAL PRIMARY KEY,
    created_at         TIMESTAMP,
    updated_at         TIMESTAMP,
    deleted_at         TIMESTAMP,
    invoice_header_id  INTEGER,
    product_id         INTEGER
);
```

### Por qué se pasan punteros (`&model.Product{}`)

GORM necesita la dirección del struct para inspeccionar su tipo mediante reflexión y derivar:
- El nombre de la tabla (`Product` → `products`)
- Los nombres de columna (snake_case de cada campo)
- Los tags `gorm:"..."` para tamaños, constraints y tipos

### Campos que agrega `gorm.Model`

Cada struct embebe `gorm.Model`, que inyecta automáticamente:

| Campo | Tipo | Comportamiento |
|---|---|---|
| `ID` | `uint` | Primary key autoincremental |
| `CreatedAt` | `time.Time` | GORM lo setea en `Create()` |
| `UpdatedAt` | `time.Time` | GORM lo actualiza en `Save()` / `Updates()` |
| `DeletedAt` | `gorm.DeletedAt` | Soft delete: marca el registro en vez de borrarlo. Las queries filtran `WHERE deleted_at IS NULL` automáticamente. Para borrado físico: `db.Unscoped().Delete(...)` |

### Asociaciones definidas en los modelos

```go
type Product struct {
    gorm.Model
    InvoiceItems []InvoiceItem   // has-many: un producto puede estar en varios items
}

type InvoiceHeader struct {
    gorm.Model
    InvoiceItems []InvoiceItem   // has-many: una cabecera tiene muchos items
}

type InvoiceItem struct {
    gorm.Model
    InvoiceHeaderID uint         // FK → invoice_headers.id (convención GORM)
    ProductID       uint         // FK → products.id
}
```

GORM resuelve las FK por convención: si el struct se llama `InvoiceHeader` y el campo es `InvoiceHeaderID`, GORM lo detecta como clave foránea sin configuración extra.

---

## Create — insertar registros

### Código

```go
product1 := model.Product{
    Name:  "Curso de Go",
    Price: 120,
}

obs := "Testing with Go"
product2 := model.Product{
    Name:         "Curso de Testing",
    Price:        150,
    Observations: &obs,
}

product3 := model.Product{
    Name:  "Curso de Python",
    Price: 250,
}

storage.DB().Create(&product1)
storage.DB().Create(&product2)
storage.DB().Create(&product3)
```

### Qué hace

`db.Create(&m)` ejecuta un `INSERT` y, al terminar, **escribe de vuelta en el struct** los valores generados por la base de datos:

- `m.ID` → el ID autoincremental asignado por la BD
- `m.CreatedAt` → timestamp seteado por GORM
- `m.UpdatedAt` → timestamp seteado por GORM

```go
storage.DB().Create(&product1)
fmt.Println(product1.ID) // ya tiene el ID real, p. ej. 1
```

### Por qué `Observations` es `*string` y no `string`

```go
Observations *string `gorm:"type:varchar(100)"`
```

Un campo `string` vacío (`""`) y un campo nulo (`NULL`) son cosas distintas en la BD. Usando puntero:

- `nil` → GORM inserta `NULL`
- `&obs` → GORM inserta el valor de la cadena

Con `string` plano no habría forma de distinguir "el usuario no envió el campo" de "el usuario envió una cadena vacía".

### Por qué se pasan punteros al Create (`&product1`)

GORM necesita la dirección del struct para:
1. Escribir el `ID`, `CreatedAt` y `UpdatedAt` de vuelta después del insert.
2. Leer los campos con reflexión para construir el `INSERT`.

Si se pasara por valor (`product1` sin `&`) los cambios se harían sobre una copia y el struct original no tendría el ID.

### SQL equivalente generado

```sql
-- product1 y product3 (Observations = nil → NULL)
INSERT INTO products (name, observations, price, created_at, updated_at)
VALUES ('Curso de Go', NULL, 120, NOW(), NOW());

-- product2 (Observations = &obs → valor)
INSERT INTO products (name, observations, price, created_at, updated_at)
VALUES ('Curso de Testing', 'Testing with Go', 150, NOW(), NOW());
```

---

## Create con asociaciones — insertar cabecera e ítems en una sola operación

### Código

```go
invoice := model.InvoiceHeader{
    Client: "Camilo Cortes",
    InvoiceItems: []model.InvoiceItem{
        {ProductID: 1},
        {ProductID: 2},
    },
}

storage.DB().Create(&invoice)
```

### Qué hace

GORM detecta que `InvoiceHeader` tiene una relación has-many con `InvoiceItem` e inserta los tres registros en una sola llamada:

```sql
INSERT INTO invoice_headers (client, created_at, updated_at)
VALUES ('Camilo Cortes', NOW(), NOW());

-- GORM toma el ID generado y lo usa como FK en los ítems
INSERT INTO invoice_items (invoice_header_id, product_id, created_at, updated_at)
VALUES (1, 1, NOW(), NOW());

INSERT INTO invoice_items (invoice_header_id, product_id, created_at, updated_at)
VALUES (1, 2, NOW(), NOW());
```

### Por qué funciona sin declarar `InvoiceHeaderID`

GORM resuelve la FK por convención: el campo `InvoiceHeaderID uint` en `InvoiceItem` coincide con el nombre del struct padre (`InvoiceHeader`) + `ID`. GORM llena ese campo automáticamente al insertar los ítems asociados.

### Por qué es mejor que insertar por separado

Insertar la cabecera primero y luego cada ítem requeriría capturar el ID generado manualmente y asignarlo a cada ítem. Con asociaciones GORM lo hace solo, el código queda más limpio y sin riesgo de inconsistencias si algún insert falla.

---

## Find — consultar productos

### Código

```go
products := make([]model.Product, 0)
storage.DB().Find(&products)

for _, product := range products {
    fmt.Printf("%d - %s\n", product.ID, product.Name)
}
```

### Qué hace

`db.Find(&slice)` ejecuta un `SELECT *` y rellena el slice con todos los registros de la tabla:

```sql
SELECT * FROM products WHERE deleted_at IS NULL;
```

El `WHERE deleted_at IS NULL` lo agrega GORM automáticamente por el soft delete de `gorm.Model`.

### Por qué `make([]model.Product, 0)` y no `var products []model.Product`

Ambas formas funcionan con GORM. La diferencia está en el valor que retornas si no hay registros:

| Declaración | Sin resultados retorna |
|---|---|
| `make([]model.Product, 0)` | `[]` — slice vacío (recomendado para APIs JSON) |
| `var products []model.Product` | `nil` |

Con `make` evitas que un `json.Marshal` serialice el resultado como `null` en vez de `[]`.

### Por qué se pasa puntero al slice (`&products`)

GORM necesita la dirección para poder escribir los registros dentro del slice original. Sin el `&` los datos se escribirían en una copia y `products` quedaría vacío.

### Variantes útiles de Find

```go
// Con condición WHERE
storage.DB().Where("price > ?", 100).Find(&products)

// Limitar columnas
storage.DB().Select("id", "name").Find(&products)

// Limit y Offset
storage.DB().Limit(10).Offset(0).Find(&products)

// Ordenar
storage.DB().Order("price desc").Find(&products)
```

---

## First — consultar un registro por ID

### Código

```go
myProduct := model.Product{}
storage.DB().First(&myProduct, 3)
fmt.Println(myProduct)
```

### Qué hace

`db.First(&m, id)` busca el registro con la primary key igual al `id` dado, ordenado por PK ascendente, y rellena el struct:

```sql
SELECT * FROM products WHERE id = 3 AND deleted_at IS NULL ORDER BY id LIMIT 1;
```

### Diferencia entre `First`, `Take` y `Last`

| Método | ORDER BY | Uso |
|---|---|---|
| `First(&m, id)` | PK `ASC` | El más común — busca por ID |
| `Take(&m, id)` | ninguno | Sin orden, ligeramente más rápido |
| `Last(&m, id)` | PK `DESC` | Trae el último registro |

### Manejo de errores — `ErrRecordNotFound`

Si el registro no existe, GORM retorna `gorm.ErrRecordNotFound`. Es importante chequearlo:

```go
result := storage.DB().First(&myProduct, 3)
if errors.Is(result.Error, gorm.ErrRecordNotFound) {
    fmt.Println("producto no encontrado")
}
```

`Find` en cambio **no** retorna error si no hay resultados — simplemente deja el slice vacío. Por eso `First` es preferible cuando buscas un registro específico y quieres saber si no existe.

### `fmt.Println(myProduct)` con `gorm.Model`

Como `model.Product` embebe `gorm.Model`, el print incluye todos los campos:

```
{{1 2024-01-01 00:00:00 +0000 UTC 2024-01-01 00:00:00 +0000 UTC {0001-01-01 false}} Curso de Python <nil> 250}
```

Para un print más limpio puedes implementar `String()` en el modelo o usar `%+v`.

---

## Updates — actualizar registros

### Código

```go
myProduct := model.Product{}
myProduct.ID = 4

storage.DB().Model(&myProduct).Updates(
    model.Product{Name: "Curso de Java", Price: 120})
```

### Qué hace

`db.Model(&m).Updates(values)` ejecuta un `UPDATE` sobre el registro identificado por la PK del struct y actualiza solo los campos **no-zero** del struct pasado a `Updates`:

```sql
UPDATE products
SET name = 'Curso de Java', price = 120, updated_at = NOW()
WHERE id = 4 AND deleted_at IS NULL;
```

GORM actualiza `UpdatedAt` automáticamente y aplica el filtro de soft delete (`deleted_at IS NULL`).

### Por qué se separa la asignación del ID

```go
myProduct := model.Product{}
myProduct.ID = 4
```

Se crea un struct vacío y se le asigna solo el `ID` para que GORM sepa qué registro actualizar sin necesidad de hacer un `SELECT` previo. Es más eficiente que usar `First` cuando ya conoces la PK.

### Campos zero-value son ignorados

`Updates` con struct ignora los campos con valor zero (`""`, `0`, `false`, `nil`). Si `Price` fuera `0`, GORM **no** lo actualizaría. Para forzar la actualización de campos zero usa `map[string]interface{}`:

```go
storage.DB().Model(&myProduct).Updates(map[string]interface{}{
    "name":  "Curso de Java",
    "price": 0,
})
```

### Diferencia entre `Save` y `Updates`

| Método | Comportamiento |
|---|---|
| `Save(&m)` | `UPDATE` de **todos** los campos (incluye zeros) |
| `Updates(struct)` | `UPDATE` solo de campos **no-zero** |
| `Updates(map)` | `UPDATE` de exactamente los campos del mapa |

---

## Delete — eliminar registros (soft delete)

### Código

```go
myProduct := model.Product{}
myProduct.ID = 4

storage.DB().Delete(&myProduct)
```

### Qué hace

`db.Delete(&m)` ejecuta un soft delete: **no borra la fila**, sino que setea `deleted_at` con el timestamp actual:

```sql
UPDATE products
SET deleted_at = NOW()
WHERE id = 4 AND deleted_at IS NULL;
```

A partir de ese momento, GORM excluye automáticamente el registro en todas las queries (`Find`, `First`, `Updates`, etc.) porque agrega `WHERE deleted_at IS NULL`.

### Por qué es soft delete y no borrado físico

`gorm.Model` embebe el campo `DeletedAt gorm.DeletedAt`. Cuando ese campo existe en el struct, GORM activa el soft delete de forma automática — no requiere configuración adicional.

Ventajas:
- Los datos se pueden recuperar
- Se mantiene integridad referencial con otras tablas
- Útil para auditoría

### Borrado físico (`Unscoped`)

Para eliminar la fila permanentemente de la base de datos:

```go
storage.DB().Unscoped().Delete(&myProduct)
```

```sql
DELETE FROM products WHERE id = 4;
```

### Diferencia entre soft delete y borrado físico

| Método | SQL generado | Recuperable |
|---|---|---|
| `Delete(&m)` | `UPDATE ... SET deleted_at = NOW()` | Sí |
| `Unscoped().Delete(&m)` | `DELETE FROM ...` | No |

---

## Delete permanente — borrado físico (`Unscoped`)

### Código

```go
myProduct := model.Product{}
myProduct.ID = 4

storage.DB().Unscoped().Delete(&myProduct)
```

### Qué hace

`db.Unscoped().Delete(&m)` elimina la fila físicamente de la base de datos, ignorando el campo `deleted_at`:

```sql
DELETE FROM products WHERE id = 4;
```

A diferencia del soft delete, este registro **no se puede recuperar**.

### Por qué se necesita `Unscoped`

Sin `Unscoped`, GORM siempre aplica soft delete cuando el modelo embebe `gorm.Model`. `Unscoped` desactiva ese comportamiento para la query puntual, tanto en lecturas como en escrituras:

```go
// Leer registros ya eliminados (soft deleted)
storage.DB().Unscoped().Find(&products)

// Borrar físicamente
storage.DB().Unscoped().Delete(&myProduct)
```

### Cuándo usar borrado físico

- Cumplimiento normativo (GDPR, "derecho al olvido")
- Limpiar registros de prueba o datos temporales
- Tablas de log donde el volumen importa y la recuperación no es necesaria

---

## Contenedores Docker usados

```bash
# MySQL 8.0
docker run --name mysql-go \
  -e MYSQL_ROOT_PASSWORD=root \
  -e MYSQL_DATABASE=mysql-go \
  -p 3306:3306 \
  -d mysql:8.0 \
  --default-authentication-plugin=mysql_native_password
```

> `--default-authentication-plugin=mysql_native_password`: fuerza el plugin de autenticación clásico de MySQL para compatibilidad con el driver Go (`go-sql-driver/mysql`). Sin este flag, MySQL 8 usa `caching_sha2_password` por defecto y la conexión puede fallar.
