### Generate Certs
```bash
docker-compose exec -it gvmd gvm-manage-certs -af
```

### Validate the Certs
```bash
docker-compose exec -it gvmd gvm-manage-certs -V
```

it should be needed only in the first time, as the CA is mounted locally

