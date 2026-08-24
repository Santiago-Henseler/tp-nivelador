echo "services:
  server:
    build:
      context: ./services/server
      dockerfile: Dockerfile
    container_name: server
    volumes:
      - ./output:/output
    ports:
      - "5678:5678"
    environment:
      - PYTHONUNBUFFERED=1
      - AGENCY_QUORUM_MIN=2
      - SERVER_HOST=server
      - SERVER_PORT=5678" > docker-compose.yml


for i in $(seq 0 $1); do
    echo "  client_"$i":
    build:
      context: ./services/client
      dockerfile: Dockerfile
    container_name: client_"$i"
    depends_on:
      - server
    volumes:
      - ./input:/input
    environment:
      - AGENCY_ID="$i"
      - SERVER_HOST=server
      - SERVER_PORT=5678
      - BATCH_SIZE=3
      - INPUT_FILE=input/input-0.csv
      - OUTPUT_FILE=output/winers.txt" >> docker-compose.yaml
done

echo "[INFO] Se creo el docker-compose con" $1 "clientes"