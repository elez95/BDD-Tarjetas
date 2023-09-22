package main

import (
	"encoding/json"
	"fmt"
	bolt "github.com/coreos/bbolt"
	"log"
	"strconv"
)

type Cliente struct {
	Nrocliente int
	Nombre     string
	Apellido   string
	Domicilio  string
	Telefono   string
}

type Tarjeta struct {
	Nrotarjeta   string
	Nrocliente   int
	Validadesde  string
	Validahasta  string
	Codseguridad string
	Limitecompra float64
	Estado       string
}

type Comercio struct {
	Nrocomercio  int
	Nombre       string
	Domicilio    string
	Codigopostal string
	Telefono     string
}

type Compra struct {
	Nrooperacion int
	Nrotarjeta   string
	Nrocomercio  int
	Fecha        string
	Monto        float64
	Pagado       bool
}

func CreateUpdate(dbb *bolt.DB, bucketName string, key []byte, val []byte) error {

	tx, err := dbb.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	b, _ := tx.CreateBucketIfNotExists([]byte(bucketName))

	err = b.Put(key, val)
	if err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func ReadUnique(dbb *bolt.DB, bucketName string, key []byte) ([]byte, error) {
	var buf []byte

	err := dbb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		buf = b.Get(key)
		return nil
	})
	return buf, err
}

func cargar_datos() {

	dbb, err := bolt.Open("dbnosql.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer dbb.Close()

	//Hago todos las variables structs

	cliente1 := Cliente{1, "Jose Maria", "Perez", "Av. T de Alvear 1299", "541126598745"}
	cliente2 := Cliente{2, "Roberto", "Rafaela", "Azcuenaga 548", "541146598787"}
	cliente3 := Cliente{3, "Cecilia", "Suarez", "Salta 1210", "541126498789"}

	comercio1 := Comercio{4, "Lo de Tito", "Vivaldi 339", "C1456NSM", "541178955412"}
	comercio2 := Comercio{5, "Moncho", "Catamarca 138", "B1600KIB", "541185749688"}
	comercio3 := Comercio{6, "Tecnico el Chapu", "Canal Beagle 1708", "B1610OIB", "541165754648"}

	tarjeta1 := Tarjeta{"4286283215095190", 1, "201709", "202208", "114", 45000.00, "vigente"}
	tarjeta2 := Tarjeta{"4532449515464319", 2, "202001", "202412", "881", 30000.00, "vigente"}
	tarjeta3 := Tarjeta{"4716905901199213", 3, "202108", "202607", "311", 15000.00, "vigente"}

	compra1 := Compra{7, "4286283215095190", 1, "2021-06-12", 293, true}
	compra2 := Compra{8, "4532449515464319", 2, "2021-06-11", 1800, true}
	compra3 := Compra{9, "4716905901199213", 3, "2021-06-13", 5500, true}

	//Paso todo a JSON

	datacl1, err := json.Marshal(cliente1)
	if err != nil {
		log.Fatal(err)
	}

	datacl2, err := json.Marshal(cliente2)
	if err != nil {
		log.Fatal(err)
	}

	datacl3, err := json.Marshal(cliente3)
	if err != nil {
		log.Fatal(err)
	}

	dataco1, err := json.Marshal(comercio1)
	if err != nil {
		log.Fatal(err)
	}

	dataco2, err := json.Marshal(comercio2)
	if err != nil {
		log.Fatal(err)
	}

	dataco3, err := json.Marshal(comercio3)
	if err != nil {
		log.Fatal(err)
	}

	datata1, err := json.Marshal(tarjeta1)
	if err != nil {
		log.Fatal(err)
	}

	datata2, err := json.Marshal(tarjeta2)
	if err != nil {
		log.Fatal(err)
	}

	datata3, err := json.Marshal(tarjeta3)
	if err != nil {
		log.Fatal(err)
	}

	datacpra1, err := json.Marshal(compra1)
	if err != nil {
		log.Fatal(err)
	}

	datacpra2, err := json.Marshal(compra2)
	if err != nil {
		log.Fatal(err)
	}

	datacpra3, err := json.Marshal(compra3)
	if err != nil {
		log.Fatal(err)
	}

	//Creo los buckets

	CreateUpdate(dbb, "cliente1", []byte(strconv.Itoa(cliente1.Nrocliente)), datacl1)
	CreateUpdate(dbb, "cliente2", []byte(strconv.Itoa(cliente2.Nrocliente)), datacl2)
	CreateUpdate(dbb, "cliente3", []byte(strconv.Itoa(cliente3.Nrocliente)), datacl3)

	CreateUpdate(dbb, "comercio1", []byte(strconv.Itoa(comercio1.Nrocomercio)), dataco1)
	CreateUpdate(dbb, "comercio2", []byte(strconv.Itoa(comercio2.Nrocomercio)), dataco2)
	CreateUpdate(dbb, "comercio3", []byte(strconv.Itoa(comercio3.Nrocomercio)), dataco3)

	CreateUpdate(dbb, "tarjeta1", []byte(tarjeta1.Nrotarjeta), datata1)
	CreateUpdate(dbb, "tarjeta2", []byte(tarjeta2.Nrotarjeta), datata2)
	CreateUpdate(dbb, "tarjeta3", []byte(tarjeta3.Nrotarjeta), datata3)

	CreateUpdate(dbb, "compra1", []byte(strconv.Itoa(compra1.Nrooperacion)), datacpra1)
	CreateUpdate(dbb, "compra2", []byte(strconv.Itoa(compra2.Nrooperacion)), datacpra2)
	CreateUpdate(dbb, "compra3", []byte(strconv.Itoa(compra3.Nrooperacion)), datacpra3)

	//Leo los buckets

	resultado1, err := ReadUnique(dbb, "cliente1", []byte(strconv.Itoa(cliente1.Nrocliente)))
	fmt.Printf("%s\n", resultado1)
	resultado2, err := ReadUnique(dbb, "cliente2", []byte(strconv.Itoa(cliente2.Nrocliente)))
	fmt.Printf("%s\n", resultado2)
	resultado3, err := ReadUnique(dbb, "cliente3", []byte(strconv.Itoa(cliente3.Nrocliente)))
	fmt.Printf("%s\n", resultado3)

	resultado4, err := ReadUnique(dbb, "comercio1", []byte(strconv.Itoa(comercio1.Nrocomercio)))
	fmt.Printf("%s\n", resultado4)
	resultado5, err := ReadUnique(dbb, "comercio2", []byte(strconv.Itoa(comercio2.Nrocomercio)))
	fmt.Printf("%s\n", resultado5)
	resultado6, err := ReadUnique(dbb, "comercio3", []byte(strconv.Itoa(comercio3.Nrocomercio)))
	fmt.Printf("%s\n", resultado6)

	resultado7, err := ReadUnique(dbb, "tarjeta1", []byte(tarjeta1.Nrotarjeta))
	fmt.Printf("%s\n", resultado7)
	resultado8, err := ReadUnique(dbb, "tarjeta2", []byte(tarjeta2.Nrotarjeta))
	fmt.Printf("%s\n", resultado8)
	resultado9, err := ReadUnique(dbb, "tarjeta3", []byte(tarjeta3.Nrotarjeta))
	fmt.Printf("%s\n", resultado9)
	
	resultado10, err := ReadUnique(dbb, "compra1", []byte(strconv.Itoa(compra1.Nrooperacion)))
	fmt.Printf("%s\n", resultado10)
	resultado11, err := ReadUnique(dbb, "compra2", []byte(strconv.Itoa(compra2.Nrooperacion)))
	fmt.Printf("%s\n", resultado11)
	resultado12, err := ReadUnique(dbb, "compra3", []byte(strconv.Itoa(compra3.Nrooperacion)))
	fmt.Printf("%s\n", resultado12)

}

func main() {

	var tecla_uno int

	fmt.Println("\n-------------------------------\nPresiona el num 1 para cargar y leer los datos")

	fmt.Scan(&tecla_uno)

	if tecla_uno == 1 {
		cargar_datos()

	} else {
		fmt.Println("Error, ingresa nuevamente")
		main()
	}
}
