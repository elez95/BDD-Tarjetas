//por ahora, para su correcto funcionamiento, la base de datos debe estar vacia, solo creada y sin conectarla

package main

//Importo paquetes
import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

//main----------------------------------------------------------------------------------------------------
func main() {

	var opcion_elegida int //numero que elegira el usuario para ejecutar una opcion

	mostrar_opciones()

	fmt.Scan(&opcion_elegida)

	ejecutar_opcion(opcion_elegida)

}

//funcion que se llama desde el main para mostrar todas las opciones del CLI-------------------------------
func mostrar_opciones() {

	fmt.Println("\n-------------------------------\nElija una opcion para ejecutar:")
	fmt.Println("1- Crear base de datos")
	fmt.Println("2- Crear tablas")
	fmt.Println("3- Ingresar datos a las tablas")
	fmt.Println("4- Crear funciones")
	fmt.Println("5- Realizar compras")
	fmt.Println("6- Generar resumen")
	fmt.Println("7- Borrar pks")
	fmt.Println("8- Salir\n")
}

//funcion que detecta la opción elegida a ejecutar---------------------------------------------------------
func ejecutar_opcion(opcion_elegida int) {

	fmt.Printf("La opcion elegida fue %v \n", opcion_elegida) //linea solo de prueba para ver que funcione

	if opcion_elegida == 1 {
		crear_bdd()
		conectar_con_bdd()
		main()

	} else if opcion_elegida == 2 {

		crear_tablas()
		main()

	} else if opcion_elegida == 3 {

		llenar_tablas()
		main()

	} else if opcion_elegida == 4 {

		crear_todas_las_funciones()
		main()

	} else if opcion_elegida == 5 {

		realizar_compras()
		main()

	} else if opcion_elegida == 6 {

		generar_resumen()
		main()
		
	} else if opcion_elegida == 7 {
		
		borrar_pks()
		main()	

	} else if opcion_elegida == 8 {

		fmt.Println("###### Fin ######")

	} else {

		fmt.Println("Error, ingresa nuevamente")
		main()

	}

}

//funcion para crear la base de datos----------------------------------------------------------------------

func crear_bdd() {
	//conectamos con postgres
	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//creamos nuestra base de datos
	_, err = db.Exec(`drop database if exists basedatos`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`create database basedatos`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n### Base de datos creada correctamente ###\n")
}

//funcion para conectar con nuestra bdd --------------------------------------------------------

func conectar_con_bdd() *sql.DB {

	db, err := sql.Open("postgres", "user=postgres host=localhost dbname=basedatos sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	return db
}

//funcion para crear las tablas-----------------------------------------------------------------------------
//Lo primero que hace es llamar a la funcion para conectar con nuestra bdd y lo guarda en la variable db
//después crea las tablas (falta completar las demás, solo hice uno de prueba)
//luego chequea los errores
//por ultimo cierra la conexion con la base (esto debe hacerse en funcion aparte porque debe permanecer abierta

func crear_tablas() {

	db := conectar_con_bdd()
	defer db.Close()

	_, err := db.Exec(`create table cliente(nrocliente  int, nombre text, apellido text, domicilio text, telefono char(12));

create table tarjeta(nrotarjeta char(16), nrocliente int, validadesde char(6), validahasta char(6), codseguridad char(4), limitecompra decimal(8,2), estado char(10));

create table comercio(nrocomercio int, nombre text, domicilio text, codigopostal char(8), telefono char(12));

create table compra(nrooperacion serial, nrotarjeta char(16), nrocomercio int, fecha timestamp, monto decimal(7,2), pagado boolean);

create table rechazo(nrorechazo serial, nrotarjeta char(16), nrocomercio int, fecha timestamp, monto decimal(7,2), motivo text);

create table cierre(año int, mes int, terminacion int, fechainicio date, fechacierre date, fechavto date);

create table cabecera(nroresumen serial, nombre text, apellido text, domicilio text, nrotarjeta char(16), desde date, hasta date, vence date, total decimal(8,2));

create table detalle(nroresumen int, nrolinea int, fecha date, nombrecomercio text, monto decimal(7,2));

create table alerta(nroalerta serial, nrotarjeta char(16), fecha timestamp, nrorechazo int, codalerta int, descripcion text);

create table consumo(nrotarjeta char(16), codseguridad char(4), nrocomercio int, monto decimal(7,2));`)

	if err != nil {
		log.Fatal(err)
	}

	crear_pk_fk()

	fmt.Printf("\n### Tablas creadas ###\n")
}

//Funcion que crea todas las pk's y fk's--------------------------------------------------------------------

func crear_pk_fk() {

	db := conectar_con_bdd()
	defer db.Close()

	_, err := db.Exec(`alter table cliente  add constraint cliente_pk  primary key (nrocliente);
alter table tarjeta  add constraint tarjeta_pk  primary key (nrotarjeta);
alter table comercio add constraint comercio_pk primary key (nrocomercio);
alter table compra   add constraint compra_pk   primary key (nrooperacion);
alter table rechazo  add constraint rechazo_pk  primary key (nrorechazo);
alter table cierre   add constraint cierre_pk   primary key (año, mes, terminacion);
alter table cabecera add constraint cabecera_pk primary key (nroresumen);
alter table detalle  add constraint detalle_pk  primary key (nroresumen, nrolinea);
alter table alerta   add constraint alerta_pk   primary key (nroalerta);

alter table tarjeta  add constraint tarjeta_nrocliente_fk  foreign key (nrocliente)  references cliente(nrocliente);

alter table compra   add constraint compra_nrotarjeta_fk   foreign key (nrotarjeta)  references tarjeta(nrotarjeta);
alter table compra   add constraint compra_nrocomercio_fk  foreign key (nrocomercio) references comercio(nrocomercio);

alter table rechazo  add constraint rechazo_nrotarjeta_fk  foreign key (nrotarjeta)  references tarjeta(nrotarjeta);
alter table rechazo  add constraint rechazo_nrocomercio_fk foreign key (nrocomercio) references comercio(nrocomercio);

alter table cabecera add constraint cabecera_nrotarjeta_fk foreign key (nrotarjeta)  references tarjeta(nrotarjeta);

alter table detalle  add constraint detalle_nroresumen_fk  foreign key (nroresumen)  references cabecera(nroresumen);

alter table alerta   add constraint alerta_nrotarjeta_fk   foreign key (nrotarjeta)  references tarjeta(nrotarjeta);
alter table alerta   add constraint alerta_nrorechazo_fk   foreign key (nrorechazo)  references rechazo(nrorechazo);

alter table consumo  add constraint consumo_nrotarjeta_fk  foreign key (nrotarjeta)  references tarjeta(nrotarjeta);
alter table consumo  add constraint consumo_nrocomercio_fk foreign key (nrocomercio) references comercio(nrocomercio);`)

	if err != nil {
		log.Fatal(err)
	}

}

//Inserts----------------------------------------------------------------------------------------------------
func llenar_tablas() {

	db := conectar_con_bdd()
	defer db.Close()

	_, err := db.Exec(`--clientes 
insert into cliente values(1,'Daniela Anabel','Oviedo','San Martin 3814','541130569988');
insert into cliente values(2,'Fernando','Ferreyra','Benito Lynch 2206','541156441305');
insert into cliente values(3, 'Elias','Goñez', 'Valparaiso 2050','541128898392');
insert into cliente values(4,'Romina','Segretin','Uruguay 790','541154085062');
insert into cliente values(5,'Fabian','García Gómez','Av. Alem 368','541140495127');
insert into cliente values(6,'Matheo Samuel','García','Av. Callao 1311','541146155914');
insert into cliente values(7,'Sabrina Rosalia','Ramirez','Av. Sáenz 945','541130360237');
insert into cliente values(8,'Sara Valeria','Hernández','Av. Cabildo 2523','541128447247');
insert into cliente values(9,'Alicia Grisel','Gómez','Cura Brochero 1053','541161206314');
insert into cliente values(10,'Joana Elizabeth','Villarreal','Palpa 1020','541164294818');
insert into cliente values(11,'Ignacio Ariel','Perez','Paso de los patos 2508','541189768847');
insert into cliente values(12,'Lucia Daniela','Benitez','Av. Rivadavia 2199','541124361554');
insert into cliente values(13,'Maximiliano Ezequiel','Fernandez','Obrien 2460','541167353600');
insert into cliente values(14,'Cristian Elias','Oviedo','Yatasto 1749','541126858087');
insert into cliente values(15,'Carolina Noelia','Diaz','Falucho 853','541151955038');
insert into cliente values(16,'Agustina','Lopez','Ricardo Rojas 1183','541141612153');
insert into cliente values(17,'Luciano Damian','Mansilla','Av. Eva Duarte de Perón 904','541136471202');
insert into cliente values(18,'Hernan Daniel','Rondelli','Nazca 1065','541127146757');
insert into cliente values(19,'Leandro David','Gimenez','Juan Maria Gutiérrez 1150','541125405212');
insert into cliente values(20,'Rodrigo Ezquiel','Palacios','Pablo Areguati 299','541124511771');


--tarjetas
insert into tarjeta values('4286283215095190', 1, '201709', '202208', '114', 45000.00, 'vigente');
insert into tarjeta values('4532449515464319', 2, '202001', '202412', '881', 30000.00, 'vigente');
insert into tarjeta values('4716905901199213', 3, '202108', '202607', '311', 15000.00, 'vigente');
insert into tarjeta values('4539760286740064', 4, '202204', '202703', '553', 35000.00, 'vigente');
insert into tarjeta values('4916197097056062', 5, '202010', '202509', '103', 45000.00, 'anulada');
insert into tarjeta values('4532157860627139', 6, '202004', '202503', '802', 42000.00, 'anulada');
insert into tarjeta values('4449942525596585', 7, '202010', '202509', '552', 12000.00, 'vigente');
insert into tarjeta values('4929028998516745', 8, '201610', '202109', '412', 11000.00, 'suspendida');
insert into tarjeta values('4916558526474988', 9, '201604', '202103', '633', 65000.00, 'vencida');
insert into tarjeta values('4456844734152285', 10, '201707', '202206', '853', 35000.00, 'anulada');
insert into tarjeta values('5305073210930499', 11, '201707', '202206', '271', 14000.00, 'vigente');
insert into tarjeta values('5115874922952014', 12, '202008', '202507', '647', 70000.00, 'suspendida');
insert into tarjeta values('5433516727758253', 13, '201802', '202301', '345', 15000.00, 'vigente');
insert into tarjeta values('5200557813577356', 14, '201707', '202206', '112', 12000.00, 'anulada');
insert into tarjeta values('5425807573408337', 15, '201712', '202211', '879', 43000.00, 'vigente');
insert into tarjeta values('5255982663365344', 16, '201906', '202405', '768', 12000.00, 'suspendida');
insert into tarjeta values('5535292533476491', 17, '201805', '202304', '876', 17000.00, 'vigente');
insert into tarjeta values('5425758312840399', 18, '202005', '202504', '881', 80000.00, 'vigente');
insert into tarjeta values('340869936801114', 17, '201907', '202406', '675', 90000.00, 'vigente'); 
insert into tarjeta values('342888106007110', 18, '202103', '202602', '127', 12000.00, 'vigente');
insert into tarjeta values('343263611209214', 19, '201909', '202408', '901', 20000.00, 'anulada');
insert into tarjeta values('377829618815820', 20, '201804', '202303', '320', 75000.00, 'suspendida');


--comercios
insert into comercio values(1, 'Coto', 'Belgrano 960', 'B1619JHU','034844458867');--comercios con el mismo CP
insert into comercio values(2, 'Sodimac','Constituyentes 1370','B1619JHU','112658423658');--comercios con el mismo CP
insert into comercio values(3, 'Buen Gusto', 'Av. Libertador 3072', 'C1245YTD','541126598965');
insert into comercio values(4, 'Cafeteria Victor', 'Juan Gutierrez 1150', 'B1613GAE', '541178451245');
insert into comercio values(5, 'Libreria Alondra', 'Mateo Churich 130', 'B1619JGB', '541125584518');
insert into comercio values(6, 'Carrefour', 'Los Andes 458', 'B1608OKL', '541126154879');
insert into comercio values(7, 'El Boulevard', 'General Peron 377', 'B1610HGU', '541128964712');
insert into comercio values(8, 'Rapanui', 'Juan Domingo Peron 1974', 'C1456NSM', '541126597841');
insert into comercio values(9, 'Ñoquis Artesanales', 'Balcarce 50', 'C1064KCF', '541143443600');
insert into comercio values(10, 'McDonalds', 'Hipolito Yrigoyen 267', 'B1610LPN','541126455468');
insert into comercio values(11, 'Starbucks', '25 de Mayo 2254', 'B1609KGJ','541148897234');
insert into comercio values(12, 'Parrilla El Chorizon', 'Blandengues 483', 'B1611HGE', '541136223879');
insert into comercio values(13, 'Optica Casimiro', 'Juan B Justo 2020', 'C1032BND', '541121145469');
insert into comercio values(14, 'Ferreteria El Cosito', 'Peru 1654', 'B1663ERF', '541136458159');
insert into comercio values(15, 'Farmacia Favaloro', 'Remedios de Escalada 392', 'B1619HJU', '541125698741');
insert into comercio values(16, 'Servicio Tecnico LG', 'Marie Curie 506', 'B1600KIB', '541125896734');
insert into comercio values(17, 'Loteria de la provincia', 'Cayetano Bourdet 2390', 'B1619PER', '541123698564');
insert into comercio values(18, 'Supermercado Puma', 'Cordoba 212', 'B1610GBF', '541128955864');
insert into comercio values(19, 'Aberturas Pepe', '9 de Julio 3004', 'C1040JUG', '541126897468');
insert into comercio values(20, 'Cinemark', 'Constituyentes 2078', 'B1620MVU', '541128969864');

--consumos

insert into consumo values('4716905901199213', '311', 10, 750.00);
insert into consumo values('5305073210930499', '271', 6, 1500.00);
insert into consumo values('5535292533476491', '876', 1, 3000.00);
insert into consumo values('5425758312840399', '881', 15, 1000.00);
insert into consumo values('4449942525596585', '552', 12, 2000.00);
insert into consumo values('4916197097056062', '103', 11, 500.00);--anulada
insert into consumo values('4449942525596585', '411', 2, 12000.00);--tarjeta mal codigo de seguridad
insert into consumo values('4916558526474988', '633', 4, 3000.00);--tarjeta vencida 
insert into consumo values('4929028998516745', '412', 5, 5000.00);--tarjeta suspendida
insert into consumo values('4286283215095190', '114', 1, 1000.00);
insert into consumo values('4286283215095190', '114', 2, 1000.00);--2 compras en menos de un minuto en comercios distintos mismo CP
insert into consumo values('5425807573408337', '879', 20, 44000.00);--compra supera el limite de la tarjeta
insert into consumo values('5425807573408337', '879', 20, 44000.00);--segunda vez rechazada por exceso del limite`)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n### Inserts creados###\n")
}

//funcion que llama a todas las funciones de sql para que se creen y se guarden en la base de datos (completar)-----------
func crear_todas_las_funciones() {

	crear_funcion_autorizar_compra()
	crear_funcion_realizar_compras()
	crear_verificar_vigencia()
	crear_llenar_cierre()
	crear_generar_resumen()
	crear_func_alerta_rechazo()
	crear_func_alerta_compra()

	fmt.Printf("\n### Funciones guardadas en la base de datos ###\n")
}

//Funcion funcierre que se guarda en la base de datos---------------------------------------------------------------------

func crear_llenar_cierre() {

	db := conectar_con_bdd()
	defer db.Close()

	_, err := db.Exec(`create or replace function llenar_cierre() returns void as $$
declare
	i int :=0;
	j int :=0;
	n int :=9;
	m int :=11;
	fecha_inicio date :='2020-12-28';
	fecha_cierre date :='2021-01-27';
	fecha_vencimiento date :='2021-02-10';
begin
for i in i..n loop
    for j in j..m loop
        insert into cierre values(2021, j+1, i, fecha_inicio, fecha_cierre, fecha_vencimiento);
        if (EXTRACT(ISOYEAR FROM fecha_vencimiento) = 2022) then
            fecha_inicio := fecha_inicio - cast('11 month' as interval);
            fecha_cierre := fecha_cierre - cast('11 month' as interval);
            fecha_vencimiento := fecha_vencimiento - cast('11 month' as interval);
        else
            fecha_inicio := fecha_inicio + cast('1 month' as interval);
            fecha_cierre := fecha_cierre+ cast('1 month' as interval);
            fecha_vencimiento := fecha_vencimiento + cast('1 month' as interval);
        end if;
    end loop;
end loop;
end;
$$ language plpgsql;`)

	if err != nil {
		log.Fatal(err)
	}

}

//Creo la funcion autorizar_compra que se va a guardar en la base de datos--------------------------------------------------------------------
func crear_funcion_autorizar_compra() {

	db := conectar_con_bdd()
	defer db.Close()

	_, err := db.Exec(`create or replace function autorizar_compra(nro_tarjeta char(16), cod_seguridad char(4), nro_comercio int, p_monto decimal(8,2)) returns boolean as $$
declare
    fecha_actual timestamp := current_timestamp(0);
    tarjeta record;
    monto_total decimal:= p_monto;
begin
    --verifico que exista alguna compra realizada por la tarjeta pasada como parametro
    if ((select count(*) from compra where nrotarjeta = nro_tarjeta ) > 0) then 
        --sumo el total de las compras realizas por esa tarjeta mas la nueva compra
        monto_total := monto_total + (select sum(monto) from compra where nrotarjeta = nro_tarjeta and pagado = false); 
    end if;
    
    select * into tarjeta from tarjeta where nrotarjeta = nro_tarjeta;
    if  not found then 
        insert into rechazo (nrotarjeta, nrocomercio, fecha, monto, motivo) 
        values(nro_tarjeta, nro_comercio, fecha_actual, p_monto, 'tarjeta no valida o no vigente');
        return false;
    
    elsif cod_seguridad != tarjeta.codseguridad then
        insert into rechazo (nrotarjeta, nrocomercio, fecha, monto, motivo) 
        values(nro_tarjeta, nro_comercio, fecha_actual, p_monto, 'codigo de seguridad invalido');
        return false;
    
    elsif (monto_total > tarjeta.limitecompra) then
        insert into rechazo (nrotarjeta, nrocomercio, fecha, monto, motivo) 
        values(nro_tarjeta, nro_comercio, fecha_actual, p_monto, 'supera limite de tarjeta');
        return false;
    
    elsif (select verificar_vigencia((tarjeta.validahasta))) then
        insert into rechazo (nrotarjeta, nrocomercio, fecha, monto, motivo) 
        values(nro_tarjeta, nro_comercio, fecha_actual, p_monto, 'plazo de vigencia expirado');
        return false;

    elsif 'suspendida' = (tarjeta.estado) then
        insert into rechazo (nrotarjeta, nrocomercio, fecha, monto, motivo) 
        values(nro_tarjeta, nro_comercio, fecha_actual, p_monto, 'la tarjeta se encuentra suspendida');
        return false;

    elsif 'anulada' = (tarjeta.estado) then
        insert into rechazo (nrotarjeta, nrocomercio, fecha, monto, motivo) 
        values(nro_tarjeta, nro_comercio, fecha_actual, p_monto, 'la tarjeta se encuentra anulada');
        return false;

    else
        insert into compra (nrotarjeta, nrocomercio, fecha, monto, pagado) 
        values(nro_tarjeta, nro_comercio, fecha_actual, p_monto, false);--se autoriza la compra
        return true;
    end if;
end;
$$ language plpgsql;`)

	if err != nil {
		log.Fatal(err)
	}

}

//Funcion func_alerta_rechazo que se guarda en la base de datos-------------------------------------------------------

func crear_func_alerta_rechazo() {

	db := conectar_con_bdd()
	defer db.Close()

	_, err := db.Exec(`create or replace function func_alerta_rechazo() returns trigger as $$
declare
    undia interval := '24:00:00';
    i record;
begin
    insert into alerta (nrotarjeta,fecha ,nrorechazo, codalerta, descripcion) 
    values(new.nrotarjeta, new.fecha, new.nrorechazo, 0, 'se produjo un rechazo');

    for i in select * from rechazo where nrotarjeta = new.nrotarjeta and motivo = 'supera limite de tarjeta' loop 
        if (new.fecha - i.fecha) < undia then

            update tarjeta set estado = 'suspendida' where nrotarjeta = new.nrotarjeta;
            
            insert into alerta (nrotarjeta,fecha ,nrorechazo, codalerta, descripcion) 
            values(new.nrotarjeta, new.fecha, new.nrorechazo, 32, 'supero el limite de compra mas una vez');

        end if; 
    end loop;   
    return new;
end;
$$ language plpgsql;

create trigger rechazo_trg
after insert on rechazo
for each row
execute procedure func_alerta_rechazo();`)

	if err != nil {
		log.Fatal(err)
	}

}

//Funcion func_alerta_compra que se guarda en la base de datos--------------------------------------------------------------

func crear_func_alerta_compra() {

	db := conectar_con_bdd()
	defer db.Close()

	_, err := db.Exec(`create function func_alerta_compra() returns trigger as $$
declare
    unminuto interval := '00:01:00';
    cincominutos interval := '00:05:00';

    i record;
    j record;
begin
    if (select count(*) from compra where nrotarjeta = new.nrotarjeta) > 1 then
            
        for i in select * from compra where nrotarjeta = new.nrotarjeta and nrocomercio in
            (select nrocomercio from comercio where nrocomercio != new.nrocomercio and codigopostal = 
             (select codigopostal from comercio where nrocomercio = new.nrocomercio)) loop

            if (new.fecha - i.fecha) <= unminuto then
            
                insert into alerta (nrotarjeta,fecha ,nrorechazo, codalerta, descripcion) 
                values(new.nrotarjeta, new.fecha, null, 1 ,'dos compras dentro del distrito en menos de un minuto'); 
         
            end if;
        end loop;
               
        for j in select fecha from compra where nrotarjeta = new.nrotarjeta and nrocomercio in
            (select nrocomercio from comercio where codigopostal != 
             (select codigopostal from comercio where nrocomercio = new.nrocomercio)) loop

            if (new.fecha - j.fecha) <= cincominutos then

                insert into alerta (nrotarjeta,fecha ,nrorechazo, codalerta, descripcion) 
                values(new.nrotarjeta, new.fecha, null, 5 ,'dos compras fuera del distrito en menos de 5 minutos');
            
            end if;
        end loop;
    end if;
    return new;
end;
$$ language plpgsql;

create trigger compra_trg
after insert on compra
for each row
execute procedure func_alerta_compra();`)

	if err != nil {
		log.Fatal(err)
	}

}

//Funcion para que recorre la tabla consumo y autoriza la compra -----------------------------------------------------------
//Esta funcion se guarda en la base de datos
func crear_funcion_realizar_compras() {

	db := conectar_con_bdd()
	defer db.Close()

	_, err := db.Exec(`create or replace function realizar_compras() returns void as $$
declare
	fila record;
begin
	for fila in select * from consumo loop
		perform autorizar_compra(fila.nrotarjeta, fila.codseguridad, fila.nrocomercio, fila.monto);
	end loop;	
	return;
end;
$$ language plpgsql;`)

	if err != nil {
		log.Fatal(err)
	}

}

//Funcion genera_resumen que se guarda en la base de datos-----------------------------------------------------------------

func crear_generar_resumen() {

	db := conectar_con_bdd()
	defer db.Close()

	_, err := db.Exec(`create or replace function generar_resumen(nro_cliente int, periodo_año int, periodo_mes int) returns void as $$
declare
    dato_cliente record;
    tarjeta_cliente record;
    dato_cierre record;
    fila_compras record;

    total_a_pagar decimal(8,2);
    contador_linea int := 1;   
begin
    select * into dato_cliente from cliente where nrocliente = nro_cliente;
    --para cada tarjeta del cliente hacemos...
    for tarjeta_cliente in select * from tarjeta where nrocliente = nro_cliente loop--obtener la terminacion de esa tarjeta 
                                                                                       
        select * into dato_cierre from cierre where terminacion = (cast(substr(tarjeta_cliente.nrotarjeta, length(tarjeta_cliente.nrotarjeta)) as int))
        and mes = periodo_mes;--obtener los datos de cierre para esa tarjeta de le cliente y para ese periodo  

        --sumamos el total de compras para esa tarjeta y ese periodo
        total_a_pagar:= (select sum(monto) from compra where nrotarjeta = tarjeta_cliente.nrotarjeta and (extract(month from fecha)) = periodo_mes and pagado = false); 

        insert into cabecera (nombre, apellido, domicilio, nrotarjeta, desde, hasta, vence, total)
        values(dato_cliente.nombre, dato_cliente.apellido, dato_cliente.domicilio, tarjeta_cliente.nrotarjeta, 
               dato_cierre.fechainicio, dato_cierre.fechacierre, dato_cierre.fechavto, total_a_pagar);
                       
        for fila_compras in select * from compra where nrotarjeta = tarjeta_cliente.nrotarjeta and (extract(month from fecha)) = periodo_mes and pagado = false loop
            
            insert into detalle values((select nroresumen from cabecera where nrotarjeta = tarjeta_cliente.nrotarjeta), contador_linea, fila_compras.fecha, 
                                        (select nombre from comercio where nrocomercio = fila_compras.nrocomercio), fila_compras.monto);
            
            contador_linea := contador_linea + 1;

        end loop;
        contador_linea := 1;
    end loop;
end;
$$ language plpgsql;`)

	if err != nil {
		log.Fatal(err)
	}

}

//Funcion verificar_vigencia que se guarda en la base de datos-------------------------------------------------------------

func crear_verificar_vigencia() {

	db := conectar_con_bdd()
	defer db.Close()

	_, err := db.Exec(`create or replace function verificar_vigencia(fecha_vencimiento char(6)) returns boolean as $$
declare
     fecha_actual date :=to_date(to_char(current_date,'YYYYMM'),'YYYYMM'); --extrae el año y mes de la fecha actual en formato date
     fecha_tarjeta date:=to_date(fecha_vencimiento, 'YYYYMM'); --extrae el año y mes de la fecha de vencimiento de la tarjeta en formato date
begin
     if (fecha_tarjeta <= fecha_actual) then --si la fecha es menor a la fecha actual esta vencida.
        return true;
     end if;
     return false;
end;
$$ language plpgsql;


--funcion que recorre la tabla consumo y va autorizando cada fila
create or replace function realizar_compras() returns void as $$
declare
	fila record;
begin
	for fila in select * from consumo loop
		perform autorizar_compra(fila.nrotarjeta, fila.codseguridad, fila.nrocomercio, fila.monto);
	end loop;	
	return;
end;
$$ language plpgsql;`)

	if err != nil {
		log.Fatal(err)
	}

}

//Funcion que llama a la funcion de realizar compras que esta en la base de datos---------------------------
func realizar_compras() {

	db := conectar_con_bdd()
	defer db.Close()

	_, err := db.Exec(`select realizar_compras()`)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n### Compras realizadas ###\n")

	//Hago el cierre
	_, err1 := db.Exec(`select llenar_cierre()`)
	if err != nil {
		log.Fatal(err1)
	}

}

//funcion para llamar a generar resumen---------------------------------------------------------------------

func generar_resumen() {
	
	db := conectar_con_bdd()
	defer db.Close()

	_, err := db.Exec(`select generar_resumen(1,2010,6)`)

	if err != nil {
		log.Fatal(err)
	}

}


//funcion para borrar pks-----------------------------------------------------------------------

func borrar_pks(){
	db := conectar_con_bdd()
	defer db.Close()
	
	_, err := db.Exec(`--borrado fk
alter table tarjeta drop constraint tarjeta_nrocliente_fk;
alter table compra drop constraint compra_nrotarjeta_fk;
alter table compra drop constraint compra_nrocomercio_fk;
alter table rechazo drop constraint rechazo_nrotarjeta_fk;
alter table rechazo drop constraint rechazo_nrocomercio_fk;
alter table cabecera drop constraint cabecera_nrotarjeta_fk;
alter table detalle drop constraint detalle_nroresumen_fk;
alter table alerta drop constraint alerta_nrotarjeta_fk;
alter table alerta drop constraint alerta_nrorechazo_fk;
alter table consumo drop constraint consumo_nrotarjeta_fk;
alter table consumo drop constraint consumo_nrocomercio_fk;
--borrado pk
alter table cliente drop constraint cliente_pk;
alter table tarjeta drop constraint tarjeta_pk;
alter table comercio drop constraint comercio_pk;
alter table compra drop constraint compra_pk;
alter table rechazo drop constraint rechazo_pk;
alter table cierre drop constraint cierre_pk;
alter table cabecera drop constraint cabecera_pk;
alter table detalle drop constraint detalle_pk;
alter table alerta drop constraint alerta_pk;`)

	if err != nil {
		log.Fatal(err)
	}
	
}
