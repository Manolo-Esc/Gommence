package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/Manolo-Esc/gommence/src/internal/server"
	"github.com/Manolo-Esc/gommence/src/tests/libtest"
)

func __TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background()) // Crear contexto con cancelación
	defer cancel()                                          // Asegurar que se cancela al final del test

	go func() {
		var args []string
		if err := server.Run(ctx, args, nil, os.Stdin, os.Stdout, os.Stderr); err != nil {
			t.Errorf("Error al ejecutar el servicio: %v", err)
		}
	}()

	// Esperar a que el servicio inicie (puedes mejorar esto con health checks)
	time.Sleep(500 * time.Millisecond)

	// Aquí irían tus tests, por ejemplo, haciendo requests HTTP al servicio
	t.Log("TestRun")

	cancel()                           // Detener el servicio
	time.Sleep(500 * time.Millisecond) // Darle tiempo para cerrarse correctamente
}

func TestRun2(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second) // Timeout para el test
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1) // Añadimos un contador al WaitGroup, indicando que esperamos un proceso

	go func() {
		defer wg.Done() // Al finalizar el servicio, decrementamos el contador
		var args []string
		if err := server.Run(ctx, args, nil, os.Stdin, os.Stdout, os.Stderr); err != nil {
			t.Errorf("Error al ejecutar el servicio: %v", err)
		}
	}()

	// Esperar a que el servidor esté online  XXX no tener hardocoded tantas cosas de la ruta!
	if err := libtest.WaitForReady(ctx, 5, "http://localhost:5080/health"); err != nil {
		t.Fatal("Server is not ready")
	}

	t.Log("Server is ready")
	// Ejecutar pruebas aquí

	cancel() // Paramos la funcion Run()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done) // Al terminar, cerramos el canal
	}()

	// Esperamos hasta que el WaitGroup sea decrementado o se acabe el timeout
	select {
	case <-done:
		t.Log("Server gracefully stopped")
		// El servicio ha terminado correctamente
	//case <-ctx.Done():
	case <-time.After(10 * time.Second):
		t.Fatal("El servicio no terminó a tiempo")
	}
}

func PostWithHeaders() {
	// Crear una nueva solicitud POST
	req, err := http.NewRequest("POST", "https://jsonplaceholder.typicode.com/posts", bytes.NewBuffer([]byte(`{"title": "foo", "body": "bar", "userId": 1}`)))
	if err != nil {
		fmt.Println("Error al crear la solicitud:", err)
		return
	}

	// Establecer encabezados
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer miToken")

	// Hacer la solicitud
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error al hacer la solicitud:", err)
		return
	}
	defer resp.Body.Close()

	// Leer el cuerpo de la respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error leyendo el cuerpo:", err)
		return
	}

	// Imprimir la respuesta
	fmt.Println("Respuesta HTTP Status:", resp.Status)
	fmt.Println("Cuerpo de la respuesta:", string(body))
}

func SimpleGet() {
	// Hacer una solicitud GET
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts/1")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close() // Asegúrate de cerrar el cuerpo de la respuesta

	// Leer el cuerpo de la respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error leyendo el cuerpo:", err)
		return
	}

	// Imprimir el contenido de la respuesta
	fmt.Println("Respuesta HTTP Status:", resp.Status)
	fmt.Println("Cuerpo de la respuesta:", string(body))
}

func SimplePost() {
	// Crear datos que vamos a enviar
	data := map[string]interface{}{
		"title":  "foo",
		"body":   "bar",
		"userId": 1,
	}

	// Convertir los datos a JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error al marshallizar datos:", err)
		return
	}

	// Hacer la solicitud POST
	resp, err := http.Post("https://jsonplaceholder.typicode.com/posts", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error al hacer la solicitud POST:", err)
		return
	}
	defer resp.Body.Close()

	// Leer el cuerpo de la respuesta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error leyendo el cuerpo:", err)
		return
	}

	// Imprimir la respuesta
	fmt.Println("Respuesta HTTP Status:", resp.Status)
	fmt.Println("Cuerpo de la respuesta:", string(body))
}
