## ILWrapper API Documentation

Diese README beschreibt die verfügbaren API-Endpunkte der ILWrapper-API.

```
# Base URL
http://blha-dimagapps-ilwrapper
```

### Authentifizierung
Die API verwendet Bearer Token Authentifizierung. Nach dem Login erhältst du ein JWT-Token, das für nachfolgende 
Requests im Authorization Header verwendet werden muss.

### Endpoints

#### Login

Authentifizierung eines Benutzers und Erhalt eines JWT-Tokens.

``` json 
Endpoint: POST /api/login
Request Body:
{
    "email": "user",
    "password": "password"
}
```

``` bash
# cURL Beispiel:
curl -X POST http://blha-dimagapps-ilwrapper/api/login \
    -H "Content-Type: application/json" \
    -d '{"email":"user","password":"password"}'
```

#### Remote-Validate

Validierung einer bereits auf dem Server vorhandenen Datei. 
``` json 
Endpoint: POST /api/create
Headers:
    Authorization: Bearer <YOUR_JWT_TOKEN>
    Content-Type: application/json
Request Body:
    {
        "filename": "data/dcf71987-d0e6-4522-af64-db163424ca71_20221116164349 Lumped model for ferroelectric transducers.pdf",
        "type": "Validate"
    }
```

``` bash 
# cURL Beispiel:
curl -X POST http://blha-dimagapps-ilwrapper/api/create \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer YOUR_JWT_TOKEN" \
    -d '{"filename":"data/example.pdf","type":"Validate"}'
```

#### Remote-Identify

Identifizierung des Dateiformats einer bereits auf dem Server vorhandenen Datei.
``` json 
Endpoint: POST /api/create
Headers:
    Authorization: Bearer <YOUR_JWT_TOKEN>
    Content-Type: application/json
Request Body:
    {
        "filename": "data/dcf71987-d0e6-4522-af64-db163424ca71_20221116164349 Lumped model for ferroelectric transducers.pdf",
        "type": "Identify"
    }
```

``` bash
# cURL Beispiel:
curl -X POST http://blha-dimagapps-ilwrapper/api/create \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer YOUR_JWT_TOKEN" \
    -d '{"filename":"data/example.pdf","type":"Identify"}'
```

#### Local-Validate
Validierung einer hochgeladenen Datei.
``` json 
Endpoint: POST /api/create
Headers:
    Authorization: Bearer <YOUR_JWT_TOKEN>
    Content-Type: multipart/form-data
Form Data:
    type: Validate
    file: (binary file)
```
``` bash
# cURL Beispiel:
curl -X POST http://blha-dimagapps-ilwrapper/api/create \
-H "Authorization: Bearer YOUR_JWT_TOKEN" \
-F "type=Validate" \
-F "file=@/path/to/your/file.pdf"
```

#### Local-Identify
Identifizierung des Dateiformats einer hochgeladenen Datei.
``` json 
Endpoint: POST /api/create
Headers:
    Authorization: Bearer <YOUR_JWT_TOKEN>
    Content-Type: multipart/form-data
Form Data:
    type: Identify
    file: (binary file)
```

``` bash
# cURL Beispiel:
curl -X POST http://blha-dimagapps-ilwrapper/api/create \
    -H "Authorization: Bearer YOUR_JWT_TOKEN" \
    -F "type=Identify" \
    -F "file=@/path/to/your/file.html"
``` 
#### Jobs
Abrufen aller Jobs.
``` json ‚
Endpoint: GET /api/jobs
Headers:
    Authorization: Bearer <YOUR_JWT_TOKEN>
```
``` bash
# cURL Beispiel:
curl -X GET http://blha-dimagapps-ilwrapper/api/jobs \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

#### Job by ID
Abrufen eines spezifischen Jobs anhand seiner ID.
``` json
Endpoint: GET /api/job/{id}
Headers:
    Authorization: Bearer <YOUR_JWT_TOKEN>
Path Parameter:
    id (integer): Job-ID
```
``` bash
# cURL Beispiel:
curl -X GET http://blha-dimagapps-ilwrapper/api/job/5 \
     -H "Authorization: Bearer YOUR_JWT_TOKEN"
```
### Operation Types

Die API unterstützt folgende Operationstypen:

- Validate: Validierung einer Datei
- Identify: Identifizierung des Dateiformats

## Workflows

### Typischer Workflow mit Remote-Dateien

# 1. Login
``` bash
TOKEN=$(curl -X POST http://blha-dimagapps-ilwrapper/api/login \
             -H "Content-Type: application/json" \
             -d '{"email":"user","password":"password"}' | jq -r '.token')
```

# 2. Datei validieren (Remote)
``` bash
curl -X POST http://blha-dimagapps-ilwrapper/api/create \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $TOKEN" \
     -d '{"filename":"data/example.pdf","type":"Validate"}'
```
# 3. Jobs abrufen
``` bash
curl -X GET http://blha-dimagapps-ilwrapper/api/jobs \
     -H "Authorization: Bearer $TOKEN"
```

### Typischer Workflow mit File-Upload

# 1. Login
``` bash
TOKEN=$(curl -X POST http://blha-dimagapps-ilwrapper/api/login \
             -H "Content-Type: application/json" \
             -d '{"email":"user","password":"password"}' | jq -r '.token')
```

# 2. Datei hochladen und validieren
``` bash
curl -X POST http://blha-dimagapps-ilwrapper/api/create \
     -H "Authorization: Bearer $TOKEN" \
     -F "type=Validate" \
     -F "file=@/path/to/file.pdf"
```

# 3. Spezifischen Job abrufen
``` bash
curl -X GET http://blha-dimagapps-ilwrapper/api/job/5 \
     -H "Authorization: Bearer $TOKEN"
```

### Unterschiede zwischen Remote und Local

#### Remote-Operationen (JSON):

- Datei muss bereits auf dem Server vorhanden sein
- Verwendung von Content-Type: application/json
- Dateipfad wird im filename Parameter übergeben

#### Local-Operationen (Multipart):

- Datei wird mit dem Request hochgeladen
- Verwendung von Content-Type: multipart/form-data
- Datei wird als file Form-Field übergeben


## Notizen

JWT-Tokens haben ein Ablaufdatum (exp claim im Token), nach Ablauf muss ein neuer Login durchgeführt werden.
Beide Operationstypen (Validate und Identify) verwenden denselben /api/create Endpoint und der Operationstyp wird 
über den type Parameter gesteuert.