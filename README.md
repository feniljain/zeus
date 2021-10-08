<p align="center">
		<img src="https://user-images.githubusercontent.com/30529572/60032255-cd613600-96c3-11e9-891f-c2e69c3ce96e.png" width=150px />
		<h2 align="center"> Zeus </h2>
		<h4 align="center"> Yet another meeting scheduler! <h4>
		</p>

## Features
- [X]  Make new meetings with participants
- [X]  Get meeting details from ID
- [X]  Get all meetings of a particular participant
- [X]  Get all meetings within a stipulated timeframe
- [X]  Clean Architecture
- [X]  Not including meetings which clash for a participant
- [X]  Race condition free(Test from golang's default tool)
- [X]  Unit tests
- [X]  Pagination(with the get all meetings of a participant)

<br>


## Instructions to run

* Pre-requisites:
	- Go
	- MongoDB

* Installation:
```bash
git clone https://github.com/feniljain/zeus.git
cd zeus
go get
```

* Execution

```bash
go run api/main.go
```

## Instructions to test

* Testing
```bash
go test api/endpoints_test.go
```

## Architecture
- Insipration from: https://medium.com/gdg-vit/clean-architecture-the-right-way-d83b81ecac6

- This project is built in *Clean Architecture*, it contains of two main modules(not packages), i.e. api and pkg.

- Service acts as usecase layer

- Repo as Repository layer

- Impl as implementation layer

- api contains all the necessary route handlings and backend supporting services(i.e. receiving requests and forwarding to proper handlers), it also initializes everything and does the important step of *dependecy injection*, it contains packages:
	- main: cotnains main.go and testing file.
	- views: contains response wrapper
	- handlers: contains all the necessaey handlers and linking with services, which in turn response using views.
- pkg contains the business logic divided into couple of packages
	- pkg: contains centralized errors.go file defining all the necessary errors which will thrown from backend and pkg
	- meeting: contains all the neccessary files for meeting business logic
	- participant: contains all the neccessary files for participant business logic
	- entities: necessary middle man structs for holding participants and meetings data from db and so forth

## Routes
-  ### Make new meetings with participants
		Route:  /meetings
		Method: POST
		Example Body:
```bash
{
   "title": "sefdkf",
   "starttime": "2009-01-02 15:04:05",
   "endtime": "2009-01-02 15:34:05",
   "participants": [
	{
		"name": "someone",
		"email": "someone@gmail.com",
		"rsvp": "Yes"
	},
	{
		"name": "someone1",
		"email": "someone1@gmail.com",
		"rsvp": "No"
	}
   ]
}
```
- ###  Get meeting details from ID
		Route:  /meetings/:id
		Method: GET
- ###  Get all meetings of a particular participant
		Route:
		 - /meetings?participant=sw
		 - /meetings?participant=sw&page=1
		Method: GET
- ###  Get all meetings within a stipulated timeframe
		Route:  /meetings?start=2005-01-02 15:04:05&end=2007-01-02 15:34:05
		Method: GET

## Contributors

<table>
<tr align="center">


<td>

Fenil Jain

<p align="center">
<img src = "https://avatars2.githubusercontent.com/u/49019259?s=460&u=231aa9f5647e68a27939c49eb7e56cfc6412db3b&v=4" width="150" height="150" alt="Your Name Here (Insert Your Image Link In Src">
</p>
<p align="center">
<a href = "https://github.com/feniljain"><img src = "http://www.iconninja.com/files/241/825/211/round-collaboration-social-github-code-circle-network-icon.svg" width="36" height = "36" alt="GitHub"/></a>
<a href = "https://www.linkedin.com/in/fenil-jain-b3711117b">
<img src = "http://www.iconninja.com/files/863/607/751/network-linkedin-social-connection-circular-circle-media-icon.svg" width="36" height="36" alt="LinkedIn"/>
</a>
</p>
</td>

</tr>
  </table>

<p align="center">
	Made with :heart: by <a href="https://github.com/feniljain">Fenil Jain</a>
</p>

