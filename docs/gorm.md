

## Valores devueltos para el select
- First(&obj)    
    gorm.ErrRecordNotFound.
    No modifica obj, pero queda con sus valores por defecto
    Ej:
    err := db.First(&user, "name = ?", "NoExiste").Error
	if err == gorm.ErrRecordNotFound { ... }

- Find
    devuelve un slice vacío

## Acceder al query SQL real
### Habilitar el modo de depuración para imprimirlo en la consola
```go
db.Debug().Find(&users)
```
### db.Statement.SQL.String()
```go
stmt := db.Session(&gorm.Session{DryRun: true}).Find(&users).Statement
fmt.Println(stmt.SQL.String()) // Muestra la consulta SQL
fmt.Println(stmt.Vars)         // Muestra los valores de los placeholders
o
tx := db.Session(&gorm.Session{DryRun: true}).Find(&users)
fmt.Println(tx.Statement.SQL.String()) // Obtiene el SQL sin ejecutarlo
```
### Configurar un Logger personalizado para ver siempre las consultas
```go
import (
    "gorm.io/gorm/logger"
    "time"
)
db, _ := gorm.Open(sqlite.Open("test.db"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info),})
```


## Generar foreign keys

Hay que incluir un objeto de la entidad foránea para que se cree la resticcion en base de datos (1 a 1 o 1 a n) o un array de objetos hijo en el lado muchos (1 a n)
  
``` go
// No crea foreign key  
type EntA struct {
    StatusID  int 
}
//Sí crea foreign key  
type EntA struct {
    StatusID  int   //  <<<--- fk
    Status    Status
}

//Sí crea foreign key  
type EntA struct {
    EntBID  int  //  <<<--- fk
}

type EntB struct {
    EntAs   []EntA
}
```

## Nombres de tablas y de structs
Las tablas se llamarán como las structs, con ciertas reinterpretaciones astutas que aplica gorm a todos los nombres. Si queremos que los tipos sean cosas como UserEntity para el programa, hay que reescribir dos métodos por cada estructura. Dado ese trabajo extra, llamamos a los tipos del ORM con el nombre "bonito", esto es, "User". De todas formas, para el mundo exterior a repos el tipo se llama repos.User