# Kumparan Microservices Generator 

u can make microservices start from proto. see the proto example in `pb/` folder

## Usage
u need to create proto file with path like this
`pb/'$service/$protoname.proto`

example
`pb/content/content_service.proto`

and create microservices
`trafo init --name content-service`
 
 - u will be asked to insert proto source path
 
 - new service will generated like this
 ```
-content-service/
    -client/    ->(Generated From Proto)
    -config/
    -console/
    -db/
    -event/
    -pb/
    -repository/
    -service/     ->(Generated From Proto)
    -worker/
    -config.yml
    -config.yml.dev
    -config.yml.example
    -config.yml.prod
    -config.yml.staging
    -go.mod       ->(mod is already with service name)
    -go.sum
    -LICENSE
    -main.go
    -Makefile
    -README.mod
 ```

## INSTALL
Clone the repository to your desired destination folder e.g :
```
cd ~ && git clone git@github.com:kumparan/trafo.git
```
Append the installation path to global PATH variable
```
echo "export PATH=\$PATH:$HOME/trafo/bin/" >> ~/.bash_profile
```
Source it
```
source ~/.bash_profile
```
It should be available as command now in terminal

## Feature
-   [x] Scaffold New Microservices
-   [x] Generate Service&test and Client From Proto
-   [ ] Go Installer
-   [ ] Proto Installer
-   [ ] Generate Repository (include model)
-   [ ] Add worker with command
-   [ ] Add Nats Subscriber with command
-   [ ] Add db migration file with command





