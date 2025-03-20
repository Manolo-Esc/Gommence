
# Puedes ver exactamente qué partes del código no están cubiertas usando **`go test -coverprofile` y `go tool cover`**.  


### **Paso a paso para ver qué líneas no están cubiertas:**

#### **1️⃣ Generar un informe detallado**
Ejecuta:
```sh
go test -coverprofile=coverage.out
```
Esto genera un archivo `coverage.out` con información sobre qué líneas fueron ejecutadas y cuáles no.

#### **2️⃣ Ver un resumen en consola**
Para ver un desglose simple en la terminal:
```sh
go tool cover -func=coverage.out
```
Ejemplo de salida:
```
github.com/tu-proyecto/main.go:12:  someFunction  80.0%
github.com/tu-proyecto/utils.go:25:  anotherFunc   50.0%
total:                              (statements)  75.6%
```
Esto te muestra **qué funciones tienen menos cobertura**.

#### **3️⃣ Ver qué líneas específicas no están cubiertas**
Para una vista interactiva en HTML:
```sh
go tool cover -html=coverage.out -o coverage.html
```
Luego, abre `coverage.html` en un navegador.  
Verás tu código resaltado:  
- **Rojo** = Código no ejecutado por los tests.  
- **Verde** = Código cubierto por los tests.  

---

### **Conclusión**
✅ **`go test -coverprofile=coverage.out`** → Genera un informe detallado.  
✅ **`go tool cover -func=coverage.out`** → Muestra cobertura por función.  
✅ **`go tool cover -html=coverage.out`** → Muestra qué líneas exactas faltan en un navegador.  

Así puedes ver **exactamente qué partes de tu código no están siendo ejecutadas por los tests** y arreglarlo. 🚀