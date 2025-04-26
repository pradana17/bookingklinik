
# bookingklinik

**System Booking Klinik / Dokter**

This project is a clinic booking system that allows users to book appointments with doctors. It handles the scheduling, user roles (patients, doctors, admins), and maintains the integrity of booking data.

## Project Overview

`bookingklinik` is a Go-based web application that allows users (patients) to book appointments with doctors. The system includes features like:

- Doctor schedules management
- User roles (patients, doctors, and admins)
- Booking creation and management
- Pagination for listing bookings

## Installation Instructions

1. Clone the repository to your local machine:
   ```bash
   git clone https://github.com/your-username/bookingklinik.git
   cd bookingklinik
   ```

2. Install dependencies:
   Make sure you have Go installed on your machine. If not, you can download it from [the official Go website](https://golang.org/).

   Then, run:
   ```bash
   go mod tidy
   ```

3. Create a `.env` file with the necessary environment variables. You can use `.env.sample` as a template.

4. Run the application:
   ```bash
   go run main.go
   ```

5. The application should now be running on `http://localhost:8080`.

## Folder Structure

```
.
├── config/             # Configuration files for the application (e.g., DB setup)
├── controllers/        # API controllers for handling requests and responses
├── middleware/         # Middleware for handling things like authentication
├── model/              # Model definitions (e.g., User, Doctor, Booking)
├── repository/         # Repository layer for interacting with the database
├── routes/             # Routes for the web application
├── services/           # Service layer for business logic
├── utils/              # Utility functions and helper methods
├── .env                # Environment variables file (you need to create this)
├── go.mod              # Go module file for dependencies
├── go.sum              # Go checksum file
├── main.go             # Entry point for the application
└── README.md           # Project documentation
```

## Usage


## Technologies

- **Go**: The backend language.
- **Gin**: The web framework for routing and handling HTTP requests.
- **GORM**: ORM for interacting with the database.
- **MySQL**: The database (or another relational DB if preferred).
- **JWT**: JSON Web Tokens for authentication.

## Authentication

- **Admin**: Can manage all bookings, view all doctors' schedules, and manage users.
- **Doctor**: Can manage their own schedule and view patient bookings.
- **Patient**: Can book appointments with available doctors and view their own bookings.
All endpoints that require authentication use **JWT tokens** for validation. You must obtain a valid token by logging in through the `/login` endpoint.

## Endpoints API

### User Routes

| **Endpoint**        | **Method** | **Description**                          | **Authentication** | **Roles**   |
|---------------------|------------|------------------------------------------|--------------------|-------------|
| `/register`         | POST       | Register a new user                      | None               | All Users   |
| `/login`            | POST       | Log in to get a JWT token                | None               | All Users   |
| `/user/password`    | PUT        | Update user password                     | Required JWT       | All Users   |

### Booking Routes

| **Endpoint**                  | **Method** | **Description**                                    | **Authentication** | **Roles**  |
|-------------------------------|------------|----------------------------------------------------|--------------------|------------|
| `/booking`                     | POST       | Create a new booking                              | Required JWT       | Patient    |
| `/booking`                     | GET        | Get all bookings                                  | Required JWT       | Patient    |
| `/booking/:id`                 | GET        | Get booking by ID                                 | Required JWT       | Patient    |
| `/booking/user/:user_id`       | GET        | Get bookings by user ID                           | Required JWT       | Patient    |
| `/booking/doctor/:doctor_id`   | GET        | Get bookings by doctor ID                         | Required JWT       | Patient    |
| `/booking/:id`                 | PUT        | Update booking by ID                              | Required JWT       | Patient    |
| `/booking/:id`                 | DELETE     | Delete booking by ID                              | Required JWT       | Patient    |

### Doctor Routes

| **Endpoint**        | **Method** | **Description**                          | **Authentication**      | **Roles** |
|---------------------|------------|------------------------------------------|-------------------------|----------|
| `/doctor`           | POST       | Create a new doctor record               | Required JWT            | Admin    |
| `/doctor`           | GET        | Get all doctors                          | Required JWT            | Admin    |
| `/doctor/:id`       | GET        | Get doctor by ID                         | Required JWT            | Admin    |
| `/doctor/:id`       | PUT        | Update doctor by ID                      | Required JWT            | Admin    |
| `/doctor/:id`       | DELETE     | Delete doctor by ID                      | Required JWT            | Admin    |

### Doctor Schedule Routes

| **Endpoint**                  | **Method** | **Description**                                    | **Authentication**      | **Roles**   |
|-------------------------------|------------|----------------------------------------------------|-------------------------|-------------|
| `/doctorschedule`              | POST       | Create a new doctor schedule                       | Required JWT            | Admin, Doctor |
| `/doctorschedule`              | GET        | Get all doctor schedules                           | Required JWT            | Admin, Doctor |
| `/doctorschedule/:id`          | GET        | Get doctor schedule by ID                          | Required JWT            | Admin, Doctor |
| `/doctorschedule/:id`          | PUT        | Update doctor schedule by ID                       | Required JWT            | Admin, Doctor |
| `/doctorschedule/:id`          | DELETE     | Delete doctor schedule by ID                       | Required JWT            | Admin, Doctor |

### Service Routes

| **Endpoint**                  | **Method** | **Description**                                    | **Authentication**      | **Roles**   |
|-------------------------------|------------|----------------------------------------------------|-------------------------|-------------|
| `/service`                     | GET        | Get all available services                         | Required JWT            | All Users   |
| `/service/:id`                 | GET        | Get service by ID                                  | Required JWT            | All Users   |
| `/service`                     | POST       | Create a new service (Admin only)                  | Required JWT            | Admin       |
| `/service/:id`                 | PUT        | Update service by ID                               | Required JWT            | Admin       |
| `/service/:id`                 | DELETE     | Delete service by ID                               | Required JWT            | Admin       |

---

### Example Collection (Postman)
Here is a sample collection for Postman to test the API:

```json
{
  "collection": {
    "info": {
      "_postman_id": "58d347ca-5811-46e1-9f78-8f0d9186e35b",
      "name": "PROJECT",
      "schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
      "updatedAt": "2025-04-26T20:47:20.000Z",
      "createdAt": "2025-04-21T01:48:42.000Z",
      "lastUpdatedBy": "39611293",
      "uid": "39611293-58d347ca-5811-46e1-9f78-8f0d9186e35b"
    },
    "item": [
      {
        "name": "Patient",
        "item": [
          {
            "name": "Create Booking",
            "id": "23b1a353-d9c7-41f8-809f-e3acc647c53f",
            "protocolProfileBehavior": {
              "disableBodyPruning": true
            },
            "request": {
              "auth": {
                "type": "bearer",
                "bearer": [
                  {
                    "key": "token",
                    "value": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjo1LCJlbWFpbCI6InNhbG1hbkBnbWFpbC5jb20iLCJyb2xlIjoicGF0aWVudCIsImV4cCI6MTc0NTY5NDEzNX0.bOuDYvZPZOFCBqWDRDospFcGRLTbYo3nWM-0wLN2o6E",
                    "type": "string"
                  }
                ]
              },
              "method": "POST",
              "header": [],
              "body": {
                "mode": "raw",
                "raw": "{\r\n    \"doctor_id\":3,\r\n    \"service_id\":1,\r\n    \"booking_date\":\"2025-04-28\",\r\n    \"booking_time\":\"14:00\",\r\n    \"notes\":\"sakit parah\"\r\n}",
                "options": {
                  "raw": {
                    "language": "json"
                  }
                }
              },
              "url": {
                "raw": "localhost:8080/booking",
                "host": [
                  "localhost"
                ],
                "port": "8080",
                "path": [
                  "booking"
                ]
              }
            },
            "response": [
              {
                "id": "d0ac6e48-aa0e-4523-97b1-9e1666527962",
                "name": "Create Booking",
                "originalRequest": {
                  "method": "POST",
                  "header": [],
                  "body": {
                    "mode": "raw",
                    "raw": "{\r\n    \"doctor_id\":3,\r\n    \"service_id\":1,\r\n    \"booking_date\":\"2025-04-28\",\r\n    \"booking_time\":\"14:00\",\r\n    \"notes\":\"sakit parah\"\r\n}",
                    "options": {
                      "raw": {
                        "language": "json"
                      }
                    }
                  },
                  "url": {
                    "raw": "localhost:8080/booking",
                    "host": [
                      "localhost"
                    ],
                    "port": "8080",
                    "path": [
                      "booking"
                    ]
                  }
                },
                "status": "OK",
                "code": 200,
                "header": [
                  {
                    "key": "Content-Type",
                    "value": "application/json; charset=utf-8"
                  },
                  {
                    "key": "Date",
                    "value": "Sat, 26 Apr 2025 18:04:18 GMT"
                  },
                  {
                    "key": "Content-Length",
                    "value": "244"
                  }
                ],
                "responseTime": null,
                "body": "{\n    \"booking\": {\n        \"id\": 14,\n        \"doctor_name\": \"John\",\n        \"service_name\": \"General Check Up\",\n        \"booking_date\": \"2025-04-28T00:00:00+07:00\",\n        \"booking_time\": \"2025-04-28T14:00:00+07:00\",\n        \"status\": \"pending\",\n        \"notes\": \"sakit parah\"\n    },\n    \"message\": \"Booking created successfully\"\n}",
                "uid": "39611293-d0ac6e48-aa0e-4523-97b1-9e1666527962"
              }
            ]
          }
        ]
      }
    ]
  }
}