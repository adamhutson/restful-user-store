# Restful User Store
Restful api for example user data store


Let's implement a restful service using golang that handles basic CRUD on a user object with a minimal nested groups slice.

The user object will contain first/last names, a user id, and a slice of groups.

The groups object will contain names of groups.

Example User Object
```json
{
    "first_name": "Adam",
    "last_name": "Hutson",
    "userid": "ahutson",
    "groups": ["admins", "users"]
}
```

# Endpoints

| Method | Endpoint | Status(es) | Description
| --- | --- | --- | --- |
| GET | /users | 200 | Return all Users w groups
| POST | /users | 201, 409 | Create new User w groups
| GET | /users/{userid} | 200, 404 | Return specified User w groups
| PUT | /users/{userid} | 200, 404 | Update specified User w groups
| DELETE | /users/{userid} | 200, 404 | Delete specified User w groups
| POST | /groups | 201, 409 | Create new Group 
| GET | /groups/{groupname} | 200, 404 | Return specified Group w userIDs
| PUT | /groups/{groupname} | 200, 404 | Update specified Group w userIDs
| DELETE | /groups/{groupname} | 200, 404 | Delete specified Group w userIDs

# Setup

Need golang installed and GOPATH env vars set. 

Need postgres installed locally, and running.  I'm on macOS, and prefer to use Homebrew.
```
brew install postgres
```

Create postgres database
```
createdb restful-user-store
```
Connect to the database and set up a user with permissions.
```
psql restful-user-store
```
While connected to the postgres command prompt, run the following sql commands:
```sql
CREATE USER adam WITH PASSWORD 'password';

GRANT ALL PRIVILEGES ON DATABASE "restful-user-store" to adam;

ALTER USER adam WITH SUPERUSER;
```

# Example Requests

```
curl -i -XGET localhost:8080/users

curl -i -XPOST localhost:8080/users -H "Content-Type, application/json" -d '{"first_name":"Adam","last_name":"Hutson","userid":"ahutson", "groups":["admins", "users"]}'

curl -i -XPOST localhost:8080/users -H "Content-Type, application/json" -d '{"first_name":"Suzanne","last_name":"Hutson","userid":"shutson"}'

curl -i -XGET localhost:8080/users/ahutson

curl -i -XGET localhost:8080/users

ccurl -i -XPOST localhost:8080/users -H "Content-Type, application/json" -d '{"first_name":"Addison","last_name":"Hutson","userid":"ahutson", "groups":["admins", "super_users"]}'

curl -i -XGET localhost:8080/users

curl -i -GET localhost:8080/groups

curl -i -POST localhost:8080/groups -d '{"name":"admin"}'

curl -i -POST localhost:8080/groups -d '{"name":"users"}'

curl -i -POST localhost:8080/groups -d '{"name":"super_users"}'

curl -i -PUT localhost:8080/groups/admins -d {"userids":["ahutson","doesntexist", "meneither"]}

curl -i -XDELETE localhost:8080/users/shutson

curl -i -XDELETE localhost:8080/groups/super_user
```

# Persistance

The service will store the user objects in a SQL store, specifically Postgres.

We will need a couple tables to store the User & Group Objects, and the relation between them.  Don't worry about creating these, as the program will create them if they don't exist.
```sql
CREATE TABLE IF NOT EXISTS users (
    FirstName TEXT, 
    LastName TEXT, 
    UserID TEXT PRIMARY KEY NOT NULL
);

CREATE TABLE IF NOT EXISTS groups (
    GroupName TEXT PRIMARY KEY NOT NULL
);

CREATE TABLE IF NOT EXISTS userGroups (
    UserID TEXT NOT NULL, 
    GroupName TEXT NOT NULL, 
    PRIMARY KEY(UserID, GroupName)
);
```

# Tests

TODO: We will also provide some tests to verify that it's all working as designed.  Yea, yea, ... I know ... it's not TDD.  
