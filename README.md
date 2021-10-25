# CRUD application using Go with MySQL and REST APIs simple demo
 
This is a demo of creating a simple CRUD application using Go, MySQL and REST APIs

Technology stack below is in use:

* Go
* MySQL
* REST APIs
* Mountebank 

## Usage


## Sending REST Requests

#GET all employees
```
http://localhost:8080/employee
```

#GET emoloyee by id 
```
http://localhost:8080/employee/:id
```

#POST create an emoloyee 
```
http://localhost:8080/employee/:id
```

#PUT update an emoloyee 
```
http://localhost:8080/employee/:id
```

#DELETE delete an emoloyee 
```
http://localhost:8080/employee/:id
```


## Mountebank for stubbing/mocking endpoints with varing condtions
To get employees' super balance, an external api is called.
You need to run both the main app and mountebank to demonstrate a proper response from endpoints 

#Run Mountebank config file
```
mb --configfile imposter.json
```

Run the following command in a terminal environment:
```
run main.go

```
