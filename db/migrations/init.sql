CREATE TABLE IF NOT EXISTS rooms (
    room_id serial primary key,
    capacity integer not null,
    room_name char(50),
    created_date timestamp not null default now(),
    updated_date timestamp not null default now()
    
);

CREATE TABLE IF NOT EXISTS users (
    user_id serial primary key,
    first_name char(50),
    last_name char(50),
    email char(50),
    ticket_remaining integer,
    room_id integer references rooms (room_id),
    created_date timestamp not null default now(),
    updated_date timestamp not null default now()
);

CREATE TABLE IF NOT EXISTS events (
    event_id serial primary key,
    event_name char(50),
    event_description text,
    event_date timestamp not null,
    room_id integer references rooms (room_id),
    created_date timestamp not null default now(),
    updated_date timestamp not null default now()
);

WITH inserted AS (
    INSERT INTO rooms (capacity, room_name) VALUES (5, 'testroom') RETURNING room_id )
INSERT INTO users (first_name, last_name, email, ticket_remaining, room_id) 
VALUES ('test', 'user', 'test_user@email.com', 100, (select room_id from inserted));