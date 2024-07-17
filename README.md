# future-app

## Introduction
This project is a backend engineering assignment for Future, designed to facilitate appointment scheduling between trainers and clients.

## Getting Started
### Prereq
Create and `.env` file in the project's root using the `.env.example` file as a template
#### ENV Variables
- `PORT`: The port to run the server on. Defaults to `8080`.
- `TEST_PORT`: The port to run the server on during testing. Defaults to `8081`.

### Running Server
1. Initialize and seed the database
```bash
make seed
```
2. Start the server
```bash
// With air
make watch

// Without air
make run
```

### Testing
```bash
make test
```

## Tech Stack
- [Go](https://go.dev)
- [Echo](https://echo.labstack.com)
- [SQLite](https://www.sqlite.org)
- [Zerolog](https://github.com/rs/zerolog)
- [Go Validator](https://github.com/go-playground/validator)

## API

**All incoming datestrings must be in [RFC-3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) date format.
Any datestring that is not in PST -08:00 will be converted as such.**

### `POST /appointments`
Creates an appointment between a user and trainer at a given timeslot

#### Request Body
- `user_id`: The user's ID. Must be GTE 1.
- `trainer_id`: The trainer's ID. Must be GTE 1.
- `starts_at`: The starting time of the appointment in RFC-3339 format (e.g. `2024-07-17T08:00:00-08:00`).
- `ends_at`: The ending time of the appointment in RFC-3339 format (e.g. `2024-07-17T08:00:00-08:00`).

##### Constraints
- Appointments can only be created M-F 8AM-5PM PST (-08:00).
- Appointments must be created at least one hour in advance.
- Appointments can only be 30 minutes long, and should be schedule at :00, :30 minutes after the hour.
- Users/Trainers are only allowed to have one scheduled appointment during a timeslot.

##### Example
```json
{
    "user_id": 1,
    "trainer_id": 1,
    "starts_at": "2030-07-08T15:00:00-08:00",
    "ends_at": "2030-07-08T15:30:00-08:00"
}
```

#### Response
The created appointment is returned in the response

##### 201 Created Example
```json
{
    "id": 10,
    "user_id": 1,
    "trainer_id": 1,
    "starts_at": "2030-07-08T15:00:00-08:00",
    "ends_at": "2030-07-08T15:30:00-08:00"
}
```

##### 400 Example
```json
{
    "message": "Timeslot is not available"
}
```

### `GET /trainers/:trainer_id/appointments`
Returns a list of a trainer's scheduled appointments within a timeframe.

#### Path Parameters
- `trainer_id`: The trainer's ID. Must be GTE 1.

#### Query Parameters
- `starts_at`: The start datetime for the search range in RFC-3339 format.
- `ends_at`: The end datetime fro the search range in RFC-3339 format.

#### Constraints
- The timeframe can be 90 days at most

#### Response
A list of the trainer's appointments within the given timeframe ordered by `starts_at` ascending.
Will return an empty list if no appointments are found.

##### Example
```json
[
    {
        "id": 1,
        "user_id": 1,
        "trainer_id": 1,
        "starts_at": "2019-01-24T09:00:00-08:00",
        "ends_at": "2019-01-24T09:30:00-08:00"
    },
    {
        "id": 2,
        "user_id": 2,
        "trainer_id": 1,
        "starts_at": "2019-01-24T10:00:00-08:00",
        "ends_at": "2019-01-24T10:30:00-08:00"
    },
    {
        "id": 3,
        "user_id": 3,
        "trainer_id": 1,
        "starts_at": "2019-01-25T10:00:00-08:00",
        "ends_at": "2019-01-25T10:30:00-08:00"
    },
    {
        "id": 4,
        "user_id": 4,
        "trainer_id": 1,
        "starts_at": "2019-01-25T10:30:00-08:00",
        "ends_at": "2019-01-25T11:00:00-08:00"
    },
    {
        "id": 5,
        "user_id": 5,
        "trainer_id": 1,
        "starts_at": "2019-01-26T10:00:00-08:00",
        "ends_at": "2019-01-26T10:30:00-08:00"
    }
]
````

### `GET /trainers/:trainer_id/availability`
Returns a list of a trainer's available timeslots in PST (-08:00).

#### Path Parameters
- `trainer_id`: The trainer's ID. Must be GTE 1.

#### Query Parameters
- `starts_at`: The start datetime for the search range in RFC-3339 format.
- `ends_at`: The end datetime fro the search range in RFC-3339 format.

#### Constraints
- The timeframe must be set in the future.
- The timeframe can be 90 days at most.

#### Response
A list of the trainer's available timeslots within the given timeframe.
Will return an empty list if the trainer has no available timeslots.

##### Example
```json
[
    {
        "starts_at": "2025-07-07T08:00:00-08:00",
        "ends_at": "2025-07-07T08:30:00-08:00"
    },
    {
        "starts_at": "2025-07-07T08:30:00-08:00",
        "ends_at": "2025-07-07T09:00:00-08:00"
    },
    {
        "starts_at": "2025-07-07T09:00:00-08:00",
        "ends_at": "2025-07-07T09:30:00-08:00"
    },
    {
        "starts_at": "2025-07-07T09:30:00-08:00",
        "ends_at": "2025-07-07T10:00:00-08:00"
    },
    {
        "starts_at": "2025-07-07T10:00:00-08:00",
        "ends_at": "2025-07-07T10:30:00-08:00"
    },
    {
        "starts_at": "2025-07-07T10:30:00-08:00",
        "ends_at": "2025-07-07T11:00:00-08:00"
    },
    {
        "starts_at": "2025-07-07T11:00:00-08:00",
        "ends_at": "2025-07-07T11:30:00-08:00"
    },
    {
        "starts_at": "2025-07-07T11:30:00-08:00",
        "ends_at": "2025-07-07T12:00:00-08:00"
    },
    {
        "starts_at": "2025-07-07T12:00:00-08:00",
        "ends_at": "2025-07-07T12:30:00-08:00"
    },
    {
        "starts_at": "2025-07-07T12:30:00-08:00",
        "ends_at": "2025-07-07T13:00:00-08:00"
    },
    {
        "starts_at": "2025-07-07T13:00:00-08:00",
        "ends_at": "2025-07-07T13:30:00-08:00"
    },
    {
        "starts_at": "2025-07-07T13:30:00-08:00",
        "ends_at": "2025-07-07T14:00:00-08:00"
    },
    {
        "starts_at": "2025-07-07T14:00:00-08:00",
        "ends_at": "2025-07-07T14:30:00-08:00"
    },
    {
        "starts_at": "2025-07-07T14:30:00-08:00",
        "ends_at": "2025-07-07T15:00:00-08:00"
    },
    {
        "starts_at": "2025-07-07T15:00:00-08:00",
        "ends_at": "2025-07-07T15:30:00-08:00"
    },
    {
        "starts_at": "2025-07-07T15:30:00-08:00",
        "ends_at": "2025-07-07T16:00:00-08:00"
    },
    {
        "starts_at": "2025-07-07T16:00:00-08:00",
        "ends_at": "2025-07-07T16:30:00-08:00"
    },
    {
        "starts_at": "2025-07-07T16:30:00-08:00",
        "ends_at": "2025-07-07T17:00:00-08:00"
    }
]
```

## Note
I changed the fields `started_at` and `ended_at` to `starts_at` and `ends_at` in the file `appointments.json` to keep it consistent with the requirements.
