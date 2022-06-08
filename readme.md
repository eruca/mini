# 说明

- cmd 运行 docker ps 实际上运行的是
```cmd
docker-compose -f nginx-goapp.yml up -d # --build
docker-compose -f nignx-goapp.yml down
```