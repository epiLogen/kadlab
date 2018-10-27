# kadlab

För att starta kademlia kör:
1. docker-compose build
2. docker-compose up --scale kademlianodes=100
3. vänta på att bootstrapingen är avslutad och slutat ge utskrifter
4. öppna upp ett nytt terminal fönster och kör docker ps
5. docker attach *servernodens ID*
6. Gör filer och ge önskat kommando tex: ./dfs store hello

För att köra test och se coverage i browsern:
1. Gå till d7024e mappen
2. kör: go test -vet="off" -coverprofile=coverage.out
3. kör: go tool cover -html=coverage.out
