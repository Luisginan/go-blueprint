docker pull mysql:latest
docker run -d --name mysql-container -p 3306:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql:latest


