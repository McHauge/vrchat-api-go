# Scarica il file openapi.yaml
Invoke-WebRequest -Uri "https://vrchatapi.github.io/specification/openapi.yaml" -OutFile "openapi.yaml"

# Installa openapi-codegen
go install github.com/mayocream/openapi-codegen@latest

# Genera il codice utilizzando openapi-codegen
openapi-codegen -i ./openapi.yaml -o . -p vrchat