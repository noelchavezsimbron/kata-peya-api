# kata-peya-api

To start the api:
```shell
docker-compose up
```

To test the connection:
```shell
curl -X GET http://localhost:8080/api/kata-peya/v1/pets
```

You should see
```json
[
    {
        "id": 1,
        "name": "Rino",
        "vaccines": [
            "distemper",
            "parvovirus"
        ],
        "age": "3 years"
    },
    {
        "id": 2,
        "name": "Braco",
        "vaccines": [
            "rabies"
        ],
        "age": "3 months"
    },
    {
        "id": 3,
        "name": "Duke",
        "vaccines": [],
        "age": "1 year"
    },
    {
        "id": 4,
        "name": "Bolt",
        "vaccines": [
            "rabies",
            "distemper",
            "parvovirus"
        ],
        "age": "3 months"
    }
]
```

Perform integration & e2e test:
```shell
make test.e2e
make test.integration
```

Perform unit test:
```shell
make test.unit
```