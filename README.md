# jsonGO
jsonGo is a golang app that handle transaction request and gonna processed to standard ISO:8583, it included database connection to save any request that it get.

## Installation
Get the program
```bash
git clone https://github.com/wrmn/jsonGo
cd jsonGo
```
Prepare package
```bash
go get github.com/gorilla/mux
go get github.com/go-sql-driver/mysql
```

Running program
```bash
go run *.go
```
Build and running program
```bash
go build -o {program name}
./{program name}
```

## Available Request (For Now)
#### GET
- `/payment` : Get all transaction data
- `/payment/{id}` : Get transaction data based on processingCode
#### POST
- `/payment` : Post and insert a new data to database based on JSON body that will send as data
#### PUT
- `/payment/{id}` : Put and update specific data with processing code that send as path with JSON body as updated data 
#### DELETE 
- `/payment/{id}` : Delete data with sended processing code from database

## ToDos
- Showing Help when go to any method that not allowed, not found method, or home page
- inspect JSON that receivedas a body and if mandatory field is empty return an error response

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)
