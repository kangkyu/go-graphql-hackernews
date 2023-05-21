# howtographql

Follow a tutorial https://www.howtographql.com/graphql-go/0-introduction/

```sh
docker run -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=dbpass -e MYSQL_DATABASE=hackernews -d mysql:latest

go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2

migrate -database mysql://root:dbpass@/hackernews -path internal/pkg/db/migrations/mysql up
```
