# Santiago Henseler 110732
Trabajo práctico N°0: Nivelador
----
El objetivo del trabajo práctico fue volver a utilizar todas las herramientas ya aprendidas en otras materias vinculándolas en un solo pequeño proyecto.

Protocolo
---- 
Para poder resolver la comunicación entre el servidor y el cliente se implemento un sencillo protocolo con 3 mensajes. Tanto el cliente como el servidor pueden recibir y enviar los 3 tipos de mensajes. Se adopto como convención del protocolo utilizar utf-8 para encodear los strings y bigEndian para los enteros (solo se utilizo enteros de 32 bits).

- Para comunicar las apuestas una por una se envía, por cada una, el mensaje `BET_MESSAGE(0x01)` el cual consta de los siguientes campos:

| BET_MESSAGE | Agency_id | Document | Number | First_name_len | Last_name_len | Birthdate_len | First_name | Last_name | Birthdate
|----|----|----|----|----|----|----|----|----|----|
| 1 byte|  4 bytes | 4 bytes | 4 bytes | 4 bytes | 4 bytes | 4 bytes | len(First_name) | len(Last_name) | len(Birthdate)

Se utilizo este formato de mensaje ya que los campos First_name, Last_name y Birthdate pueden ser variables.
Luego de enviar todos las apuestas se debe enviar el mensaje `END_MESSAGE(0x00)`

- Para comunicar las apuestas en chunks de apuestas se envia el mensaje `BATCH_MESSAGE(0x02)` el cual consta de los siguientes campos:

| BATCH_MESSAGE | len | size | apuestas
|----|----|----|---|
| 1 byte| 4 bytes | 4 bytes | size

Siendo len la cantidad de apuestas que hay en el chunk y size el tamaño total que ocupan. Dentro de apuestas, cada apuesta se ordena de igual manera que al enviar las apuestas individuales:
 Agency_id | Document | Number | First_name_len | Last_name_len | Birthdate_len | First_name | Last_name | Birthdate
|----|----|----|----|----|----|----|----|----|
|  4 bytes | 4 bytes | 4 bytes | 4 bytes | 4 bytes | 4 bytes | len(First_name) | len(Last_name) | len(Birthdate)

Luego de enviar todos los chunks se debe enviar el mensaje `END_MESSAGE(0x00)`

Sincronización
---- 
Para poder esperar a un mínimo de agencias para realizar el sorteo respetando el `AGENCY_QUORUM_MIN` se utilizo una CondVar:
``` Python
  with self.condVar:
    self.bets += 1
    if self.bets == self.quorum:
      self.condVar.notify_all()
    
    while self.bets < self.quorum:
      self.condVar.wait()
```
Con la cual frena a todos los procesos hasta que haya una cantidad de `self.quorum` agencias que ya hayan cargado todas sus apuestas. Y luego de superar ese limite el resto de procesos van a pasar directamente sin tener que esperar. Se podría haber utilizado la herramienta de sincronización Barrier pero al reiniciarse cada vez que se alcanza `self.quorum` bloquearia a threads que no deberian bloquearse.
